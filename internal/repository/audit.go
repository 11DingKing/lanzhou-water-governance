package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type Audit struct{ DB *sql.DB }

func (r Audit) Record(ctx context.Context, actorID int64, requestID, objectType, objectID, action, result string, details map[string]any) error {
	raw, err := json.Marshal(details)
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, `INSERT INTO audit_events(actor_id,request_id,object_type,object_id,action,result,details_json,created_at) VALUES(?,?,?,?,?,?,?,?)`, actorID, requestID, objectType, objectID, action, result, string(raw), time.Now().UTC().Format(time.RFC3339Nano))
	return err
}
func (r Audit) RecordTx(ctx context.Context, tx *sql.Tx, actorID int64, requestID, objectType, objectID, action, result string, details map[string]any) error {
	raw, err := json.Marshal(details)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO audit_events(actor_id,request_id,object_type,object_id,action,result,details_json,created_at) VALUES(?,?,?,?,?,?,?,?)`, actorID, requestID, objectType, objectID, action, result, string(raw), time.Now().UTC().Format(time.RFC3339Nano))
	return err
}
func (r Audit) Count(ctx context.Context, objectType, objectID string) (int, error) {
	var count int
	err := r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_events WHERE object_type=? AND object_id=?`, objectType, objectID).Scan(&count)
	return count, err
}
func (r Audit) Last(ctx context.Context, objectType, objectID string) (string, error) {
	var action string
	err := r.DB.QueryRowContext(ctx, `SELECT action FROM audit_events WHERE object_type=? AND object_id=? ORDER BY created_at DESC,id DESC LIMIT 1`, objectType, objectID).Scan(&action)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return action, err
}
