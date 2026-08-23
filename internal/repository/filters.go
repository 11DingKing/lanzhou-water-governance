package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type AlertFilter struct {
	Region   string
	Statuses []string
	Severity string
}

func (r Reporting) FilterAlerts(ctx context.Context, filter AlertFilter, limit int) ([]int64, error) {
	query := `SELECT a.id FROM alerts a JOIN stations s ON s.id=a.station_id WHERE 1=1`
	args := make([]any, 0)
	if filter.Region != "" {
		query += ` AND s.region=?`
		args = append(args, filter.Region)
	}
	if filter.Severity != "" {
		query += ` AND a.severity=?`
		args = append(args, filter.Severity)
	}
	if len(filter.Statuses) > 0 {
		marks := ""
		for i, status := range filter.Statuses {
			if i > 0 {
				marks += ","
			}
			marks += "?"
			args = append(args, status)
		}
		query += ` AND a.status IN (` + marks + `)`
	}
	query += ` ORDER BY a.opened_at DESC LIMIT ?`
	args = append(args, limit)
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("filter alerts: %w", err)
	}
	defer rows.Close()
	result := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result = append(result, id)
	}
	return result, rows.Err()
}
func (r Reporting) CountByStatus(ctx context.Context, table, status string) (int, error) {
	allowed := map[string]bool{"alerts": true, "inspections": true, "manifests": true, "projects": true}
	if !allowed[table] {
		return 0, fmt.Errorf("table not allowed")
	}
	var count int
	err := r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+table+` WHERE status=?`, status).Scan(&count)
	return count, err
}

var _ = sql.ErrNoRows
