package domain

type ChecklistItem struct {
	Code     string
	Required bool
	Complete bool
	Evidence string
}

func ChecklistComplete(items []ChecklistItem) bool {
	for _, item := range items {
		if item.Required && !item.Complete {
			return false
		}
	}
	return true
}
func MissingChecklist(items []ChecklistItem) []string {
	missing := make([]string, 0)
	for _, item := range items {
		if item.Required && !item.Complete {
			missing = append(missing, item.Code)
		}
	}
	return missing
}
func CompleteChecklist(items []ChecklistItem, code, evidence string) []ChecklistItem {
	result := append([]ChecklistItem(nil), items...)
	for i := range result {
		if result[i].Code == code {
			result[i].Complete = true
			result[i].Evidence = evidence
		}
	}
	return result
}
func ValidateChecklist(items []ChecklistItem) error {
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		if item.Code == "" {
			return ErrConflict
		}
		if _, ok := seen[item.Code]; ok {
			return ErrConflict
		}
		seen[item.Code] = struct{}{}
	}
	return nil
}
func RequiredChecklist(items []ChecklistItem) int {
	count := 0
	for _, item := range items {
		if item.Required {
			count++
		}
	}
	return count
}
