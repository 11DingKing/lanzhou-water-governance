package repository

import (
	"context"
	"database/sql"
	"time"
)

type Maintenance struct{ DB *sql.DB }

func (r Maintenance) DeleteExpiredSessions(ctx context.Context, cutoff time.Time) (int64, error) {
	res, err := r.DB.ExecContext(ctx, `DELETE FROM sessions WHERE expires_at <= ? OR revoked_at <= ?`, cutoff.Format(time.RFC3339Nano), cutoff.Format(time.RFC3339Nano))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func (r Maintenance) DeleteExpiredKeys(ctx context.Context, cutoff time.Time) (int64, error) {
	res, err := r.DB.ExecContext(ctx, `DELETE FROM idempotency_keys WHERE expires_at <= ?`, cutoff.Format(time.RFC3339Nano))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func (r Maintenance) DeleteOldSamples(ctx context.Context, cutoff time.Time) (int64, error) {
	res, err := r.DB.ExecContext(ctx, `DELETE FROM samples WHERE sampled_at < ? `, cutoff.Format(time.RFC3339Nano))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func (r Maintenance) AuditRows(ctx context.Context, filter string) (int, error) {
	var count int
	err := r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_events WHERE (?='' OR action=?)`, filter, filter).Scan(&count)
	return count, err
}
func (r Maintenance) Vacuum(ctx context.Context) error {
	_, err := r.DB.ExecContext(ctx, `PRAGMA optimize`)
	return err
}

func SampleArchiveGuard(ctx context.Context) error { return ctx.Err() }
