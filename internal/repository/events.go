package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"time"
)

type Events struct{ DB *sql.DB }

func (r Events) Append(ctx context.Context, event domain.Event) error {
	raw, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, `INSERT INTO audit_events(actor_id,request_id,object_type,object_id,action,result,details_json,created_at) VALUES(?,?,?,?,?,?,?,?)`, event.ActorID, event.ID, event.ObjectType, event.ObjectID, string(event.Type), "pending", string(raw), event.OccurredAt.Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("append event: %w", err)
	}
	return nil
}
func (r Events) MarkDelivered(ctx context.Context, eventID string) error {
	res, err := r.DB.ExecContext(ctx, `UPDATE audit_events SET result='delivered' WHERE object_id=?`, eventID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r Events) Pending(ctx context.Context, limit int) ([]domain.Event, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT request_id,object_type,object_id,actor_id,action,details_json,created_at FROM audit_events WHERE result='pending' ORDER BY id LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.Event, 0)
	for rows.Next() {
		var e domain.Event
		var raw, created string
		if err := rows.Scan(&e.ID, &e.ObjectType, &e.ObjectID, &e.ActorID, &e.Type, &raw, &created); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(raw), &e.Payload)
		e.OccurredAt, _ = time.Parse(time.RFC3339Nano, created)
		result = append(result, e)
	}
	return result, rows.Err()
}
