package domain

import "strings"

var allowedStatuses = map[string]map[string]struct{}{"alert": {"open": {}, "investigating": {}, "resolved": {}}, "inspection": {"pending": {}, "running": {}, "completed": {}, "failed": {}}, "manifest": {"created": {}, "in_transit": {}, "accepted": {}, "disposed": {}}, "project": {"planned": {}, "building": {}, "accepted": {}}}

func ValidStatus(object, status string) bool {
	_, ok := allowedStatuses[object][strings.ToLower(status)]
	return ok
}
func Statuses(object string) []string {
	set := allowedStatuses[object]
	result := make([]string, 0, len(set))
	for status := range set {
		result = append(result, status)
	}
	return result
}
func NormalizeStatus(status string) string { return strings.ToLower(strings.TrimSpace(status)) }
