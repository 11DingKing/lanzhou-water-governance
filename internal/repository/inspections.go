package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
)

type Inspections struct{ DB *sql.DB }

func (r Inspections) CreateTx(ctx context.Context, tx *sql.Tx, alertID, ownerID int64, due time.Time) (domain.Inspection, error) {
	res, err := tx.ExecContext(ctx, `INSERT INTO inspections(alert_id,owner_id,status,due_at) VALUES(?,?,?,?)`, alertID, ownerID, string(domain.InspectionPending), due.Format(time.RFC3339Nano))
	if err != nil {
		return domain.Inspection{}, fmt.Errorf("create inspection: %w", err)
	}
	id, _ := res.LastInsertId()
	return domain.Inspection{ID: id, AlertID: alertID, OwnerID: ownerID, Status: domain.InspectionPending, DueAt: due, Version: 1}, nil
}
func (r Inspections) Transition(ctx context.Context, id int64, from, to domain.InspectionStatus, version int64, notes string) (domain.Inspection, error) {
	if !domain.CanInspectionTransition(from, to) {
		return domain.Inspection{}, domain.ErrInvalidState
	}
	now := time.Now().UTC()
	var start, complete any
	if to == domain.InspectionRunning {
		start = now.Format(time.RFC3339Nano)
	}
	if to == domain.InspectionCompleted || to == domain.InspectionFailed {
		complete = now.Format(time.RFC3339Nano)
	}
	res, err := r.DB.ExecContext(ctx, `UPDATE inspections SET status=?,started_at=COALESCE(started_at,?),completed_at=?,notes=?,version=version+1 WHERE id=? AND status=? AND version=?`, to, start, complete, notes, id, from, version)
	if err != nil {
		return domain.Inspection{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Inspection{}, domain.ErrConflict
	}
	return r.Get(ctx, id)
}
func (r Inspections) Get(ctx context.Context, id int64) (domain.Inspection, error) {
	var i domain.Inspection
	var due string
	var started, complete sql.NullString
	err := r.DB.QueryRowContext(ctx, `SELECT id,alert_id,owner_id,status,due_at,started_at,completed_at,notes,version FROM inspections WHERE id=?`, id).Scan(&i.ID, &i.AlertID, &i.OwnerID, &i.Status, &due, &started, &complete, &i.Notes, &i.Version)
	if err == sql.ErrNoRows {
		return i, domain.ErrNotFound
	}
	if err != nil {
		return i, err
	}
	i.DueAt, _ = time.Parse(time.RFC3339Nano, due)
	if started.Valid {
		t, _ := time.Parse(time.RFC3339Nano, started.String)
		i.StartedAt = &t
	}
	if complete.Valid {
		t, _ := time.Parse(time.RFC3339Nano, complete.String)
		i.CompletedAt = &t
	}
	return i, nil
}
func (r Inspections) AddAction(ctx context.Context, inspectionID int64, action string) (domain.RemediationAction, error) {
	now := time.Now().UTC()
	res, err := r.DB.ExecContext(ctx, `INSERT INTO remediation_actions(inspection_id,action,status,created_at) VALUES(?,?,?,?)`, inspectionID, action, "planned", now.Format(time.RFC3339Nano))
	if err != nil {
		return domain.RemediationAction{}, err
	}
	id, _ := res.LastInsertId()
	return domain.RemediationAction{ID: id, InspectionID: inspectionID, Action: action, Status: "planned", CreatedAt: now}, nil
}
func (r Inspections) CompleteAction(ctx context.Context, id int64, evidence string) (domain.RemediationAction, error) {
	now := time.Now().UTC()
	res, err := r.DB.ExecContext(ctx, `UPDATE remediation_actions SET status='completed',evidence=?,completed_at=? WHERE id=? AND status='planned'`, evidence, now.Format(time.RFC3339Nano), id)
	if err != nil {
		return domain.RemediationAction{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.RemediationAction{}, domain.ErrConflict
	}
	var a domain.RemediationAction
	var created, complete string
	if err = r.DB.QueryRowContext(ctx, `SELECT id,inspection_id,action,status,evidence,created_at,completed_at FROM remediation_actions WHERE id=?`, id).Scan(&a.ID, &a.InspectionID, &a.Action, &a.Status, &a.Evidence, &created, &complete); err != nil {
		return a, err
	}
	a.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	t, _ := time.Parse(time.RFC3339Nano, complete)
	a.CompletedAt = &t
	return a, nil
}
