package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"time"
)

type ExportRepo struct{ DB *sql.DB }

func (r ExportRepo) Save(ctx context.Context, job domain.ExportJob) error {
	_, err := r.DB.ExecContext(ctx, `INSERT INTO audit_events(actor_id,request_id,object_type,object_id,action,result,details_json,created_at) VALUES(NULL,?,?,?,?,?,?,?)`, job.ID, "export", job.ID, "export", job.Status, fmt.Sprintf(`{"format":%q,"rows":%d}`, job.Format, job.Rows), time.Now().UTC().Format(time.RFC3339Nano))
	return err
}
func (r ExportRepo) SetStatus(ctx context.Context, id, status string) error {
	res, err := r.DB.ExecContext(ctx, `UPDATE audit_events SET result=? WHERE action='export' AND object_id=?`, status, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r ExportRepo) Status(ctx context.Context, id string) (string, error) {
	var status string
	err := r.DB.QueryRowContext(ctx, `SELECT result FROM audit_events WHERE action='export' AND object_id=? ORDER BY id DESC LIMIT 1`, id).Scan(&status)
	return status, err
}
func (r ExportRepo) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM audit_events WHERE action='export' AND object_id=?`, id)
	return err
}
