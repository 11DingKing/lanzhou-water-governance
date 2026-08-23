package service

import (
	"context"
	"fmt"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
)

type Reporting struct{ Repo repository.Reporting }

func (s Reporting) StationSummary(ctx context.Context, user domain.User, region string, page domain.Page) ([]repository.ReportRow, error) {
	if user.Role != domain.RoleAdmin && user.Region != region && region != "" {
		return nil, domain.ErrForbidden
	}
	return s.Repo.StationSummary(ctx, region, page)
}
func (s Reporting) Exposure(ctx context.Context, user domain.User, region string) (map[string]int64, error) {
	if user.Role != domain.RoleAdmin && user.Region != region {
		return nil, domain.ErrForbidden
	}
	waste, err := s.Repo.WasteByRegion(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("exposure waste: %w", err)
	}
	return map[string]int64{"waste_kg": waste}, nil
}
func (s Reporting) Progress(ctx context.Context, user domain.User, projectID int64) (float64, error) {
	if user.Role != domain.RoleAdmin && user.Role != domain.RoleRegional {
		return 0, domain.ErrForbidden
	}
	return s.Repo.ProjectProgress(ctx, projectID)
}
