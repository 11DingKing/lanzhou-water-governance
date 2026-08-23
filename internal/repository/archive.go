package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Archive struct{ DB *sql.DB }

func (r Archive) ArchiveAudit(ctx context.Context, cutoff time.Time) (int64, error) {
	res, err := r.DB.ExecContext(ctx, `UPDATE audit_events SET result='archived' WHERE created_at < ? AND result='delivered'`, cutoff.Format(time.RFC3339Nano))
	if err != nil {
		return 0, fmt.Errorf("archive audit: %w", err)
	}
	return res.RowsAffected()
}
func (r Archive) ArchiveCompensations(ctx context.Context, cutoff time.Time) (int64, error) {
	res, err := r.DB.ExecContext(ctx, `UPDATE compensations SET status='archived' WHERE settled_at < ? AND status='settled'`, cutoff.Format(time.RFC3339Nano))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func (r Archive) ArchivedCount(ctx context.Context, table string) (int, error) {
	allowed := map[string]bool{"audit_events": true, "compensations": true}
	if !allowed[table] {
		return 0, fmt.Errorf("unsupported archive table")
	}
	var count int
	err := r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+table+` WHERE result='archived' OR status='archived'`).Scan(&count)
	return count, err
}

var _ *sql.DB
