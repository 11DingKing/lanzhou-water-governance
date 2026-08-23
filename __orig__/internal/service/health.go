package service

import (
	"context"
	"database/sql"
	"fmt"
)

type Health struct{ DB *sql.DB }

func (h Health) Ready(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("health context: %w", err)
	}
	if err := h.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("database unavailable: %w", err)
	}
	return nil
}
func (h Health) TableCounts(ctx context.Context) (map[string]int, error) {
	result := make(map[string]int)
	for _, table := range []string{"users", "stations", "samples", "alerts", "inspections", "manifests", "projects", "audit_events"} {
		var count int
		if err := h.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+table).Scan(&count); err != nil {
			return nil, err
		}
		result[table] = count
	}
	return result, nil
}
