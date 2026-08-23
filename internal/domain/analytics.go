package domain

import "math"

type Histogram struct {
	Buckets []float64
	Counts  []int
}

func NewHistogram(buckets []float64) Histogram {
	return Histogram{Buckets: append([]float64(nil), buckets...), Counts: make([]int, len(buckets)+1)}
}
func (h *Histogram) Add(value float64) {
	index := len(h.Buckets)
	for i, bucket := range h.Buckets {
		if value < bucket {
			index = i
			break
		}
	}
	h.Counts[index]++
}
func (h Histogram) Total() int {
	total := 0
	for _, count := range h.Counts {
		total += count
	}
	return total
}
func (h Histogram) Peak() int {
	peak := 0
	for _, count := range h.Counts {
		if count > peak {
			peak = count
		}
	}
	return peak
}
func (h Histogram) Mean() float64 {
	total := 0.0
	count := 0
	for i, value := range h.Buckets {
		total += value * float64(h.Counts[i])
		count += h.Counts[i]
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
func Percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	copyValues := append([]float64(nil), values...)
	for i := 0; i < len(copyValues); i++ {
		for j := i + 1; j < len(copyValues); j++ {
			if copyValues[j] < copyValues[i] {
				copyValues[i], copyValues[j] = copyValues[j], copyValues[i]
			}
		}
	}
	index := int(math.Round(Clamp(p, 0, 1) * float64(len(copyValues)-1)))
	return copyValues[index]
}
func Average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}
func Variance(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	mean := Average(values)
	total := 0.0
	for _, value := range values {
		delta := value - mean
		total += delta * delta
	}
	return total / float64(len(values))
}
func StdDev(values []float64) float64 { return math.Sqrt(Variance(values)) }
func MovingAverage(values []float64, window int) []float64 {
	if window < 1 {
		window = 1
	}
	result := make([]float64, len(values))
	for i := range values {
		start := i - window + 1
		if start < 0 {
			start = 0
		}
		result[i] = Average(values[start : i+1])
	}
	return result
}
func ExponentialAverage(values []float64, alpha float64) []float64 {
	alpha = Clamp(alpha, 0, 1)
	result := make([]float64, len(values))
	if len(values) == 0 {
		return result
	}
	result[0] = values[0]
	for i := 1; i < len(values); i++ {
		result[i] = alpha*values[i] + (1-alpha)*result[i-1]
	}
	return result
}
func Difference(values []float64) []float64 {
	if len(values) < 2 {
		return []float64{}
	}
	result := make([]float64, len(values)-1)
	for i := 1; i < len(values); i++ {
		result[i-1] = values[i] - values[i-1]
	}
	return result
}
func Sum(values []float64) float64 {
	total := 0.0
	for _, value := range values {
		total += value
	}
	return total
}
func Min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	value := values[0]
	for _, candidate := range values[1:] {
		if candidate < value {
			value = candidate
		}
	}
	return value
}
func Max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	value := values[0]
	for _, candidate := range values[1:] {
		if candidate > value {
			value = candidate
		}
	}
	return value
}
func Range(values []float64) float64 { return Max(values) - Min(values) }
func Normalize(value, min, max float64) float64 {
	if max == min {
		return 0
	}
	return Clamp((value-min)/(max-min), 0, 1)
}
func Denormalize(value, min, max float64) float64 { return min + Clamp(value, 0, 1)*(max-min) }
func Distance(a, b []float64) float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	total := 0.0
	for i := 0; i < n; i++ {
		delta := a[i] - b[i]
		total += delta * delta
	}
	return math.Sqrt(total)
}
func Dot(a, b []float64) float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	total := 0.0
	for i := 0; i < n; i++ {
		total += a[i] * b[i]
	}
	return total
}
func Cosine(a, b []float64) float64 {
	denom := math.Sqrt(Dot(a, a)) * math.Sqrt(Dot(b, b))
	if denom == 0 {
		return 0
	}
	return Dot(a, b) / denom
}
func Logistic(value float64) float64 {
	if value >= 0 {
		z := math.Exp(-value)
		return 1 / (1 + z)
	}
	z := math.Exp(value)
	return z / (1 + z)
}
func Sigmoid(value, center, slope float64) float64 { return Logistic((value - center) * slope) }
func Softmax(values []float64) []float64 {
	if len(values) == 0 {
		return []float64{}
	}
	maxValue := Max(values)
	result := make([]float64, len(values))
	total := 0.0
	for i, value := range values {
		result[i] = math.Exp(value - maxValue)
		total += result[i]
	}
	for i := range result {
		result[i] /= total
	}
	return result
}
func Entropy(probabilities []float64) float64 {
	total := 0.0
	for _, probability := range probabilities {
		if probability > 0 {
			total -= probability * math.Log(probability)
		}
	}
	return total
}
func NormalizeProbabilities(values []float64) []float64 {
	total := Sum(values)
	result := make([]float64, len(values))
	if total <= 0 {
		return result
	}
	for i, value := range values {
		if value > 0 {
			result[i] = value / total
		}
	}
	return result
}
func Round(value float64, places int) float64 {
	factor := math.Pow10(places)
	return math.Round(value*factor) / factor
}
func IsFinite(value float64) bool { return !math.IsInf(value, 0) && !math.IsNaN(value) }
func Clean(values []float64) []float64 {
	result := make([]float64, 0, len(values))
	for _, value := range values {
		if IsFinite(value) {
			result = append(result, value)
		}
	}
	return result
}
func Fill(value float64, count int) []float64 {
	if count < 0 {
		count = 0
	}
	result := make([]float64, count)
	for i := range result {
		result[i] = value
	}
	return result
}
func Zip(a, b []float64) [][2]float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	result := make([][2]float64, n)
	for i := 0; i < n; i++ {
		result[i] = [2]float64{a[i], b[i]}
	}
	return result
}
func Unzip(values [][2]float64) ([]float64, []float64) {
	a := make([]float64, len(values))
	b := make([]float64, len(values))
	for i, value := range values {
		a[i] = value[0]
		b[i] = value[1]
	}
	return a, b
}
func Add(a, b []float64) []float64 {
	n := len(a)
	if len(b) > n {
		n = len(b)
	}
	result := make([]float64, n)
	for i := range result {
		if i < len(a) {
			result[i] += a[i]
		}
		if i < len(b) {
			result[i] += b[i]
		}
	}
	return result
}
func Scale(values []float64, factor float64) []float64 {
	result := make([]float64, len(values))
	for i, value := range values {
		result[i] = value * factor
	}
	return result
}
func FilterFinite(values []float64) []float64 { return Clean(values) }
func CountAbove(values []float64, threshold float64) int {
	count := 0
	for _, value := range values {
		if value > threshold {
			count++
		}
	}
	return count
}
func CountBelow(values []float64, threshold float64) int {
	count := 0
	for _, value := range values {
		if value < threshold {
			count++
		}
	}
	return count
}
func CountBetween(values []float64, min, max float64) int {
	count := 0
	for _, value := range values {
		if value >= min && value <= max {
			count++
		}
	}
	return count
}
func FractionAbove(values []float64, threshold float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return float64(CountAbove(values, threshold)) / float64(len(values))
}
func TrendDirection(values []float64) string {
	if len(values) < 2 {
		return "stable"
	}
	delta := values[len(values)-1] - values[0]
	if delta > 0 {
		return "rising"
	}
	if delta < 0 {
		return "falling"
	}
	return "stable"
}
func TrendStrength(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	return math.Abs(values[len(values)-1]-values[0]) / (1 + StdDev(values))
}
func IsOutlier(value float64, values []float64, multiplier float64) bool {
	return math.Abs(value-Average(values)) > multiplier*StdDev(values)
}
func ReplaceOutliers(values []float64, multiplier float64) []float64 {
	result := append([]float64(nil), values...)
	mean := Average(values)
	for i, value := range result {
		if IsOutlier(value, values, multiplier) {
			result[i] = mean
		}
	}
	return result
}
func Smooth(values []float64, window int) []float64 {
	return MovingAverage(ReplaceOutliers(values, 3), window)
}
func Covariance(a, b []float64) float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	if n == 0 {
		return 0
	}
	meanA, meanB := Average(a[:n]), Average(b[:n])
	total := 0.0
	for i := 0; i < n; i++ {
		total += (a[i] - meanA) * (b[i] - meanB)
	}
	return total / float64(n)
}
func Correlation(a, b []float64) float64 {
	denom := StdDev(a) * StdDev(b)
	if denom == 0 {
		return 0
	}
	return Covariance(a, b) / denom
}
func PercentChange(previous, current float64) float64 {
	if previous == 0 {
		return 0
	}
	return (current - previous) / math.Abs(previous)
}
func GrowthRate(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	return PercentChange(values[0], values[len(values)-1])
}
func Accumulate(values []float64) []float64 {
	result := make([]float64, len(values))
	total := 0.0
	for i, value := range values {
		total += value
		result[i] = total
	}
	return result
}
func NormalizeCumulative(values []float64) []float64 {
	result := Accumulate(values)
	total := Sum(values)
	if total == 0 {
		return result
	}
	for i := range result {
		result[i] /= total
	}
	return result
}
func Quantize(value, step float64) float64 {
	if step <= 0 {
		return value
	}
	return math.Round(value/step) * step
}
func HistogramCounts(values, buckets []float64) []int {
	hist := NewHistogram(buckets)
	for _, value := range values {
		hist.Add(value)
	}
	return hist.Counts
}
func BucketIndex(value float64, buckets []float64) int {
	for i, bucket := range buckets {
		if value < bucket {
			return i
		}
	}
	return len(buckets)
}
func WeightedAverage(values, weights []float64) float64 {
	n := len(values)
	if len(weights) < n {
		n = len(weights)
	}
	total, weight := 0.0, 0.0
	for i := 0; i < n; i++ {
		total += values[i] * weights[i]
		weight += weights[i]
	}
	if weight == 0 {
		return 0
	}
	return total / weight
}
func TrimmedMean(values []float64, ratio float64) float64 {
	if len(values) == 0 {
		return 0
	}
	copyValues := append([]float64(nil), values...)
	for i := 0; i < len(copyValues); i++ {
		for j := i + 1; j < len(copyValues); j++ {
			if copyValues[j] < copyValues[i] {
				copyValues[i], copyValues[j] = copyValues[j], copyValues[i]
			}
		}
	}
	cut := int(float64(len(copyValues)) * Clamp(ratio, 0, 0.49))
	return Average(copyValues[cut : len(copyValues)-cut])
}
