package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
)

type Waste struct{ DB *sql.DB }

func (r Waste) Create(ctx context.Context, m domain.Manifest) (domain.Manifest, error) {
	now := time.Now().UTC()
	res, err := r.DB.ExecContext(ctx, `INSERT INTO manifests(manifest_no,producer_region,carrier_region,facility_region,waste_type,weight_kg,status,created_at) VALUES(?,?,?,?,?,?,?,?)`, m.Number, m.ProducerRegion, m.CarrierRegion, m.FacilityRegion, m.WasteType, m.WeightKG, string(domain.ManifestCreated), now.Format(time.RFC3339Nano))
	if err != nil {
		return m, fmt.Errorf("create manifest: %w", err)
	}
	m.ID, _ = res.LastInsertId()
	m.Status = domain.ManifestCreated
	m.CreatedAt = &now
	m.Version = 1
	return m, nil
}
func (r Waste) Transition(ctx context.Context, id int64, from, to domain.ManifestStatus, version int64) (domain.Manifest, error) {
	if !domain.CanManifestTransition(from, to) {
		return domain.Manifest{}, domain.ErrInvalidState
	}
	now := time.Now().UTC()
	var accepted, disposed any
	if to == domain.ManifestAccepted {
		accepted = now.Format(time.RFC3339Nano)
	}
	if to == domain.ManifestDisposed {
		disposed = now.Format(time.RFC3339Nano)
	}
	res, err := r.DB.ExecContext(ctx, `UPDATE manifests SET status=?,accepted_at=COALESCE(accepted_at,?),disposed_at=?,version=version+1 WHERE id=? AND status=? AND version=?`, to, accepted, disposed, id, from, version)
	if err != nil {
		return domain.Manifest{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Manifest{}, domain.ErrConflict
	}
	return r.Get(ctx, id)
}
func (r Waste) Get(ctx context.Context, id int64) (domain.Manifest, error) {
	var m domain.Manifest
	var created string
	var accepted, disposed sql.NullString
	err := r.DB.QueryRowContext(ctx, `SELECT id,manifest_no,producer_region,carrier_region,facility_region,waste_type,weight_kg,status,created_at,accepted_at,disposed_at,version FROM manifests WHERE id=?`, id).Scan(&m.ID, &m.Number, &m.ProducerRegion, &m.CarrierRegion, &m.FacilityRegion, &m.WasteType, &m.WeightKG, &m.Status, &created, &accepted, &disposed, &m.Version)
	if err == sql.ErrNoRows {
		return m, domain.ErrNotFound
	}
	if err != nil {
		return m, err
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, created)
	m.CreatedAt = &createdAt
	if accepted.Valid {
		t, _ := time.Parse(time.RFC3339Nano, accepted.String)
		m.AcceptedAt = &t
	}
	if disposed.Valid {
		t, _ := time.Parse(time.RFC3339Nano, disposed.String)
		m.DisposedAt = &t
	}
	return m, nil
}
func (r Waste) FindByNumber(ctx context.Context, number string) (domain.Manifest, error) {
	var m domain.Manifest
	var created string
	err := r.DB.QueryRowContext(ctx, `SELECT id,manifest_no,producer_region,carrier_region,facility_region,waste_type,weight_kg,status,created_at,version FROM manifests WHERE manifest_no=?`, number).Scan(&m.ID, &m.Number, &m.ProducerRegion, &m.CarrierRegion, &m.FacilityRegion, &m.WasteType, &m.WeightKG, &m.Status, &created, &m.Version)
	if err == sql.ErrNoRows {
		return m, domain.ErrNotFound
	}
	if err != nil {
		return m, err
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, created)
	m.CreatedAt = &createdAt
	return m, nil
}
