package domain

import "time"

type Calendar struct {
	Location         *time.Location
	BusinessDayStart int
	BusinessDayEnd   int
}

func NewCalendar(name string) Calendar {
	loc, err := time.LoadLocation(name)
	if err != nil {
		loc = time.UTC
	}
	return Calendar{Location: loc, BusinessDayStart: 8, BusinessDayEnd: 18}
}
func (c Calendar) StartOfDay(t time.Time) time.Time {
	local := t.In(c.Location)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, c.Location)
}
func (c Calendar) EndOfDay(t time.Time) time.Time {
	return c.StartOfDay(t).Add(24 * time.Hour).Add(-time.Nanosecond)
}
func (c Calendar) IsBusinessHour(t time.Time) bool {
	local := t.In(c.Location)
	hour := local.Hour()
	return hour >= c.BusinessDayStart && hour < c.BusinessDayEnd && local.Weekday() != time.Saturday && local.Weekday() != time.Sunday
}
func (c Calendar) AddBusinessHours(t time.Time, hours int) time.Time {
	if hours <= 0 {
		return t
	}
	current := t
	remaining := hours
	for remaining > 0 {
		current = current.Add(time.Hour)
		if c.IsBusinessHour(current) {
			remaining--
		}
	}
	return current
}
func (c Calendar) Quarter(t time.Time) int { return (int(t.In(c.Location).Month())-1)/3 + 1 }
func (c Calendar) Period(t time.Time) string {
	local := t.In(c.Location)
	return local.Format("2006-01")
}
func (c Calendar) Parse(value string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02T15:04:05", value, c.Location)
}
func (c Calendar) Format(t time.Time) string { return t.In(c.Location).Format(time.RFC3339) }
func (c Calendar) SameDay(left, right time.Time) bool {
	return c.StartOfDay(left).Equal(c.StartOfDay(right))
}
func (c Calendar) DaysBetween(left, right time.Time) int {
	return int(c.StartOfDay(right).Sub(c.StartOfDay(left)) / (24 * time.Hour))
}
