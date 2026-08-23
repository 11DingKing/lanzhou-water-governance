package domain_test

import (
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"math"
	"testing"
	"time"
)

func TestAnalyticsHistogram(t *testing.T) {
	h := domain.NewHistogram([]float64{1, 2, 3})
	for _, v := range []float64{0.2, 1.2, 2.2, 4} {
		h.Add(v)
	}
	if h.Total() != 4 || h.Peak() != 1 {
		t.Fatal(h)
	}
}
func TestAnalyticsStatistics(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5}
	if domain.Average(values) != 3 || domain.Min(values) != 1 || domain.Max(values) != 5 {
		t.Fatal(values)
	}
	if domain.Range(values) != 4 || domain.Sum(values) != 15 {
		t.Fatal(values)
	}
}
func TestAnalyticsPercentiles(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5}
	if domain.Percentile(values, 0) != 1 || domain.Percentile(values, 1) != 5 {
		t.Fatal(values)
	}
	if domain.Percentile(values, .5) != 3 {
		t.Fatal(values)
	}
}
func TestAnalyticsVariance(t *testing.T) {
	values := []float64{1, 2, 3}
	if domain.Round(domain.Variance(values), 3) != .667 {
		t.Fatal(domain.Variance(values))
	}
	if domain.StdDev(values) <= 0 {
		t.Fatal("stddev")
	}
}
func TestAnalyticsMoving(t *testing.T) {
	values := []float64{1, 2, 3, 4}
	got := domain.MovingAverage(values, 2)
	if got[0] != 1 || got[1] != 1.5 || got[3] != 3.5 {
		t.Fatal(got)
	}
}
func TestAnalyticsExponential(t *testing.T) {
	got := domain.ExponentialAverage([]float64{1, 3, 5}, .5)
	if len(got) != 3 || got[1] != 2 {
		t.Fatal(got)
	}
}
func TestAnalyticsDifferences(t *testing.T) {
	got := domain.Difference([]float64{1, 4, 9})
	if len(got) != 2 || got[0] != 3 || got[1] != 5 {
		t.Fatal(got)
	}
}
func TestAnalyticsNormalize(t *testing.T) {
	if domain.Normalize(5, 0, 10) != .5 || domain.Denormalize(.5, 0, 10) != 5 {
		t.Fatal("normalize")
	}
	if domain.Normalize(5, 5, 5) != 0 {
		t.Fatal("equal")
	}
}
func TestAnalyticsVectors(t *testing.T) {
	a := []float64{1, 2}
	b := []float64{2, 3}
	if domain.Dot(a, b) != 8 {
		t.Fatal("dot")
	}
	if domain.Distance(a, b) != math.Sqrt(2) {
		t.Fatal("distance")
	}
	if domain.Cosine(a, b) <= .9 {
		t.Fatal("cosine")
	}
}
func TestAnalyticsLogistic(t *testing.T) {
	if domain.Logistic(0) != .5 {
		t.Fatal("logistic")
	}
	soft := domain.Softmax([]float64{1, 2, 3})
	if domain.Round(domain.Sum(soft), 5) != 1 {
		t.Fatal(soft)
	}
}
func TestAnalyticsEntropy(t *testing.T) {
	if domain.Entropy([]float64{.5, .5}) <= .6 {
		t.Fatal("entropy")
	}
	if len(domain.NormalizeProbabilities([]float64{1, 1})) != 2 {
		t.Fatal("prob")
	}
}
func TestAnalyticsClean(t *testing.T) {
	clean := domain.Clean([]float64{1, math.NaN(), math.Inf(1), 2})
	if len(clean) != 2 {
		t.Fatal(clean)
	}
	if !domain.IsFinite(2) {
		t.Fatal("finite")
	}
}
func TestAnalyticsZip(t *testing.T) {
	pairs := domain.Zip([]float64{1, 2}, []float64{3, 4})
	a, b := domain.Unzip(pairs)
	if len(a) != 2 || b[1] != 4 {
		t.Fatal(a, b)
	}
}
func TestAnalyticsAddScale(t *testing.T) {
	if got := domain.Add([]float64{1, 2}, []float64{3}); got[0] != 4 || got[1] != 2 {
		t.Fatal(got)
	}
	if domain.Scale([]float64{2, 3}, 2)[1] != 6 {
		t.Fatal("scale")
	}
}
func TestAnalyticsCounts(t *testing.T) {
	values := []float64{1, 2, 3, 4}
	if domain.CountAbove(values, 2) != 2 || domain.CountBelow(values, 3) != 2 || domain.CountBetween(values, 2, 3) != 2 {
		t.Fatal("counts")
	}
}
func TestAnalyticsTrend(t *testing.T) {
	if domain.TrendDirection([]float64{1, 2, 3}) != "rising" {
		t.Fatal("trend")
	}
	if domain.TrendDirection([]float64{3, 2, 1}) != "falling" {
		t.Fatal("trend")
	}
	if domain.TrendDirection([]float64{2, 2}) != "stable" {
		t.Fatal("trend")
	}
}
func TestAnalyticsOutliers(t *testing.T) {
	values := []float64{1, 1, 1, 100}
	if !domain.IsOutlier(100, values, 1) {
		t.Fatal("outlier")
	}
	clean := domain.ReplaceOutliers(values, 1)
	if clean[3] == 100 {
		t.Fatal(clean)
	}
}
func TestAnalyticsCorrelation(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{2, 4, 6}
	if domain.Correlation(a, b) < .99 {
		t.Fatal("correlation")
	}
	if domain.Covariance(a, b) <= 0 {
		t.Fatal("covariance")
	}
}
func TestAnalyticsGrowth(t *testing.T) {
	if domain.PercentChange(10, 15) != .5 {
		t.Fatal("change")
	}
	if domain.GrowthRate([]float64{10, 20}) != 1 {
		t.Fatal("growth")
	}
}
func TestAnalyticsAccumulate(t *testing.T) {
	got := domain.Accumulate([]float64{1, 2, 3})
	if got[2] != 6 {
		t.Fatal(got)
	}
	if domain.NormalizeCumulative([]float64{1, 1})[1] != 1 {
		t.Fatal("cumulative")
	}
}
func TestAnalyticsQuantize(t *testing.T) {
	if math.Abs(domain.Quantize(1.24, .1)-1.2) > 1e-9 {
		t.Fatal("quantize")
	}
	if math.Abs(domain.Quantize(1.26, .1)-1.3) > 1e-9 {
		t.Fatal("quantize")
	}
}
func TestAnalyticsBuckets(t *testing.T) {
	if domain.BucketIndex(2, []float64{1, 3}) != 1 {
		t.Fatal("bucket")
	}
	counts := domain.HistogramCounts([]float64{1, 2, 4}, []float64{2, 3})
	if len(counts) != 3 {
		t.Fatal(counts)
	}
}
func TestAnalyticsWeighted(t *testing.T) {
	if domain.WeightedAverage([]float64{1, 3}, []float64{1, 3}) != 2.5 {
		t.Fatal("weighted")
	}
	if domain.TrimmedMean([]float64{1, 2, 3, 100}, .25) != 2.5 {
		t.Fatal("trimmed")
	}
}
func TestAnalyticsResample(t *testing.T) {
	got := domain.Resample([]float64{0, 10}, 3)
	if len(got) != 3 || got[1] != 5 {
		t.Fatal(got)
	}
}
func TestAnalyticsIntegrateDerivative(t *testing.T) {
	if domain.Integrate([]float64{1, 1, 1}, 1) != 2 {
		t.Fatal("integrate")
	}
	if domain.Derivative([]float64{1, 3, 6}, 1)[1] != 3 {
		t.Fatal("derivative")
	}
}
func TestAnalyticsProjection(t *testing.T) {
	projection := domain.Projection([]float64{1, 2}, []float64{1, 0})
	if projection[0] != 1 || projection[1] != 0 {
		t.Fatal(projection)
	}
	if domain.Energy([]float64{3, 4}) != 25 {
		t.Fatal("energy")
	}
}
func TestAnalyticsStringSimilarity(t *testing.T) {
	if domain.Similarity("river", "river") != 1 {
		t.Fatal("similarity")
	}
	if domain.Levenshtein("kitten", "sitting") != 3 {
		t.Fatal("levenshtein")
	}
	if domain.Hamming([]byte("abc"), []byte("axc")) != 1 {
		t.Fatal("hamming")
	}
}
func TestAnalyticsSetOperations(t *testing.T) {
	if domain.Jaccard([]string{"a", "b"}, []string{"b", "c"}) != 1.0/3.0 {
		t.Fatal("jaccard")
	}
	if len(domain.UnionStrings([]string{"a"}, []string{"a", "b"})) != 2 {
		t.Fatal("union")
	}
	if len(domain.IntersectStrings([]string{"a", "b"}, []string{"b"})) != 1 {
		t.Fatal("intersection")
	}
}
func TestAnalyticsHigherOrder(t *testing.T) {
	mapped := domain.MapFloat([]float64{1, 2}, func(v float64) float64 { return v * 2 })
	if mapped[1] != 4 {
		t.Fatal(mapped)
	}
	if domain.ReduceFloat(mapped, 0, func(a, b float64) float64 { return a + b }) != 6 {
		t.Fatal("reduce")
	}
	if !domain.AnyFloat(mapped, func(v float64) bool { return v == 4 }) {
		t.Fatal("any")
	}
}
func TestAnalyticsChunks(t *testing.T) {
	chunks := domain.ChunkFloat([]float64{1, 2, 3, 4, 5}, 2)
	if len(chunks) != 3 {
		t.Fatal(chunks)
	}
	if len(domain.FlattenFloat(chunks)) != 5 {
		t.Fatal("flatten")
	}
}
func TestAnalyticsOrdering(t *testing.T) {
	if domain.MaxIndex([]float64{1, 5, 3}) != 1 || domain.MinIndex([]float64{1, 5, 3}) != 0 {
		t.Fatal("index")
	}
	if !domain.IsSorted([]float64{1, 2, 3}) || domain.IsStrictlySorted([]float64{1, 1}) {
		t.Fatal("sorted")
	}
}
func TestAnalyticsSelection(t *testing.T) {
	if len(domain.Take([]float64{1, 2, 3}, 2)) != 2 || len(domain.Drop([]float64{1, 2, 3}, 2)) != 1 {
		t.Fatal("take/drop")
	}
	if domain.Reverse([]float64{1, 2})[0] != 2 {
		t.Fatal("reverse")
	}
}
func TestAnalyticsUnique(t *testing.T) {
	if len(domain.UniqueFloat([]float64{1, 1, 2})) != 2 {
		t.Fatal("unique")
	}
	if len(domain.ArgSort([]float64{3, 1, 2})) != 3 {
		t.Fatal("argsort")
	}
}
func TestAnalyticsQuantiles(t *testing.T) {
	if len(domain.Quantiles([]float64{1, 2, 3}, []float64{.25, .75})) != 2 {
		t.Fatal("quantiles")
	}
	if len(domain.ZScores([]float64{1, 2, 3})) != 3 {
		t.Fatal("zscores")
	}
	if len(domain.Winsorize([]float64{1, 2, 100}, .1, .9)) != 3 {
		t.Fatal("winsor")
	}
}
func TestDomainCalendars(t *testing.T) {
	calendar := domain.NewCalendar("UTC")
	now := time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)
	if !calendar.IsBusinessHour(now) {
		t.Fatal("business hour")
	}
	if calendar.Quarter(now) != 3 || calendar.Period(now) != "2026-08" {
		t.Fatal("period")
	}
	if !calendar.SameDay(now, now.Add(time.Hour)) {
		t.Fatal("day")
	}
}
func TestDomainRules(t *testing.T) {
	rules := domain.DefaultRules()
	if err := domain.ValidateRules(rules); err != nil {
		t.Fatal(err)
	}
	if !domain.ApplyRule(domain.FindRule(rules, domain.QualityIII), map[string]float64{"cod": 40}) {
		t.Fatal("rule")
	}
	if len(domain.RuleCodes(rules)) != 4 {
		t.Fatal("codes")
	}
}
func TestDomainAllocation(t *testing.T) {
	allocation := domain.Allocation{Capacity: 100}
	if err := allocation.Reserve(40); err != nil {
		t.Fatal(err)
	}
	if err := allocation.Commit(20); err != nil {
		t.Fatal(err)
	}
	if allocation.Available() != 60 {
		t.Fatal(allocation)
	}
}
func TestDomainSerialization(t *testing.T) {
	raw, err := domain.EncodeEnvelope("sample", map[string]any{"id": 1})
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := domain.DecodeEnvelope(raw)
	if err != nil || envelope.Type != "sample" {
		t.Fatal(envelope, err)
	}
	if len(domain.MapKeys(domain.MergeMaps(map[string]any{"a": 1}, map[string]any{"b": 2}))) != 2 {
		t.Fatal("merge")
	}
}
func TestDomainQueues(t *testing.T) {
	now := time.Now()
	items := []domain.QueueItem{{ID: "a", Priority: 1, AvailableAt: now}, {ID: "b", Priority: 2, AvailableAt: now}}
	ready := domain.QueueReady(items, now)
	if ready[0].ID != "b" {
		t.Fatal(ready)
	}
	if len(domain.RemoveQueueItem(items, "a")) != 1 {
		t.Fatal("remove")
	}
}
func TestDomainChecklist(t *testing.T) {
	items := []domain.ChecklistItem{{Code: "a", Required: true}, {Code: "b", Required: false}}
	if domain.ChecklistComplete(items) {
		t.Fatal("incomplete")
	}
	items = domain.CompleteChecklist(items, "a", "photo")
	if !domain.ChecklistComplete(items) {
		t.Fatal("complete")
	}
}
func TestDomainNotification(t *testing.T) {
	notification := domain.Notification{Recipient: "a", Channel: "email", Body: "hello"}
	if !notification.Valid() {
		t.Fatal("notification")
	}
	if len(domain.BatchNotifications([]domain.Notification{notification})) != 1 {
		t.Fatal("batch")
	}
}
func TestDomainExport(t *testing.T) {
	job := domain.ExportJob{ID: "1", Format: "csv", Status: "queued"}
	if err := job.Start(); err != nil {
		t.Fatal(err)
	}
	if err := job.Complete(2); err != nil || !job.CanDownload() {
		t.Fatal(job, err)
	}
	if job.Filename("report") == "" {
		t.Fatal("filename")
	}
}
func TestDomainThresholds(t *testing.T) {
	if domain.ClassifyMetric("cod", 10) != domain.QualityII {
		t.Fatal("classify")
	}
	if domain.WorstClass(domain.QualityII, domain.QualityIV) != domain.QualityIV {
		t.Fatal("worst")
	}
	if !domain.IsSafeForDischarge(domain.QualityII) {
		t.Fatal("safe")
	}
}
func TestDomainActors(t *testing.T) {
	actor := domain.Actor{User: domain.User{Role: domain.RoleInspector}}
	if actor.Can("admin") {
		t.Fatal("scope")
	}
	actor = actor.AddScope("inspect")
	if !actor.Can("inspect") {
		t.Fatal("add scope")
	}
	actor = actor.RemoveScope("inspect")
	if actor.Can("inspect") {
		t.Fatal("remove scope")
	}
}
func TestDomainOperations(t *testing.T) {
	operation := domain.Operation{Status: "queued"}
	now := time.Now()
	if err := operation.Start(now); err != nil {
		t.Fatal(err)
	}
	if err := operation.Finish(now.Add(time.Second)); err != nil || operation.Duration() != time.Second {
		t.Fatal(operation, err)
	}
}
func TestDomainCursor(t *testing.T) {
	cursor := domain.NextCursor("now", 5)
	decoded, err := domain.DecodeCursor(cursor.Encode())
	if err != nil || decoded.ID != 5 {
		t.Fatal(decoded, err)
	}
}
func TestDomainLimits(t *testing.T) {
	limits := domain.DefaultLimits()
	if !limits.Valid() || !limits.CanSample(1) || limits.CanAlert(limits.MaxAlertsPerStation) {
		t.Fatal(limits)
	}
}
func TestDomainGeo(t *testing.T) {
	point := domain.Coordinate{Latitude: 36, Longitude: 104}
	if !point.Valid() || !point.Within(point, 0) {
		t.Fatal(point)
	}
	if domain.AverageCoordinate([]domain.Coordinate{point, point}) != point {
		t.Fatal("average")
	}
}
func TestDomainStatus(t *testing.T) {
	if !domain.ValidStatus("alert", "open") {
		t.Fatal("status")
	}
	if domain.NormalizeStatus(" OPEN ") != "open" {
		t.Fatal("normalize")
	}
}
func TestDomainWeights(t *testing.T) {
	if domain.WeightedRisk(map[string]float64{"cod": 2}, map[string]float64{"cod": 2}) != 4 {
		t.Fatal("weight")
	}
	if len(domain.TopRisk(map[string]float64{"cod": 2, "ammonia": 3}, 1)) != 1 {
		t.Fatal("top")
	}
}
