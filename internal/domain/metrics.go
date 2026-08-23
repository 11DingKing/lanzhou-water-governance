package domain

import "math"

func NormalizeMetrics(input map[string]float64) map[string]float64 {
	result := make(map[string]float64, len(input))
	for key, value := range input {
		if !math.IsNaN(value) && !math.IsInf(value, 0) {
			result[key] = value
		}
	}
	return result
}
func RiskScore(metrics map[string]float64) float64 {
	score := 0.0
	for key, value := range metrics {
		weight := 1.0
		if key == "ammonia" {
			weight = 1.5
		}
		if key == "cod" {
			weight = 1.2
		}
		score += math.Max(0, value) * weight
	}
	return score
}
