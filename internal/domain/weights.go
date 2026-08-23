package domain

func WeightedRisk(metrics map[string]float64, weights map[string]float64) float64 {
	total := 0.0
	for metric, value := range metrics {
		weight := weights[metric]
		if weight == 0 {
			weight = 1
		}
		if value > 0 {
			total += value * weight
		}
	}
	return total
}
func NormalizeWeights(weights map[string]float64) map[string]float64 {
	result := make(map[string]float64, len(weights))
	total := 0.0
	for _, value := range weights {
		if value > 0 {
			total += value
		}
	}
	if total == 0 {
		return result
	}
	for key, value := range weights {
		if value > 0 {
			result[key] = value / total
		}
	}
	return result
}
func TopRisk(metrics map[string]float64, limit int) []string {
	keys := make([]string, 0, len(metrics))
	for key := range metrics {
		keys = append(keys, key)
	}
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if metrics[keys[j]] > metrics[keys[i]] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	if limit > len(keys) {
		limit = len(keys)
	}
	return keys[:limit]
}
