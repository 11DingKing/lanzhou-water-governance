package service

import (
	"context"
	"fmt"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
)

type Settlement struct {
	Collaboration repository.Collaboration
	Audit         repository.Audit
}

func (s Settlement) Calculate(ctx context.Context, user domain.User, agreement domain.Agreement, input domain.SettlementInput, period, direction string) (domain.Compensation, error) {
	if user.Role != domain.RoleRegional && user.Role != domain.RoleAdmin {
		return domain.Compensation{}, domain.ErrForbidden
	}
	amount := domain.SettlementAmount(input)
	if amount == 0 {
		return domain.Compensation{}, fmt.Errorf("settlement blocked: %s", domain.SettlementStatus(amount))
	}
	c, err := s.Collaboration.CreateCompensationTxDirect(ctx, agreement, period, direction, "formula", amount)
	if err != nil {
		return c, err
	}
	_ = s.Audit.Record(ctx, user.ID, requestID(ctx), "compensation", fmt.Sprint(c.ID), "calculate", "ok", map[string]any{"status": domain.SettlementStatus(amount)})
	return c, nil
}
