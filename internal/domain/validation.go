package domain

import (
	"regexp"
	"strings"
)

var stationCodePattern = regexp.MustCompile(`^[A-Z]{2,6}-[0-9]{1,6}$`)

func ValidStationCode(code string) bool {
	return stationCodePattern.MatchString(strings.TrimSpace(code))
}
func ValidQualityClass(class QualityClass) bool {
	switch class {
	case QualityI, QualityII, QualityIII, QualityIV, QualityV:
		return true
	}
	return false
}
func ValidRole(role Role) bool {
	return role == RoleAdmin || role == RoleInspector || role == RoleRegional
}
func ValidRegion(region string) bool {
	normalized := NormalizeRegion(region)
	return normalized == "LANZHOU" || normalized == "BAIYIN" || normalized == "LINXIA" || normalized == "Wuwei" || normalized == "HAIDONG" || normalized == "XINING"
}
func SanitizeNotes(notes string) string {
	notes = strings.TrimSpace(notes)
	if len(notes) > 2000 {
		return notes[:2000]
	}
	return notes
}
func RequireFields(values ...string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return false
		}
	}
	return true
}
