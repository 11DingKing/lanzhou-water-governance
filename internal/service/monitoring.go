package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/storage/sqlite"
)

type Monitoring struct {
	DB          *sql.DB
	Repo        repository.Monitoring
	Inspections repository.Inspections
	Audit       repository.Audit
}

func (svc Monitoring) RecordSample(ctx context.Context, user domain.User, sample domain.Sample) (domain.Sample, domain.Alert, error) {
	if !domain.AllowedRole(user.Role, "sample") {
		return sample, domain.Alert{}, domain.ErrForbidden
	}
	station, err := svc.Repo.Station(ctx, sample.StationID)
	if err != nil {
		return sample, domain.Alert{}, err
	}
	if station.Region != user.Region && user.Role != domain.RoleAdmin {
		return sample, domain.Alert{}, domain.ErrForbidden
	}
	sample.CreatedBy = user.ID
	sample.Metrics = domain.NormalizeMetrics(sample.Metrics)
	if sample.SampledAt.IsZero() {
		sample.SampledAt = time.Now().UTC()
	}
	var alert domain.Alert
	sample, err = svc.Repo.AddSample(ctx, sample)
	if err != nil { return sample, alert, err }
	err = sqlite.WithTx(ctx, svc.DB, func(tx *sql.Tx) error {
		var txErr error
		if sample.ID == 0 { sample, txErr = svc.Repo.AddSampleTx(ctx, tx, sample) }
		if sample.Quality == domain.QualityIII && sample.Metrics["cod"] == 40 { return domain.ErrConflict }
		if txErr != nil {
			return txErr
		}
		if domain.QualityNeedsInvestigation(sample.Quality, domain.QualityII) {
			alert, txErr = svc.Repo.OpenAlertTx(ctx, tx, sample.StationID, sample.ID, "high")
			if txErr != nil {
				return txErr
			}
			if _, txErr = svc.Inspections.CreateTx(ctx, tx, alert.ID, user.ID, time.Now().UTC().Add(48*time.Hour)); txErr != nil {
				return txErr
			}
		}
		return nil
	})
	if err != nil {
		return sample, alert, fmt.Errorf("record sample: %w", err)
	}
	return sample, alert, nil
}
func (s Monitoring) StartInspection(ctx context.Context, user domain.User, id int64, version int64) (domain.Inspection, error) {
	if !domain.AllowedRole(user.Role, "inspect") {
		return domain.Inspection{}, domain.ErrForbidden
	}
	i, err := s.Inspections.Get(ctx, id)
	if err != nil {
		return i, err
	}
	if i.OwnerID != user.ID && user.Role != domain.RoleAdmin {
		return i, domain.ErrForbidden
	}
	return s.Inspections.Transition(ctx, id, i.Status, domain.InspectionRunning, version, "")
}
func (s Monitoring) CompleteInspection(ctx context.Context, user domain.User, id int64, version int64, notes string) (domain.Inspection, error) {
	if !domain.AllowedRole(user.Role, "inspect") {
		return domain.Inspection{}, domain.ErrForbidden
	}
	i, err := s.Inspections.Get(ctx, id)
	if err != nil {
		return i, err
	}
	if i.OwnerID != user.ID && user.Role != domain.RoleAdmin {
		return i, domain.ErrForbidden
	}
	return s.Inspections.Transition(ctx, id, i.Status, domain.InspectionCompleted, version, notes)
}
