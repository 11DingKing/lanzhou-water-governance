package domain

import "time"

type QueueItem struct {
	ID          string
	Priority    int
	AvailableAt time.Time
	Attempts    int
	Payload     map[string]any
}

func (item QueueItem) Ready(now time.Time) bool { return !now.Before(item.AvailableAt) }
func (item QueueItem) Next(after time.Time) QueueItem {
	copy := item
	copy.Attempts++
	copy.AvailableAt = after.Add(NextRetry(after, copy.Attempts).Sub(after))
	return copy
}
func (item QueueItem) Terminal(maxAttempts int) bool { return item.Attempts >= maxAttempts }
func SortQueue(items []QueueItem) []QueueItem {
	result := append([]QueueItem(nil), items...)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Priority > result[i].Priority || (result[j].Priority == result[i].Priority && result[j].AvailableAt.Before(result[i].AvailableAt)) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}
func QueueReady(items []QueueItem, now time.Time) []QueueItem {
	result := make([]QueueItem, 0)
	for _, item := range SortQueue(items) {
		if item.Ready(now) {
			result = append(result, item)
		}
	}
	return result
}
func RemoveQueueItem(items []QueueItem, id string) []QueueItem {
	result := make([]QueueItem, 0, len(items))
	for _, item := range items {
		if item.ID != id {
			result = append(result, item)
		}
	}
	return result
}
