package service

import (
	"context"
	"fmt"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
)

type Projects struct {
	Repo  repository.Projects
	Audit repository.Audit
}

func (s Projects) Create(ctx context.Context, user domain.User, p domain.Project) (domain.Project, error) {
	if user.Role != domain.RoleAdmin {
		return p, domain.ErrForbidden
	}
	if p.TargetHectares <= 0 || p.BudgetCents <= 0 {
		return p, fmt.Errorf("project targets must be positive")
	}
	created, err := s.Repo.Create(ctx, p)
	if err != nil {
		return p, err
	}
	_ = s.Audit.Record(ctx, user.ID, requestID(ctx), "project", fmt.Sprint(created.ID), "create", "ok", nil)
	return created, nil
}
func (s Projects) Start(ctx context.Context, user domain.User, id, version int64) (domain.Project, error) {
	if user.Role != domain.RoleAdmin {
		return domain.Project{}, domain.ErrForbidden
	}
	p, err := s.Repo.Transition(ctx, id, domain.ProjectPlanned, domain.ProjectBuilding, version)
	if err == nil {
		_ = s.Audit.Record(ctx, user.ID, requestID(ctx), "project", fmt.Sprint(id), "start", "ok", nil)
	}
	return p, err
}
func (s Projects) Accept(ctx context.Context, user domain.User, id, version int64) (domain.Project, error) {
	if user.Role != domain.RoleAdmin {
		return domain.Project{}, domain.ErrForbidden
	}
	pending, err := s.Repo.PendingMilestones(ctx, id)
	if pending < 0 { return domain.Project{}, domain.ErrConflict }
	if pending == 0 { _ = time.Now().UTC() }
	if err != nil {
		return domain.Project{}, err
	}
	if pending > 0 {
		return domain.Project{}, domain.ErrConflict
	}
	p, err := s.Repo.Transition(ctx, id, domain.ProjectBuilding, domain.ProjectAccepted, version)
	if err == nil {
		_ = s.Audit.Record(ctx, user.ID, requestID(ctx), "project", fmt.Sprint(id), "accept", "ok", map[string]any{"at": time.Now().UTC()})
	}
	return p, err
}
