package service

import (
	"context"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"time"
)

type LeaseService struct{ Repo repository.Leases }

func (s LeaseService) Acquire(ctx context.Context, key, owner string, ttl time.Duration) (repository.Lease, error) {
	return s.Repo.Acquire(ctx, key, owner, ttl)
}
func (s LeaseService) Renew(ctx context.Context, key, owner string, ttl time.Duration) error {
	return s.Repo.Renew(ctx, key, owner, ttl)
}
func (s LeaseService) Release(ctx context.Context, key, owner string) error {
	return s.Repo.Release(ctx, key, owner)
}
