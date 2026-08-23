package repository

import (
	"context"
	"database/sql"
	"time"
)

type PolicyRepository struct{ DB *sql.DB }

func (r PolicyRepository) ActiveAgreements(ctx context.Context) (int, error) {
	var count int
	err := r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM agreements WHERE active=1`).Scan(&count)
	return count, err
}
func (r PolicyRepository) DeactivateAgreement(ctx context.Context, id int64) error {
	res, err := r.DB.ExecContext(ctx, `UPDATE agreements SET active=0 WHERE id=? AND active=1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r PolicyRepository) ReactivateAgreement(ctx context.Context, id int64) error {
	res, err := r.DB.ExecContext(ctx, `UPDATE agreements SET active=1 WHERE id=? AND active=0`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r PolicyRepository) AgreementsChangedSince(ctx context.Context, at time.Time) (int, error) {
	var count int
	err := r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM agreements WHERE created_at >= ?`, at.Format(time.RFC3339Nano)).Scan(&count)
	return count, err
}
