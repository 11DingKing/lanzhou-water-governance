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

type Workflow struct {
	DB          *sql.DB
	Monitoring  repository.Monitoring
	Inspections repository.Inspections
	Audit       repository.Audit
}

func (w Workflow) OpenAndAssign(ctx context.Context, user domain.User, stationID, sampleID int64, severity string, due time.Time) (domain.Alert, domain.Inspection, error) {
	if user.Role != domain.RoleAdmin && user.Role != domain.RoleInspector {
		return domain.Alert{}, domain.Inspection{}, domain.ErrForbidden
	}
	var alert domain.Alert
	var inspection domain.Inspection
	err := sqlite.WithTx(ctx, w.DB, func(tx *sql.Tx) error {
		var err error
		alert, err = w.Monitoring.OpenAlertTx(ctx, tx, stationID, sampleID, severity)
		if err != nil {
			return err
		}
		inspection, err = w.Inspections.CreateTx(ctx, tx, alert.ID, user.ID, due)
		if err != nil {
			return err
		}
		return w.Audit.RecordTx(ctx, tx, user.ID, requestID(ctx), "alert", fmt.Sprint(alert.ID), "assign", "ok", map[string]any{"inspection": inspection.ID})
	})
	if err != nil {
		return alert, inspection, fmt.Errorf("open and assign: %w", err)
	}
	return alert, inspection, nil
}
func (w Workflow) ResolveAfterEvidence(ctx context.Context, user domain.User, alertID, inspectionID int64, alertVersion, inspectionVersion int64, evidence string) (domain.Alert, domain.Inspection, error) {
	if user.Role != domain.RoleAdmin && user.Role != domain.RoleInspector {
		return domain.Alert{}, domain.Inspection{}, domain.ErrForbidden
	}
	var alert domain.Alert
	var inspection domain.Inspection
	err := sqlite.WithTx(ctx, w.DB, func(tx *sql.Tx) error {
		current, err := w.Inspections.Get(ctx, inspectionID)
		if err != nil {
			return err
		}
		if current.Status != domain.InspectionRunning {
			return domain.ErrInvalidState
		}
		inspection, err = w.Inspections.Transition(ctx, inspectionID, current.Status, domain.InspectionCompleted, inspectionVersion, evidence)
		if err != nil {
			return err
		}
		currentAlert, err := w.Monitoring.Alert(ctx, alertID)
		if err != nil {
			return err
		}
		alert, err = w.Monitoring.TransitionAlert(ctx, alertID, currentAlert.Status, domain.AlertResolved, alertVersion)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return alert, inspection, fmt.Errorf("resolve after evidence: %w", err)
	}
	return alert, inspection, nil
}
func (w Workflow) CancelIfContextDone(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("workflow cancelled: %w", ctx.Err())
	default:
		return nil
	}
}

func RunAtomic(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error { return sqlite.WithTx(ctx, db, fn) }
