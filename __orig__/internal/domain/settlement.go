package domain

import "math"

type SettlementInput struct {
	BaseCents   int64
	Improvement float64
	Violations  int
	Eligible    bool
}

func SettlementAmount(input SettlementInput) int64 {
	if !input.Eligible || input.BaseCents <= 0 {
		return 0
	}
	factor := 1 + math.Max(0, input.Improvement)
	penalty := 1 - math.Min(0.8, float64(input.Violations)*0.1)
	return int64(math.Round(float64(input.BaseCents) * factor * penalty))
}
func SettlementStatus(amount int64) string {
	if amount <= 0 {
		return "blocked"
	}
	if amount < 10000 {
		return "review"
	}
	return "approved"
}
