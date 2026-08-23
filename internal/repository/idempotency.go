package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
)

type Idempotency struct{ DB *sql.DB }

func (r Idempotency) Load(ctx context.Context, key, hash string) (string, error) {
	var savedHash, response, expires string
	err := r.DB.QueryRowContext(ctx, `SELECT request_hash,response_json,expires_at FROM idempotency_keys WHERE key=?`, key).Scan(&savedHash, &response, &expires)
	if err == sql.ErrNoRows {
		return "", domain.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	expiry, _ := time.Parse(time.RFC3339Nano, expires)
	if !time.Now().UTC().Before(expiry) {
		_, _ = r.DB.ExecContext(ctx, `DELETE FROM idempotency_keys WHERE key=?`, key)
		return "", domain.ErrNotFound
	}
	if savedHash != hash {
		return "", domain.ErrDuplicate
	}
	return response, nil
}
func (r Idempotency) Save(ctx context.Context, key, hash, response string, ttl time.Duration) error {
	now := time.Now().UTC()
	_, err := r.DB.ExecContext(ctx, `INSERT INTO idempotency_keys(key,request_hash,response_json,created_at,expires_at) VALUES(?,?,?,?,?)`, key, hash, response, now.Format(time.RFC3339Nano), now.Add(ttl).Format(time.RFC3339Nano))
	return err
}
func (r Idempotency) Purge(ctx context.Context) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM idempotency_keys WHERE expires_at <= ?`, time.Now().UTC().Format(time.RFC3339Nano))
	return err
}
