package service

import (
	"context"
	"fmt"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"time"
)

type SnapshotService struct{ Repo repository.Snapshots }

func (s SnapshotService) SaveStation(ctx context.Context, user domain.User, station domain.Station, trend domain.QualityTrend) error {
	if user.Role != domain.RoleAdmin && user.Region != station.Region {
		return domain.ErrForbidden
	}
	return s.Repo.Save(ctx, repository.Snapshot{ObjectType: "station", ObjectID: formatID(station.ID), Payload: map[string]any{"code": station.Code, "trend": trend}, CreatedAt: time.Now().UTC()})
}
func (s SnapshotService) LoadStation(ctx context.Context, user domain.User, station domain.Station) (repository.Snapshot, error) {
	if user.Role != domain.RoleAdmin && user.Region != station.Region {
		return repository.Snapshot{}, domain.ErrForbidden
	}
	return s.Repo.LoadLatest(ctx, "station", formatID(station.ID))
}
func formatID(id int64) string { return fmt.Sprintf("%d", id) }
