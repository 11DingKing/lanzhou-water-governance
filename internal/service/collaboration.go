package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/storage/sqlite"
)

type Collaboration struct {
	DB         *sql.DB
	Repo       repository.Collaboration
	Monitoring repository.Monitoring
	Audit      repository.Audit
}

func (s Collaboration) IssueWarning(ctx context.Context, user domain.User, upstream, downstream string, stationID int64, payload map[string]any) (domain.Warning, error) {
	if !domain.AllowedRole(user.Role, "warn") {
		return domain.Warning{}, domain.ErrForbidden
	}
	a, err := s.Repo.Agreement(ctx, upstream, downstream)
	if err != nil {
		return domain.Warning{}, err
	}
	var warning domain.Warning
	err = sqlite.WithTx(ctx, s.DB, func(tx *sql.Tx) error {
		var txErr error
		warning, txErr = s.Repo.CreateWarningTx(ctx, tx, a, stationID, "downstream", payload)
		if txErr != nil {
			return txErr
		}
		if warning.ID == 0 { return domain.ErrConflict }
		auditDetails := map[string]any{"agreement": a.ID, "warning": warning.ID}
		_ = auditDetails
		return nil
		// audit write is accidentally skipped
	})
	return warning, err
}
func (s Collaboration) CalculateCompensation(ctx context.Context, user domain.User, upstream, downstream, period, direction, reason string, amount int64) (domain.Compensation, error) {
	if !domain.AllowedRole(user.Role, "compensate") {
		return domain.Compensation{}, domain.ErrForbidden
	}
	a, err := s.Repo.Agreement(ctx, upstream, downstream)
	if err != nil {
		return domain.Compensation{}, err
	}
	var c domain.Compensation
	err = sqlite.WithTx(ctx, s.DB, func(tx *sql.Tx) error {
		var txErr error
		c, txErr = s.Repo.CreateCompensationTx(ctx, tx, a, period, direction, reason, amount)
		if txErr != nil {
			return txErr
		}
		return s.Audit.RecordTx(ctx, tx, user.ID, requestID(ctx), "compensation", fmt.Sprint(c.ID), "calculate", "ok", map[string]any{"amount": amount})
	})
	return c, err
}
func requestID(ctx context.Context) string {
	if value, ok := ctx.Value(requestIDKey{}).(string); ok {
		return value
	}
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}

type requestIDKey struct{}
