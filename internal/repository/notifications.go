package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"time"
)

type NotificationRepo struct{ DB *sql.DB }

func (r NotificationRepo) Enqueue(ctx context.Context, n domain.Notification) error {
	raw, err := json.Marshal(n)
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, `INSERT INTO audit_events(actor_id,request_id,object_type,object_id,action,result,details_json,created_at) VALUES(NULL,?,?,?,?,?,?,?)`, n.Recipient, n.Channel, n.DedupKey(), "notification", "queued", string(raw), time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("enqueue notification: %w", err)
	}
	return nil
}
func (r NotificationRepo) Count(ctx context.Context, channel, status string) (int, error) {
	var count int
	err := r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_events WHERE object_type=? AND result=? AND action=?`, channel, status, "notification").Scan(&count)
	return count, err
}
func (r NotificationRepo) MarkSent(ctx context.Context, key string) error {
	res, err := r.DB.ExecContext(ctx, `UPDATE audit_events SET result='sent' WHERE object_id=? AND action='notification' AND result='queued'`, key)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r NotificationRepo) Pending(ctx context.Context, limit int) ([]domain.Notification, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT details_json FROM audit_events WHERE action='notification' AND result='queued' ORDER BY id LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.Notification, 0)
	for rows.Next() {
		var raw string
		var n domain.Notification
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(raw), &n); err != nil {
			return nil, err
		}
		result = append(result, n)
	}
	return result, rows.Err()
}
