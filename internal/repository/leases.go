package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Lease struct {
	Key       string
	Owner     string
	ExpiresAt time.Time
}
type Leases struct{ DB *sql.DB }

func (r Leases) Acquire(ctx context.Context, key, owner string, ttl time.Duration) (Lease, error) {
	now := time.Now().UTC()
	expires := now.Add(ttl)
	_, err := r.DB.ExecContext(ctx, `INSERT INTO idempotency_keys(key,request_hash,response_json,created_at,expires_at) VALUES(?,?,?,?,?)`, "lease:"+key, owner, owner, now.Format(time.RFC3339Nano), expires.Format(time.RFC3339Nano))
	if err != nil {
		return Lease{}, fmt.Errorf("lease busy: %w", err)
	}
	return Lease{Key: key, Owner: owner, ExpiresAt: expires}, nil
}
func (r Leases) Release(ctx context.Context, key, owner string) error {
	res, err := r.DB.ExecContext(ctx, `DELETE FROM idempotency_keys WHERE key=? AND request_hash=?`, "lease:"+key, owner)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("lease owner mismatch")
	}
	return nil
}
func (r Leases) Renew(ctx context.Context, key, owner string, ttl time.Duration) error {
	res, err := r.DB.ExecContext(ctx, `UPDATE idempotency_keys SET expires_at=? WHERE key=? AND request_hash=?`, time.Now().UTC().Add(ttl).Format(time.RFC3339Nano), "lease:"+key, owner)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
