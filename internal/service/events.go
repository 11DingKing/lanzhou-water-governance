package service

import (
	"context"
	"fmt"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"time"
)

type EventBus struct{ Repo repository.Events }

func (s EventBus) Publish(ctx context.Context, user domain.User, event domain.Event) error {
	if event.ActorID == 0 {
		event.ActorID = user.ID
	}
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now().UTC()
	}
	if !event.IsValid() {
		return fmt.Errorf("invalid event")
	}
	return s.Repo.Append(ctx, event)
}
func (s EventBus) DeliverPending(ctx context.Context, limit int) error {
	events, err := s.Repo.Pending(ctx, limit)
	if err != nil {
		return err
	}
	for _, event := range events {
		if err := s.Repo.MarkDelivered(ctx, event.ID); err != nil {
			return fmt.Errorf("deliver %s: %w", event.ID, err)
		}
	}
	return nil
}
func (s EventBus) PendingCount(ctx context.Context) (int, error) {
	events, err := s.Repo.Pending(ctx, 1000)
	return len(events), err
}
