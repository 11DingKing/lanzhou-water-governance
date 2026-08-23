package domain

import "time"

type EventType string

const (
	EventSampleRecorded  EventType = "sample.recorded"
	EventAlertOpened     EventType = "alert.opened"
	EventWarningSent     EventType = "warning.sent"
	EventManifestMoved   EventType = "manifest.moved"
	EventProjectAccepted EventType = "project.accepted"
)

type Event struct {
	ID         string
	Type       EventType
	ObjectType string
	ObjectID   string
	ActorID    int64
	OccurredAt time.Time
	Payload    map[string]any
}

func (e Event) IsValid() bool {
	return e.ID != "" && e.Type != "" && e.ObjectType != "" && e.ObjectID != "" && !e.OccurredAt.IsZero()
}
func EventNeedsRetry(event Event) bool {
	switch event.Type {
	case EventWarningSent, EventManifestMoved:
		return true
	default:
		return false
	}
}
func EventKey(event Event) string {
	return string(event.Type) + ":" + event.ObjectType + ":" + event.ObjectID
}

func SnapshotEvent(event Event) Event { return event }
