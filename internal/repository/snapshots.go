package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Snapshot struct {
	ObjectType string
	ObjectID   string
	Payload    map[string]any
	CreatedAt  time.Time
}
type Snapshots struct{ DB *sql.DB }

func (r Snapshots) Save(ctx context.Context, s Snapshot) error {
	raw, err := json.Marshal(s.Payload)
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, `INSERT INTO audit_events(actor_id,request_id,object_type,object_id,action,result,details_json,created_at) VALUES(NULL,?,?,?,?,?,?,?)`, "snapshot", s.ObjectType, s.ObjectID, "snapshot", "ok", string(raw), s.CreatedAt.Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("save snapshot: %w", err)
	}
	return nil
}
func (r Snapshots) LoadLatest(ctx context.Context, objectType, objectID string) (Snapshot, error) {
	var raw, created string
	err := r.DB.QueryRowContext(ctx, `SELECT details_json,created_at FROM audit_events WHERE action='snapshot' AND object_type=? AND object_id=? ORDER BY created_at DESC,id DESC LIMIT 1`, objectType, objectID).Scan(&raw, &created)
	if err == sql.ErrNoRows {
		return Snapshot{}, fmt.Errorf("snapshot not found")
	}
	if err != nil {
		return Snapshot{}, err
	}
	var payload map[string]any
	if err = json.Unmarshal([]byte(raw), &payload); err != nil {
		return Snapshot{}, err
	}
	at, _ := time.Parse(time.RFC3339Nano, created)
	return Snapshot{ObjectType: objectType, ObjectID: objectID, Payload: payload, CreatedAt: at}, nil
}
func (r Snapshots) PurgeBefore(ctx context.Context, deadline time.Time) (int64, error) {
	res, err := r.DB.ExecContext(ctx, `DELETE FROM audit_events WHERE action='snapshot' AND created_at < ?`, deadline.Format(time.RFC3339Nano))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func SnapshotKey(objectType, objectID string) string { return objectType + ":" + objectID }
