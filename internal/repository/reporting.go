package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
)

type ReportRow struct {
	StationID   int64
	StationCode string
	Region      string
	SampleCount int
	OpenAlerts  int
	LastSample  *time.Time
}
type Reporting struct{ DB *sql.DB }

func (r Reporting) StationSummary(ctx context.Context, region string, page domain.Page) ([]ReportRow, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT s.id,s.code,s.region,COUNT(DISTINCT sm.id),COUNT(DISTINCT CASE WHEN a.status IN ('open','investigating') THEN a.id END),MAX(sm.sampled_at) FROM stations s LEFT JOIN samples sm ON sm.station_id=s.id LEFT JOIN alerts a ON a.station_id=s.id WHERE (?='' OR s.region=?) GROUP BY s.id,s.code,s.region ORDER BY s.code LIMIT ? OFFSET ?`, region, region, page.Limit(), page.Offset())
	if err != nil {
		return nil, fmt.Errorf("station summary: %w", err)
	}
	defer rows.Close()
	result := make([]ReportRow, 0)
	for rows.Next() {
		var row ReportRow
		var last sql.NullString
		if err := rows.Scan(&row.StationID, &row.StationCode, &row.Region, &row.SampleCount, &row.OpenAlerts, &last); err != nil {
			return nil, err
		}
		if last.Valid {
			parsed, _ := time.Parse(time.RFC3339Nano, last.String)
			row.LastSample = &parsed
		}
		result = append(result, row)
	}
	return result, rows.Err()
}
func (r Reporting) OpenInspectionCount(ctx context.Context, ownerID int64) (int, error) {
	var count int
	err := r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM inspections WHERE owner_id=? AND status IN ('pending','running')`, ownerID).Scan(&count)
	return count, err
}
func (r Reporting) CompensationTotal(ctx context.Context, agreementID int64, period string) (int64, error) {
	var total sql.NullInt64
	err := r.DB.QueryRowContext(ctx, `SELECT SUM(amount_cents) FROM compensations WHERE agreement_id=? AND period=?`, agreementID, period).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total.Int64, nil
}
func (r Reporting) WasteByRegion(ctx context.Context, region string) (int64, error) {
	var total sql.NullInt64
	err := r.DB.QueryRowContext(ctx, `SELECT SUM(weight_kg) FROM manifests WHERE producer_region=? OR facility_region=?`, region, region).Scan(&total)
	return total.Int64, err
}
func (r Reporting) ProjectProgress(ctx context.Context, projectID int64) (float64, error) {
	var completed, total int
	err := r.DB.QueryRowContext(ctx, `SELECT SUM(CASE WHEN status='completed' THEN 1 ELSE 0 END),COUNT(*) FROM milestones WHERE project_id=?`, projectID).Scan(&completed, &total)
	if err != nil {
		return 0, err
	}
	if total == 0 {
		return 0, nil
	}
	return float64(completed) / float64(total), nil
}
