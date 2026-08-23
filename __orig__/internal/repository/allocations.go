package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CapacityRow struct {
	Region                   string
	Capacity, Used, Reserved int64
}
type CapacityRepo struct{ DB *sql.DB }

func (r CapacityRepo) Ensure(ctx context.Context, region string, capacity int64) error {
	_, err := r.DB.ExecContext(ctx, `INSERT INTO audit_events(actor_id,request_id,object_type,object_id,action,result,details_json,created_at) VALUES(NULL,?,?,?,?,?,?,?)`, region, "capacity", region, "capacity", fmt.Sprint(capacity), `{}`, time.Now().UTC().Format(time.RFC3339Nano))
	return err
}
func (r CapacityRepo) Adjust(ctx context.Context, region string, delta int64) error {
	res, err := r.DB.ExecContext(ctx, `UPDATE audit_events SET result=CAST(CAST(result AS INTEGER)+? AS TEXT) WHERE action='capacity' AND object_id=?`, delta, region)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r CapacityRepo) Read(ctx context.Context, region string) (CapacityRow, error) {
	var value string
	err := r.DB.QueryRowContext(ctx, `SELECT result FROM audit_events WHERE action='capacity' AND object_id=?`, region).Scan(&value)
	if err != nil {
		return CapacityRow{}, err
	}
	var capacity int64
	_, err = fmt.Sscan(value, &capacity)
	return CapacityRow{Region: region, Capacity: capacity}, err
}
func (r CapacityRepo) Delete(ctx context.Context, region string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM audit_events WHERE action='capacity' AND object_id=?`, region)
	return err
}
