package domain

import "math"

func NormalizeAngle(value float64) float64 {
	for value < 0 {
		value += 360
	}
	for value >= 360 {
		value -= 360
	}
	return value
}
func AngleDistance(left, right float64) float64 {
	delta := math.Abs(NormalizeAngle(left) - NormalizeAngle(right))
	if delta > 180 {
		return 360 - delta
	}
	return delta
}
func LinearInterpolate(left, right, fraction float64) float64 {
	return left + (right-left)*Clamp(fraction, 0, 1)
}
func Resample(values []float64, size int) []float64 {
	if size <= 0 || len(values) == 0 {
		return []float64{}
	}
	if size == 1 {
		return []float64{values[0]}
	}
	result := make([]float64, size)
	for i := range result {
		fraction := float64(i) / float64(size-1)
		position := fraction * float64(len(values)-1)
		left := int(position)
		right := left + 1
		if right >= len(values) {
			right = len(values) - 1
		}
		result[i] = LinearInterpolate(values[left], values[right], position-float64(left))
	}
	return result
}
func Integrate(values []float64, step float64) float64 {
	if len(values) < 2 || step <= 0 {
		return 0
	}
	total := 0.0
	for i := 1; i < len(values); i++ {
		total += (values[i-1] + values[i]) * step / 2
	}
	return total
}
func Derivative(values []float64, step float64) []float64 {
	if len(values) < 2 || step <= 0 {
		return []float64{}
	}
	result := make([]float64, len(values)-1)
	for i := 1; i < len(values); i++ {
		result[i-1] = (values[i] - values[i-1]) / step
	}
	return result
}
func Energy(values []float64) float64 {
	total := 0.0
	for _, value := range values {
		total += value * value
	}
	return total
}
func NormalizeVector(values []float64) []float64 {
	length := math.Sqrt(Energy(values))
	if length == 0 {
		return Fill(0, len(values))
	}
	return Scale(values, 1/length)
}
func Projection(vector, basis []float64) []float64 {
	denom := Dot(basis, basis)
	if denom == 0 {
		return Fill(0, len(basis))
	}
	factor := Dot(vector, basis) / denom
	return Scale(basis, factor)
}
func Residual(vector, basis []float64) []float64 {
	return Add(vector, Scale(Projection(vector, basis), -1))
}
func Hamming(left, right []byte) int {
	n := len(left)
	if len(right) < n {
		n = len(right)
	}
	distance := absInt(len(left) - len(right))
	for i := 0; i < n; i++ {
		if left[i] != right[i] {
			distance++
		}
	}
	return distance
}
func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
func Levenshtein(left, right string) int {
	a, b := []rune(left), []rune(right)
	previous := make([]int, len(b)+1)
	for j := range previous {
		previous[j] = j
	}
	for i := 1; i <= len(a); i++ {
		current := make([]int, len(b)+1)
		current[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			current[j] = minInt(current[j-1]+1, minInt(previous[j]+1, previous[j-1]+cost))
		}
		previous = current
	}
	return previous[len(b)]
}
func minInt(left, right int) int {
	if left < right {
		return left
	}
	return right
}
func Similarity(left, right string) float64 {
	distance := Levenshtein(left, right)
	maxLen := len([]rune(left))
	if len([]rune(right)) > maxLen {
		maxLen = len([]rune(right))
	}
	if maxLen == 0 {
		return 1
	}
	return 1 - float64(distance)/float64(maxLen)
}
func Jaccard(left, right []string) float64 {
	leftSet := map[string]struct{}{}
	rightSet := map[string]struct{}{}
	for _, value := range left {
		leftSet[value] = struct{}{}
	}
	for _, value := range right {
		rightSet[value] = struct{}{}
	}
	intersection := 0
	union := len(leftSet)
	for value := range rightSet {
		if _, ok := leftSet[value]; ok {
			intersection++
		} else {
			union++
		}
	}
	if union == 0 {
		return 1
	}
	return float64(intersection) / float64(union)
}
func UnionStrings(left, right []string) []string {
	return UniqueStrings(append(append([]string(nil), left...), right...))
}
func IntersectStrings(left, right []string) []string {
	rightSet := map[string]struct{}{}
	for _, value := range right {
		rightSet[value] = struct{}{}
	}
	result := make([]string, 0)
	for _, value := range left {
		if _, ok := rightSet[value]; ok {
			result = append(result, value)
		}
	}
	return UniqueStrings(result)
}
func DifferenceStrings(left, right []string) []string {
	rightSet := map[string]struct{}{}
	for _, value := range right {
		rightSet[value] = struct{}{}
	}
	result := make([]string, 0)
	for _, value := range left {
		if _, ok := rightSet[value]; !ok {
			result = append(result, value)
		}
	}
	return UniqueStrings(result)
}
func MapFloat(values []float64, fn func(float64) float64) []float64 {
	result := make([]float64, len(values))
	for i, value := range values {
		result[i] = fn(value)
	}
	return result
}
func ReduceFloat(values []float64, initial float64, fn func(float64, float64) float64) float64 {
	result := initial
	for _, value := range values {
		result = fn(result, value)
	}
	return result
}
func AnyFloat(values []float64, fn func(float64) bool) bool {
	for _, value := range values {
		if fn(value) {
			return true
		}
	}
	return false
}
func AllFloat(values []float64, fn func(float64) bool) bool {
	for _, value := range values {
		if !fn(value) {
			return false
		}
	}
	return true
}
func FirstFloat(values []float64, fn func(float64) bool) (float64, bool) {
	for _, value := range values {
		if fn(value) {
			return value, true
		}
	}
	return 0, false
}
func ChunkFloat(values []float64, size int) [][]float64 {
	if size < 1 {
		size = 1
	}
	result := make([][]float64, 0)
	for start := 0; start < len(values); start += size {
		end := start + size
		if end > len(values) {
			end = len(values)
		}
		result = append(result, append([]float64(nil), values[start:end]...))
	}
	return result
}
func FlattenFloat(chunks [][]float64) []float64 {
	total := 0
	for _, chunk := range chunks {
		total += len(chunk)
	}
	result := make([]float64, 0, total)
	for _, chunk := range chunks {
		result = append(result, chunk...)
	}
	return result
}
func RepeatFloat(value float64, count int) []float64 { return Fill(value, count) }
func ZipWith(a, b []float64, fn func(float64, float64) float64) []float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	result := make([]float64, n)
	for i := 0; i < n; i++ {
		result[i] = fn(a[i], b[i])
	}
	return result
}
func MaxIndex(values []float64) int {
	if len(values) == 0 {
		return -1
	}
	index := 0
	for i := 1; i < len(values); i++ {
		if values[i] > values[index] {
			index = i
		}
	}
	return index
}
func MinIndex(values []float64) int {
	if len(values) == 0 {
		return -1
	}
	index := 0
	for i := 1; i < len(values); i++ {
		if values[i] < values[index] {
			index = i
		}
	}
	return index
}
func ArgSort(values []float64) []int {
	indices := make([]int, len(values))
	for i := range indices {
		indices[i] = i
	}
	for i := 0; i < len(indices); i++ {
		for j := i + 1; j < len(indices); j++ {
			if values[indices[j]] < values[indices[i]] {
				indices[i], indices[j] = indices[j], indices[i]
			}
		}
	}
	return indices
}
func Take(values []float64, count int) []float64 {
	if count < 0 {
		count = 0
	}
	if count > len(values) {
		count = len(values)
	}
	return append([]float64(nil), values[:count]...)
}
func Drop(values []float64, count int) []float64 {
	if count < 0 {
		count = 0
	}
	if count > len(values) {
		count = len(values)
	}
	return append([]float64(nil), values[count:]...)
}
func Reverse(values []float64) []float64 {
	result := make([]float64, len(values))
	for i, value := range values {
		result[len(values)-1-i] = value
	}
	return result
}
func IsSorted(values []float64) bool {
	for i := 1; i < len(values); i++ {
		if values[i] < values[i-1] {
			return false
		}
	}
	return true
}
func IsStrictlySorted(values []float64) bool {
	for i := 1; i < len(values); i++ {
		if values[i] <= values[i-1] {
			return false
		}
	}
	return true
}
func UniqueFloat(values []float64) []float64 {
	seen := map[float64]struct{}{}
	result := make([]float64, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; !ok {
			seen[value] = struct{}{}
			result = append(result, value)
		}
	}
	return result
}
func Quantiles(values []float64, points []float64) []float64 {
	result := make([]float64, len(points))
	for i, point := range points {
		result[i] = Percentile(values, point)
	}
	return result
}
func ZScores(values []float64) []float64 {
	mean, deviation := Average(values), StdDev(values)
	result := make([]float64, len(values))
	if deviation == 0 {
		return result
	}
	for i, value := range values {
		result[i] = (value - mean) / deviation
	}
	return result
}
func Winsorize(values []float64, low, high float64) []float64 {
	result := append([]float64(nil), values...)
	if len(values) == 0 {
		return result
	}
	lowValue, highValue := Percentile(values, low), Percentile(values, high)
	for i, value := range result {
		result[i] = Clamp(value, lowValue, highValue)
	}
	return result
}
func NormalizeRange(values []float64) []float64 {
	minValue, maxValue := Min(values), Max(values)
	result := make([]float64, len(values))
	for i, value := range values {
		result[i] = Normalize(value, minValue, maxValue)
	}
	return result
}
