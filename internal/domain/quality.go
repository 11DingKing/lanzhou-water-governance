package domain

import (
	"sort"
	"time"
)

type Reading struct {
	At      time.Time
	Class   QualityClass
	Metrics map[string]float64
}
type QualityTrend struct {
	StationID   int64
	Samples     int
	Worst       QualityClass
	AverageRisk float64
	Stable      bool
}

func BuildTrend(stationID int64, readings []Reading) QualityTrend {
	trend := QualityTrend{StationID: stationID, Stable: true}
	if len(readings) == 0 {
		return trend
	}
	ranks := make([]int, 0, len(readings))
	total := 0.0
	worst := QualityI
	for _, reading := range readings {
		rank := qualityRank(reading.Class)
		ranks = append(ranks, rank)
		if rank > qualityRank(worst) {
			worst = reading.Class
		}
		total += RiskScore(reading.Metrics)
	}
	sort.Ints(ranks)
	trend.Samples = len(readings)
	trend.Worst = worst
	trend.AverageRisk = total / float64(len(readings))
	if len(ranks) > 1 && ranks[len(ranks)-1]-ranks[0] > 1 {
		trend.Stable = false
	}
	return trend
}
func IsSafeForDischarge(class QualityClass) bool { return class == QualityI || class == QualityII }
func RequiredInspectionDeadline(class QualityClass, now time.Time) time.Time {
	hours := 72
	switch class {
	case QualityIII:
		hours = 48
	case QualityIV:
		hours = 24
	case QualityV:
		hours = 6
	}
	return now.Add(time.Duration(hours) * time.Hour)
}
func MergeMetrics(previous, current map[string]float64) map[string]float64 {
	result := make(map[string]float64, len(previous)+len(current))
	for k, v := range previous {
		result[k] = v
	}
	for k, v := range current {
		result[k] = (result[k] + v) / 2
	}
	return result
}
func MetricDelta(previous, current map[string]float64) map[string]float64 {
	delta := make(map[string]float64)
	for key, value := range current {
		if old, ok := previous[key]; ok {
			delta[key] = value - old
		} else {
			delta[key] = value
		}
	}
	return delta
}
