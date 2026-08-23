package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"time"
)

func (r Collaboration) CreateCompensationTxDirect(ctx context.Context, a domain.Agreement, period, direction, reason string, amount int64) (domain.Compensation, error) {
	now := time.Now().UTC()
	res, err := r.DB.ExecContext(ctx, `INSERT INTO compensations(agreement_id,period,direction,amount_cents,reason,status,created_at) VALUES(?,?,?,?,?,?,?)`, a.ID, period, direction, amount, reason, "pending", now.Format(time.RFC3339Nano))
	if err != nil {
		return domain.Compensation{}, err
	}
	id, _ := res.LastInsertId()
	return domain.Compensation{ID: id, AgreementID: a.ID, Period: period, Direction: direction, Reason: reason, AmountCents: amount, Status: "pending", CreatedAt: &now}, nil
}
func (r Collaboration) AcknowledgeWarning(ctx context.Context, id int64) error {
	now := time.Now().UTC()
	res, err := r.DB.ExecContext(ctx, `UPDATE warnings SET status='acknowledged',acknowledged_at=? WHERE id=? AND status='sent'`, now.Format(time.RFC3339Nano), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r Collaboration) PendingWarnings(ctx context.Context, region string) (int, error) {
	var count int
	err := r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM warnings w JOIN agreements a ON a.id=w.agreement_id WHERE w.status='sent' AND (a.upstream_region=? OR a.downstream_region=?)`, region, region).Scan(&count)
	return count, err
}
