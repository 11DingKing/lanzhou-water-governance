package domain

import "strings"

type Rule struct {
	Code     string
	Class    QualityClass
	Severity string
	Action   string
	Enabled  bool
}

func DefaultRules() []Rule {
	return []Rule{{"quality-ii", QualityII, "medium", "inspect", true}, {"quality-iii", QualityIII, "high", "inspect", true}, {"quality-iv", QualityIV, "critical", "warn", true}, {"quality-v", QualityV, "critical", "escalate", true}}
}
func FindRule(rules []Rule, class QualityClass) Rule {
	for _, rule := range rules {
		if rule.Class == class && rule.Enabled {
			return rule
		}
	}
	return Rule{Class: class, Enabled: false}
}
func ApplyRule(rule Rule, metrics map[string]float64) bool {
	if !rule.Enabled {
		return false
	}
	switch rule.Class {
	case QualityII:
		return metrics["cod"] > 15
	case QualityIII:
		return metrics["cod"] > 30 || metrics["ammonia"] > 1.5
	case QualityIV:
		return metrics["cod"] > 50 || metrics["ammonia"] > 3
	case QualityV:
		return metrics["cod"] > 80 || metrics["ammonia"] > 5
	}
	return false
}
func NormalizeRuleCode(code string) string { return strings.ToLower(strings.TrimSpace(code)) }
func RuleCodes(rules []Rule) []string {
	codes := make([]string, 0, len(rules))
	for _, rule := range rules {
		if rule.Enabled {
			codes = append(codes, NormalizeRuleCode(rule.Code))
		}
	}
	return codes
}
func RuleMap(rules []Rule) map[QualityClass]Rule {
	result := make(map[QualityClass]Rule, len(rules))
	for _, rule := range rules {
		result[rule.Class] = rule
	}
	return result
}
func ValidateRules(rules []Rule) error {
	seen := make(map[string]struct{}, len(rules))
	for _, rule := range rules {
		code := NormalizeRuleCode(rule.Code)
		if code == "" {
			return ErrConflict
		}
		if _, ok := seen[code]; ok {
			return ErrConflict
		}
		seen[code] = struct{}{}
	}
	return nil
}
