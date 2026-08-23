package domain

import "strings"

type AuditFilter struct{ ObjectType, ObjectID, Action, Result string }

func (f AuditFilter) Match(objectType, objectID, action, result string) bool {
	if f.ObjectType != "" && f.ObjectType != objectType {
		return false
	}
	if f.ObjectID != "" && f.ObjectID != objectID {
		return false
	}
	if f.Action != "" && f.Action != action {
		return false
	}
	if f.Result != "" && f.Result != result {
		return false
	}
	return true
}
func Redact(value string) string {
	value = strings.TrimSpace(value)
	if len(value) <= 4 {
		return "****"
	}
	return value[:2] + "****" + value[len(value)-2:]
}
func NormalizeAction(action string) string { return strings.ToLower(strings.TrimSpace(action)) }
