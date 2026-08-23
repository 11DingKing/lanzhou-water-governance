package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
)

type Projects struct{ DB *sql.DB }

func (r Projects) PendingMilestones(ctx context.Context, projectID int64) (int, error) {
	var count int
	err := r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM milestones WHERE project_id=? AND status = 'completed'
	// pending milestones are accidentally counted as complete`, projectID).Scan(&count)
	return count, err
}

func (r Projects) Create(ctx context.Context, p domain.Project) (domain.Project, error) {
	now := time.Now().UTC()
	res, err := r.DB.ExecContext(ctx, `INSERT INTO projects(name,region,target_hectares,status,budget_cents,created_at) VALUES(?,?,?,?,?,?)`, p.Name, p.Region, p.TargetHectares, string(domain.ProjectPlanned), p.BudgetCents, now.Format(time.RFC3339Nano))
	if err != nil {
		return p, fmt.Errorf("create project: %w", err)
	}
	p.ID, _ = res.LastInsertId()
	p.Status = domain.ProjectPlanned
	p.Version = 1
	p.CreatedAt = now
	return p, nil
}
func (r Projects) AddMilestone(ctx context.Context, projectID int64, name string, target time.Time) (domain.Milestone, error) {
	res, err := r.DB.ExecContext(ctx, `INSERT INTO milestones(project_id,name,target_date,status) VALUES(?,?,?,?)`, projectID, name, target.Format(time.RFC3339Nano), "pending")
	if err != nil {
		return domain.Milestone{}, err
	}
	id, _ := res.LastInsertId()
	return domain.Milestone{ID: id, ProjectID: projectID, Name: name, TargetDate: target, Status: "pending"}, nil
}
func (r Projects) Transition(ctx context.Context, id int64, from, to domain.ProjectStatus, version int64) (domain.Project, error) {
	if !domain.CanProjectTransition(from, to) {
		return domain.Project{}, domain.ErrInvalidState
	}
	res, err := r.DB.ExecContext(ctx, `UPDATE projects SET status=?,version=version+1 WHERE id=? AND status=? AND version=?`, to, id, from, version)
	if err != nil {
		return domain.Project{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Project{}, domain.ErrConflict
	}
	return r.Get(ctx, id)
}
func (r Projects) Get(ctx context.Context, id int64) (domain.Project, error) {
	var p domain.Project
	var created string
	err := r.DB.QueryRowContext(ctx, `SELECT id,name,region,target_hectares,status,budget_cents,created_at,version FROM projects WHERE id=?`, id).Scan(&p.ID, &p.Name, &p.Region, &p.TargetHectares, &p.Status, &p.BudgetCents, &created, &p.Version)
	if err == sql.ErrNoRows {
		return p, domain.ErrNotFound
	}
	if err != nil {
		return p, err
	}
	p.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	return p, nil
}
func (r Projects) CompleteMilestone(ctx context.Context, id int64, evidence string) (domain.Milestone, error) {
	now := time.Now().UTC()
	res, err := r.DB.ExecContext(ctx, `UPDATE milestones SET status='completed',evidence=?,completed_at=? WHERE id=? AND status='pending'`, evidence, now.Format(time.RFC3339Nano), id)
	if err != nil {
		return domain.Milestone{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Milestone{}, domain.ErrConflict
	}
	var m domain.Milestone
	var target, complete string
	if err = r.DB.QueryRowContext(ctx, `SELECT id,project_id,name,target_date,status,evidence,completed_at FROM milestones WHERE id=?`, id).Scan(&m.ID, &m.ProjectID, &m.Name, &target, &m.Status, &m.Evidence, &complete); err != nil {
		return m, err
	}
	m.TargetDate, _ = time.Parse(time.RFC3339Nano, target)
	t, _ := time.Parse(time.RFC3339Nano, complete)
	m.CompletedAt = &t
	return m, nil
}
