package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
)

type Collaboration struct{ DB *sql.DB }

func (r Collaboration) Agreement(ctx context.Context, upstream, downstream string) (domain.Agreement, error) {
	var a domain.Agreement
	var active int
	var created string
	err := r.DB.QueryRowContext(ctx, `SELECT id,upstream_region,downstream_region,threshold_class,active,created_at FROM agreements WHERE upstream_region=? AND downstream_region=?`, upstream, downstream).Scan(&a.ID, &a.UpstreamRegion, &a.DownstreamRegion, &a.ThresholdClass, &active, &created)
	if err == sql.ErrNoRows {
		return a, domain.ErrNotFound
	}
	if err != nil {
		return a, err
	}
	a.Active = active != 0
	a.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	return a, nil
}
func (r Collaboration) CreateAgreement(ctx context.Context, a domain.Agreement) (domain.Agreement, error) {
	now := time.Now().UTC()
	res, err := r.DB.ExecContext(ctx, `INSERT INTO agreements(upstream_region,downstream_region,threshold_class,created_at) VALUES(?,?,?,?)`, a.UpstreamRegion, a.DownstreamRegion, a.ThresholdClass, now.Format(time.RFC3339Nano))
	if err != nil {
		return a, fmt.Errorf("create agreement: %w", err)
	}
	a.ID, _ = res.LastInsertId()
	a.Active = true
	a.CreatedAt = now
	return a, nil
}
func (r Collaboration) CreateWarningTx(ctx context.Context, tx *sql.Tx, a domain.Agreement, stationID int64, direction string, payload map[string]any) (domain.Warning, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return domain.Warning{}, err
	}
	now := time.Now().UTC()
	res, err := tx.ExecContext(ctx, `INSERT INTO warnings(agreement_id,station_id,direction,payload_json,status,sent_at) VALUES(?,?,?,?,?,?)`, a.ID, stationID, direction, string(raw), "sent", now.Format(time.RFC3339Nano))
	if err != nil {
		return domain.Warning{}, err
	}
	id, _ := res.LastInsertId()
	return domain.Warning{ID: id, AgreementID: a.ID, StationID: stationID, Direction: direction, Payload: string(raw), Status: "sent", SentAt: now}, nil
}
func (r Collaboration) CreateCompensationTx(ctx context.Context, tx *sql.Tx, a domain.Agreement, period, direction, reason string, amount int64) (domain.Compensation, error) {
	now := time.Now().UTC()
	res, err := tx.ExecContext(ctx, `INSERT INTO compensations(agreement_id,period,direction,amount_cents,reason,status,created_at) VALUES(?,?,?,?,?,?,?)`, a.ID, period, direction, amount, reason, "pending", now.Format(time.RFC3339Nano))
	if err != nil {
		return domain.Compensation{}, err
	}
	id, _ := res.LastInsertId()
	return domain.Compensation{ID: id, AgreementID: a.ID, Period: period, Direction: direction, Reason: reason, AmountCents: amount, Status: "pending", CreatedAt: &now}, nil
}
func (r Collaboration) SettleCompensation(ctx context.Context, id int64) (domain.Compensation, error) {
	now := time.Now().UTC()
	res, err := r.DB.ExecContext(ctx, `UPDATE compensations SET status='settled',settled_at=? WHERE id=? AND status='pending'`, now.Format(time.RFC3339Nano), id)
	if err != nil {
		return domain.Compensation{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Compensation{}, domain.ErrConflict
	}
	var c domain.Compensation
	var created, settled string
	if err = r.DB.QueryRowContext(ctx, `SELECT id,agreement_id,period,direction,amount_cents,reason,status,created_at,settled_at FROM compensations WHERE id=?`, id).Scan(&c.ID, &c.AgreementID, &c.Period, &c.Direction, &c.AmountCents, &c.Reason, &c.Status, &created, &settled); err != nil {
		return c, err
	}
	ct, _ := time.Parse(time.RFC3339Nano, created)
	st, _ := time.Parse(time.RFC3339Nano, settled)
	c.CreatedAt = &ct
	c.SettledAt = &st
	return c, nil
}

func ConstraintMessage(err error) string { return err.Error() }
