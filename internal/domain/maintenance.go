package domain

import "time"

type MaintenanceWindow struct {
	Start, End time.Time
	Reason     string
}

func (w MaintenanceWindow) Active(at time.Time) bool      { return !at.Before(w.Start) && at.Before(w.End) }
func (w MaintenanceWindow) AllowsWrite(at time.Time) bool { return !w.Active(at) }

type RetentionPolicy struct{ AuditDays, SampleDays, SessionDays int }

func (p RetentionPolicy) Cutoffs(now time.Time) (audit, samples, sessions time.Time) {
	return now.AddDate(0, 0, -p.AuditDays), now.AddDate(0, 0, -p.SampleDays), now.AddDate(0, 0, -p.SessionDays)
}
func ValidRetention(p RetentionPolicy) bool {
	return p.AuditDays >= 30 && p.SampleDays >= 365 && p.SessionDays >= 1
}
func NextRetry(at time.Time, attempt int) time.Time {
	if attempt < 1 {
		attempt = 1
	}
	if attempt > 8 {
		attempt = 8
	}
	return at.Add(time.Duration(1<<attempt) * time.Second)
}
func Retryable(status string) bool {
	return status == "temporary" || status == "timeout" || status == "unavailable"
}

func ArchiveCutoff(now time.Time, days int) time.Time { return now.AddDate(0,0,-days) }
