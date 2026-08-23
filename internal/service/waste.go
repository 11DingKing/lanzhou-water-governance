package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
)

type Waste struct {
	Repo  repository.Waste
	Audit repository.Audit
}

func (s Waste) Create(ctx context.Context, user domain.User, m domain.Manifest) (domain.Manifest, error) {
	if !domain.AllowedRole(user.Role, "manifest") {
		return m, domain.ErrForbidden
	}
	if m.Number == "" {
		m.Number = manifestNumber(m)
	}
	if m.WeightKG <= 0 {
		return m, fmt.Errorf("weight must be positive")
	}
	created, err := s.Repo.Create(ctx, m)
	if err != nil {
		return m, err
	}
	_ = s.Audit.Record(ctx, user.ID, requestID(ctx), "manifest", fmt.Sprint(created.ID), "create", "ok", map[string]any{"number": created.Number})
	return created, nil
}
func (s Waste) Advance(ctx context.Context, user domain.User, id int64, from, to domain.ManifestStatus, version int64) (domain.Manifest, error) {
	if !domain.AllowedRole(user.Role, "manifest") {
		return domain.Manifest{}, domain.ErrForbidden
	}
	current, err := s.Repo.Get(ctx, id)
	if err != nil {
		return domain.Manifest{}, err
	}
	if !domain.ManifestRegionsValid(current) {
		return domain.Manifest{}, domain.ErrConflict
	}
	m, err := s.Repo.Transition(ctx, id, from, to, version)
	if err != nil {
		return m, err
	}
	_ = s.Audit.Record(ctx, user.ID, requestID(ctx), "manifest", fmt.Sprint(id), string(to), "ok", nil)
	return m, nil
}
func manifestNumber(m domain.Manifest) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s-%d-%d", m.WasteType, m.WeightKG, time.Now().UnixNano())))
	return "LZ-" + hex.EncodeToString(sum[:])[:12]
}
