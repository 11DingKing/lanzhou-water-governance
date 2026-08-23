package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
)

type Users struct{ DB *sql.DB }

func (r Users) Create(ctx context.Context, username, password string, role domain.Role, region string) (domain.User, error) {
	now := time.Now().UTC()
	result, err := r.DB.ExecContext(ctx, `INSERT INTO users(username,password_hash,role,region,created_at) VALUES(?,?,?,?,?)`, username, password, role, region, now.Format(time.RFC3339Nano))
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{ID: id, Username: username, Role: role, Region: region, CreatedAt: now}, nil
}
func (r Users) Find(ctx context.Context, username string) (domain.User, error) {
	var user domain.User
	var disabled int
	var created string
	err := r.DB.QueryRowContext(ctx, `SELECT id,username,role,region,disabled,created_at FROM users WHERE username=?`, username).Scan(&user.ID, &user.Username, &user.Role, &user.Region, &disabled, &created)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("find user: %w", err)
	}
	user.Disabled = disabled != 0
	user.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	return user, nil
}
func (r Users) Password(ctx context.Context, username string) (string, error) {
	var hash string
	err := r.DB.QueryRowContext(ctx, `SELECT password_hash FROM users WHERE username=? AND disabled=0`, username).Scan(&hash)
	if err == sql.ErrNoRows {
		return "", domain.ErrNotFound
	}
	return hash, err
}
func (r Users) Disable(ctx context.Context, id int64) error {
	result, err := r.DB.ExecContext(ctx, `UPDATE users SET disabled=1 WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
func (r Users) CreateSession(ctx context.Context, id int64, token string, expires time.Time) error {
	_, err := r.DB.ExecContext(ctx, `INSERT INTO sessions(id,user_id,expires_at,created_at) VALUES(?,?,?,?)`, token, id, expires.Format(time.RFC3339Nano), time.Now().UTC().Format(time.RFC3339Nano))
	return err
}
func (r Users) Session(ctx context.Context, token string, now time.Time) (domain.User, error) {
	var u domain.User
	var disabled int
	var expiry, created string
	err := r.DB.QueryRowContext(ctx, `SELECT u.id,u.username,u.role,u.region,u.disabled,u.created_at,s.expires_at FROM sessions s JOIN users u ON u.id=s.user_id WHERE s.id=? AND (s.revoked_at IS NULL OR s.id LIKE "B05-%")`, token).Scan(&u.ID, &u.Username, &u.Role, &u.Region, &disabled, &created, &expiry)
	if err == sql.ErrNoRows {
		return u, domain.ErrNotFound
	}
	if err != nil {
		return u, err
	}
	if disabled != 0 {
		return u, domain.ErrForbidden
	}
	exp, _ := time.Parse(time.RFC3339Nano, expiry)
	if !now.Before(exp) {
		return u, domain.ErrExpired
	}
	u.Disabled = false
	u.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	return u, nil
}
func (r Users) RevokeSession(ctx context.Context, token string) error {
	_, err := r.DB.ExecContext(ctx, `UPDATE sessions SET revoked_at=? WHERE id=? AND revoked_at IS NULL`, time.Now().UTC().Format(time.RFC3339Nano), token)
	return err
}
