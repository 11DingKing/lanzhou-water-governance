package worker

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
)

type Reconciler struct {
	DB          *sql.DB
	Idempotency repository.Idempotency
	Logger      *slog.Logger
}

func (r Reconciler) Job() Job {
	return func(ctx context.Context) error {
		if err := r.Idempotency.Purge(ctx); err != nil {
			return err
		}
		r.Logger.Info("reconciled expired idempotency keys", "at", time.Now().UTC())
		return nil
	}
}
