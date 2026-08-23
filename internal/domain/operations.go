package domain

import "time"

type Operation struct {
	ID                    string
	Name                  string
	Status                string
	StartedAt, FinishedAt *time.Time
	OwnerID               int64
	Error                 string
}

func (o Operation) Running() bool { return o.Status == "running" }
func (o *Operation) Start(now time.Time) error {
	if o.Status != "queued" {
		return ErrConflict
	}
	o.Status = "running"
	o.StartedAt = &now
	return nil
}
func (o *Operation) Finish(now time.Time) error {
	if o.Status != "running" {
		return ErrConflict
	}
	o.Status = "completed"
	o.FinishedAt = &now
	return nil
}
func (o *Operation) Fail(now time.Time, err string) error {
	if o.Status != "running" {
		return ErrConflict
	}
	o.Status = "failed"
	o.FinishedAt = &now
	o.Error = err
	return nil
}
func (o Operation) Duration() time.Duration {
	if o.StartedAt == nil || o.FinishedAt == nil {
		return 0
	}
	return o.FinishedAt.Sub(*o.StartedAt)
}
