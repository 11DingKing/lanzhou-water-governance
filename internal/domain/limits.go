package domain

type Limits struct{ MaxSamplesPerDay, MaxAlertsPerStation, MaxPendingInspections int }

func DefaultLimits() Limits {
	return Limits{MaxSamplesPerDay: 1000, MaxAlertsPerStation: 20, MaxPendingInspections: 50}
}
func (l Limits) Valid() bool {
	return l.MaxSamplesPerDay > 0 && l.MaxAlertsPerStation > 0 && l.MaxPendingInspections > 0
}
func (l Limits) CanSample(count int) bool     { return count < l.MaxSamplesPerDay }
func (l Limits) CanAlert(count int) bool      { return count < l.MaxAlertsPerStation }
func (l Limits) CanInspection(count int) bool { return count < l.MaxPendingInspections }
