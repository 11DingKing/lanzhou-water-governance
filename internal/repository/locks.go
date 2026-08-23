package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type LockRepo struct{ DB *sql.DB }

func (r LockRepo) Acquire(ctx context.Context, key, owner string, ttl time.Duration) error {
	_, err := r.DB.ExecContext(ctx, `INSERT INTO idempotency_keys(key,request_hash,response_json,created_at,expires_at) VALUES(?,?,?,?,?)`, "lock:"+key, owner, owner, time.Now().UTC().Format(time.RFC3339Nano), time.Now().UTC().Add(ttl).Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("acquire lock: %w", err)
	}
	return nil
}
func (r LockRepo) Release(ctx context.Context, key, owner string) error {
	res, err := r.DB.ExecContext(ctx, `DELETE FROM idempotency_keys WHERE key=? AND request_hash=?`, "lock:"+key, owner)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r LockRepo) Held(ctx context.Context, key string) (bool, error) {
	var count int
	err := r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM idempotency_keys WHERE key=? AND expires_at>?`, "lock:"+key, time.Now().UTC().Format(time.RFC3339Nano)).Scan(&count)
	return count > 0, err
}
