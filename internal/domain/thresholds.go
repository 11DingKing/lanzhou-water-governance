package domain

type Threshold struct {
	Metric    string
	Class     QualityClass
	Min, Max  float64
	Inclusive bool
}

func (t Threshold) Matches(value float64) bool {
	if t.Inclusive {
		return value >= t.Min && value <= t.Max
	}
	return value > t.Min && value < t.Max
}
func DefaultThresholds() []Threshold {
	return []Threshold{{"cod", QualityII, 0, 15, true}, {"cod", QualityIII, 15, 30, false}, {"cod", QualityIV, 30, 50, false}, {"cod", QualityV, 50, 100000, true}, {"ammonia", QualityII, 0, 0.5, true}, {"ammonia", QualityIII, 0.5, 1.5, false}, {"ammonia", QualityIV, 1.5, 3, false}, {"ammonia", QualityV, 3, 100000, true}}
}
func ClassifyMetric(metric string, value float64) QualityClass {
	for _, threshold := range DefaultThresholds() {
		if threshold.Metric == metric && threshold.Matches(value) {
			return threshold.Class
		}
	}
	return QualityV
}
func WorstClass(classes ...QualityClass) QualityClass {
	worst := QualityI
	for _, class := range classes {
		if qualityRank(class) > qualityRank(worst) {
			worst = class
		}
	}
	return worst
}
func ClassifyMetrics(metrics map[string]float64) QualityClass {
	classes := make([]QualityClass, 0, len(metrics))
	for metric, value := range metrics {
		classes = append(classes, ClassifyMetric(metric, value))
	}
	return WorstClass(classes...)
}
