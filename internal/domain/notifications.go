package domain

import "strings"

type Notification struct {
	Recipient string
	Channel   string
	Subject   string
	Body      string
	Urgent    bool
}

func (n Notification) Valid() bool {
	return strings.TrimSpace(n.Recipient) != "" && strings.TrimSpace(n.Body) != "" && (n.Channel == "sms" || n.Channel == "email" || n.Channel == "webhook")
}
func (n Notification) DedupKey() string { return n.Recipient + ":" + n.Channel + ":" + n.Subject }
func TruncateBody(body string, max int) string {
	if max < 1 {
		return ""
	}
	runes := []rune(body)
	if len(runes) <= max {
		return body
	}
	return string(runes[:max])
}
func BatchNotifications(items []Notification) map[string][]Notification {
	result := make(map[string][]Notification)
	for _, item := range items {
		result[item.Channel] = append(result[item.Channel], item)
	}
	return result
}
