package domain

import "time"

func AllowedRole(role Role, action string) bool {
	switch role {
	case RoleAdmin:
		return true
	case RoleInspector:
		return action == "sample" || action == "inspect" || action == "remediate"
	case RoleRegional:
		return action == "warn" || action == "compensate" || action == "manifest"
	}
	return false
}
func WithinWindow(now, due time.Time) bool { return !now.After(due) }
func QualityNeedsInvestigation(class QualityClass, threshold QualityClass) bool {
	return qualityRank(class) > qualityRank(threshold)
}
func qualityRank(c QualityClass) int {
	switch c {
	case QualityI:
		return 1
	case QualityII:
		return 2
	case QualityIII:
		return 3
	case QualityIV:
		return 4
	case QualityV:
		return 5
	}
	return 99
}
