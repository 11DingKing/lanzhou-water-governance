package domain

import "strings"

type TextFilter struct {
	Query  string
	Fields []string
}

func (f TextFilter) Match(values map[string]string) bool {
	query := strings.ToLower(strings.TrimSpace(f.Query))
	if query == "" {
		return true
	}
	for _, field := range f.Fields {
		if strings.Contains(strings.ToLower(values[field]), query) {
			return true
		}
	}
	return false
}
func NormalizeFilter(value string) string {
	return strings.Join(strings.Fields(strings.ToLower(value)), " ")
}
func ParseCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
func UniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; !ok {
			seen[value] = struct{}{}
			result = append(result, value)
		}
	}
	return result
}
