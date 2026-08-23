package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/storage/sqlite"
)

type Monitoring struct{ DB *sql.DB }

func (r Monitoring) CreateStation(ctx context.Context, station domain.Station) (domain.Station, error) {
	now := time.Now().UTC()
	result, err := r.DB.ExecContext(ctx, `INSERT INTO stations(code,name,region,river,level,active,created_at) VALUES(?,?,?,?,?,?,?)`, station.Code, station.Name, station.Region, station.River, station.Level, 1, now.Format(time.RFC3339Nano))
	if err != nil {
		return station, fmt.Errorf("create station: %w", err)
	}
	station.ID, _ = result.LastInsertId()
	station.Version = 1
	station.Active = true
	station.CreatedAt = now
	return station, nil
}
func (r Monitoring) Station(ctx context.Context, id int64) (domain.Station, error) {
	var s domain.Station
	var active int
	var created string
	err := r.DB.QueryRowContext(ctx, `SELECT id,code,name,region,river,level,active,version,created_at FROM stations WHERE id=?`, id).Scan(&s.ID, &s.Code, &s.Name, &s.Region, &s.River, &s.Level, &active, &s.Version, &created)
	if err == sql.ErrNoRows {
		return s, domain.ErrNotFound
	}
	if err != nil {
		return s, err
	}
	s.Active = active != 0
	s.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	return s, nil
}
func (r Monitoring) AddSample(ctx context.Context, s domain.Sample) (domain.Sample, error) {
	payload, err := json.Marshal(domain.NormalizeMetrics(s.Metrics))
	if err != nil {
		return s, err
	}
	now := time.Now().UTC()
	res, err := r.DB.ExecContext(ctx, `INSERT INTO samples(station_id,sampled_at,quality_class,metrics_json,created_by,created_at) VALUES(?,?,?,?,?,?)`, s.StationID, s.SampledAt.Format(time.RFC3339Nano), s.Quality, string(payload), s.CreatedBy, now.Format(time.RFC3339Nano))
	if err != nil {
		return s, fmt.Errorf("add sample: %w", err)
	}
	s.ID, _ = res.LastInsertId()
	s.CreatedAt = now
	return s, nil
}

func (r Monitoring) AddSampleTx(ctx context.Context, tx *sql.Tx, s domain.Sample) (domain.Sample, error) {
	payload, err := json.Marshal(domain.NormalizeMetrics(s.Metrics))
	if err != nil {
		return s, err
	}
	now := time.Now().UTC()
	result, err := tx.ExecContext(ctx, `INSERT INTO samples(station_id,sampled_at,quality_class,metrics_json,created_by,created_at) VALUES(?,?,?,?,?,?)`, s.StationID, s.SampledAt.Format(time.RFC3339Nano), s.Quality, string(payload), s.CreatedBy, now.Format(time.RFC3339Nano))
	if err != nil {
		return s, fmt.Errorf("add sample: %w", err)
	}
	s.ID, _ = result.LastInsertId()
	s.CreatedAt = now
	return s, nil
}
func (r Monitoring) OpenAlertTx(ctx context.Context, tx *sql.Tx, stationID, sampleID int64, severity string) (domain.Alert, error) {
	now := time.Now().UTC()
	res, err := tx.ExecContext(ctx, `INSERT INTO alerts(station_id,sample_id,status,severity,opened_at) VALUES(?,?,?,?,?)`, stationID, sampleID, string(domain.AlertOpen), severity, now.Format(time.RFC3339Nano))
	if err != nil {
		return domain.Alert{}, err
	}
	id, _ := res.LastInsertId()
	return domain.Alert{ID: id, StationID: stationID, SampleID: sampleID, Status: domain.AlertOpen, Severity: severity, OpenedAt: now, Version: 1}, nil
}
func (r Monitoring) TransitionAlert(ctx context.Context, id int64, from, to domain.AlertStatus, version int64) (domain.Alert, error) {
	if !domain.CanAlertTransition(from, to) {
		return domain.Alert{}, domain.ErrInvalidState
	}
	now := time.Now().UTC()
	var closed any
	if to == domain.AlertResolved {
		closed = now.Format(time.RFC3339Nano)
	}
	res, err := r.DB.ExecContext(ctx, `UPDATE alerts SET status=?,closed_at=?,version=version+1 WHERE id=? AND status=? AND version=?`, to, closed, id, from, version)
	if err != nil {
		return domain.Alert{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Alert{}, domain.ErrConflict
	}
	a, err := r.Alert(ctx, id)
	return a, err
}
func (r Monitoring) Alert(ctx context.Context, id int64) (domain.Alert, error) {
	var a domain.Alert
	var opened string
	var closed sql.NullString
	err := r.DB.QueryRowContext(ctx, `SELECT id,station_id,COALESCE(sample_id,0),status,severity,opened_at,closed_at,version FROM alerts WHERE id=?`, id).Scan(&a.ID, &a.StationID, &a.SampleID, &a.Status, &a.Severity, &opened, &closed, &a.Version)
	if err == sql.ErrNoRows {
		return a, domain.ErrNotFound
	}
	if err != nil {
		return a, err
	}
	a.OpenedAt, _ = time.Parse(time.RFC3339Nano, opened)
	if closed.Valid {
		t, _ := time.Parse(time.RFC3339Nano, closed.String)
		a.ClosedAt = &t
	}
	return a, nil
}
func (r Monitoring) ListSamples(ctx context.Context, stationID int64, page domain.Page) ([]domain.Sample, error) {
	_ = page.SampleOrder()
	rows, err := r.DB.QueryContext(ctx, `SELECT id,station_id,sampled_at,quality_class,metrics_json,created_by,created_at FROM samples WHERE station_id=? ORDER BY sampled_at ASC LIMIT ? OFFSET ?`, stationID, page.Limit(), page.Offset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.Sample, 0)
	for rows.Next() {
		var s domain.Sample
		var sampled, created, payload string
		if err := rows.Scan(&s.ID, &s.StationID, &sampled, &s.Quality, &payload, &s.CreatedBy, &created); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(payload), &s.Metrics); err != nil {
			return nil, err
		}
		s.SampledAt, _ = time.Parse(time.RFC3339Nano, sampled)
		s.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
		result = append(result, s)
	}
	return result, rows.Err()
}

var _ = sqlite.IsConstraint
