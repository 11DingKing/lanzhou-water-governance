package domain_test

import (
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"testing"
	"time"
)

func TestBusinessValidation(t *testing.T) {
	station := domain.Station{ID: 1, Region: "Lanzhou"}
	user := domain.User{ID: 2, Role: domain.RoleInspector, Region: "Lanzhou"}
	sample := domain.Sample{StationID: 1, CreatedBy: 2, SampledAt: time.Now(), Metrics: map[string]float64{"cod": 1}}
	if len(domain.ValidateSample(sample, station, user)) != 0 {
		t.Fatal("valid sample rejected")
	}
	if len(domain.ValidateSample(domain.Sample{StationID: 2}, station, user)) == 0 {
		t.Fatal("invalid sample accepted")
	}
}
func TestManifestValidation(t *testing.T) {
	valid := domain.Manifest{Number: "M", WeightKG: 10, ProducerRegion: "Lanzhou", FacilityRegion: "Baiyin"}
	if len(domain.ValidateManifest(valid)) != 0 {
		t.Fatal("manifest")
	}
	if len(domain.ValidateManifest(domain.Manifest{})) == 0 {
		t.Fatal("empty manifest")
	}
}
func TestProjectValidation(t *testing.T) {
	if err := domain.ValidateProject(domain.Project{Name: "p", TargetHectares: 1, BudgetCents: 1}); err != nil {
		t.Fatal(err)
	}
	if err := domain.ValidateProject(domain.Project{}); err == nil {
		t.Fatal("empty project")
	}
}
func TestSettlement(t *testing.T) {
	input := domain.SettlementInput{BaseCents: 10000, Improvement: .2, Eligible: true}
	if domain.SettlementAmount(input) <= 10000 {
		t.Fatal("settlement")
	}
	if domain.SettlementStatus(0) != "blocked" {
		t.Fatal("status")
	}
	if domain.SettlementStatus(20000) != "approved" {
		t.Fatal("status")
	}
}
func TestMaintenance(t *testing.T) {
	calendar := domain.NewCalendar("UTC")
	window := domain.MaintenanceWindow{Start: time.Now(), End: time.Now().Add(time.Hour)}
	if !window.Active(time.Now().Add(time.Minute)) {
		t.Fatal("window")
	}
	if calendar.DaysBetween(time.Now(), time.Now().Add(48*time.Hour)) != 2 {
		t.Fatal("days")
	}
	if !domain.ValidRetention(domain.RetentionPolicy{AuditDays: 30, SampleDays: 365, SessionDays: 1}) {
		t.Fatal("retention")
	}
}
func TestEvents(t *testing.T) {
	event := domain.Event{ID: "1", Type: domain.EventSampleRecorded, ObjectType: "sample", ObjectID: "1", OccurredAt: time.Now()}
	if !event.IsValid() || domain.EventKey(event) != "sample.recorded:sample:1" {
		t.Fatal(event)
	}
	if domain.EventNeedsRetry(domain.Event{Type: domain.EventWarningSent}) != true {
		t.Fatal("retry")
	}
}
func TestFilters(t *testing.T) {
	filter := domain.TextFilter{Query: "river", Fields: []string{"name"}}
	if !filter.Match(map[string]string{"name": "Yellow River"}) {
		t.Fatal("filter")
	}
	if domain.NormalizeFilter(" A  B ") != "a b" {
		t.Fatal("normalize")
	}
	if len(domain.ParseCSV("a,b,a")) != 3 || len(domain.UniqueStrings([]string{"a", "a"})) != 1 {
		t.Fatal("strings")
	}
}
func TestNotifications(t *testing.T) {
	n := domain.Notification{Recipient: "r", Channel: "sms", Body: "hello"}
	if !n.Valid() || n.DedupKey() == "" {
		t.Fatal(n)
	}
	if domain.TruncateBody("abcdef", 3) != "abc" {
		t.Fatal("truncate")
	}
}
func TestExportFailure(t *testing.T) {
	job := domain.ExportJob{ID: "x", Status: "queued"}
	if err := job.Fail("bad"); err == nil {
		t.Fatal("fail before start")
	}
	if err := job.Start(); err != nil {
		t.Fatal(err)
	}
	if err := job.Fail("bad"); err != nil || job.Status != "failed" {
		t.Fatal(job, err)
	}
}
func TestActorRegion(t *testing.T) {
	actor := domain.Actor{User: domain.User{Role: domain.RoleRegional, Region: "Lanzhou"}}
	if !actor.CanRegion("Lanzhou") || actor.CanRegion("Baiyin") {
		t.Fatal("region")
	}
	if len(actor.ScopeSet()) != 0 {
		t.Fatal("scopes")
	}
}
func TestOperationFailure(t *testing.T) {
	op := domain.Operation{Status: "queued"}
	now := time.Now()
	if err := op.Start(now); err != nil {
		t.Fatal(err)
	}
	if err := op.Fail(now.Add(time.Second), "x"); err != nil || op.Error != "x" {
		t.Fatal(op, err)
	}
}
func TestCursorEmpty(t *testing.T) {
	if !domain.EmptyCursor("") || domain.EmptyCursor("x") {
		t.Fatal("cursor")
	}
}
func TestLimits(t *testing.T) {
	limits := domain.DefaultLimits()
	if limits.CanInspection(limits.MaxPendingInspections) {
		t.Fatal("limit")
	}
	if limits.CanAlert(0) == false {
		t.Fatal("alert")
	}
}
func TestCoordinates(t *testing.T) {
	invalid := domain.Coordinate{Latitude: 100}
	if invalid.Valid() {
		t.Fatal("invalid coordinate")
	}
	left := domain.Coordinate{Latitude: 1, Longitude: 1}
	right := domain.Coordinate{Latitude: 1, Longitude: 2}
	if left.Distance(right) != 1 {
		t.Fatal("distance")
	}
}
func TestWeights(t *testing.T) {
	normalized := domain.NormalizeWeights(map[string]float64{"a": 1, "b": 1})
	if normalized["a"] != .5 {
		t.Fatal(normalized)
	}
	if len(domain.TopRisk(map[string]float64{}, 1)) != 0 {
		t.Fatal("top")
	}
}
func TestTrend(t *testing.T) {
	trend := domain.BuildTrend(1, []domain.Reading{{Class: domain.QualityII, Metrics: map[string]float64{"cod": 1}}, {Class: domain.QualityIII, Metrics: map[string]float64{"cod": 2}}})
	if trend.Samples != 2 || trend.Worst != domain.QualityIII || !trend.Stable {
		t.Fatal(trend)
	}
}
func TestRegion(t *testing.T) {
	if domain.NormalizeRegion(" lanzhou ") != "LANZHOU" {
		t.Fatal("region")
	}
	if domain.CompensationDirection(domain.Upstream) != "upstream-to-downstream" {
		t.Fatal("direction")
	}
	if !domain.ValidRegion("Lanzhou") {
		t.Fatal("valid region")
	}
}
func TestRulesDisabled(t *testing.T) {
	rules := []domain.Rule{{Code: "x", Class: domain.QualityII, Enabled: false}}
	if domain.ApplyRule(rules[0], map[string]float64{"cod": 100}) {
		t.Fatal("disabled rule")
	}
	if domain.FindRule(rules, domain.QualityV).Enabled {
		t.Fatal("missing")
	}
}
func TestAllocationRebalance(t *testing.T) {
	items := []domain.Allocation{{Region: "a", Capacity: 10}, {Region: "b", Capacity: 20}}
	result := domain.Rebalance(items, 15)
	if domain.TotalAvailable(result) != 15 {
		t.Fatal(result)
	}
	if len(domain.Balance(result)) != 2 {
		t.Fatal("balance")
	}
}
func TestSerializationClones(t *testing.T) {
	source := map[string]any{"a": 1}
	copy := domain.CloneMap(source)
	copy["b"] = 2
	if len(source) != 1 {
		t.Fatal("clone")
	}
	if len(domain.CloneStrings([]string{"a", "b"})) != 2 {
		t.Fatal("strings")
	}
}
func TestQueueRetry(t *testing.T) {
	item := domain.QueueItem{ID: "x", Attempts: 1, AvailableAt: time.Now()}
	next := item.Next(time.Now())
	if next.Attempts != 2 || next.Terminal(3) {
		t.Fatal(next)
	}
}
func TestChecklistValidation(t *testing.T) {
	items := []domain.ChecklistItem{{Code: "a", Required: true}}
	if err := domain.ValidateChecklist(items); err != nil {
		t.Fatal(err)
	}
	if domain.RequiredChecklist(items) != 1 || len(domain.MissingChecklist(items)) != 1 {
		t.Fatal("checklist")
	}
}
func TestStatusList(t *testing.T) {
	statuses := domain.Statuses("manifest")
	if len(statuses) != 4 {
		t.Fatal(statuses)
	}
	if domain.ValidStatus("unknown", "x") {
		t.Fatal("unknown")
	}
}
func TestThresholdClassify(t *testing.T) {
	if !domain.ValidQualityClass(domain.QualityV) || !domain.ValidRole(domain.RoleAdmin) {
		t.Fatal("valid")
	}
	if !domain.ValidStationCode("MON-1") {
		t.Fatal("code")
	}
	if domain.ValidStationCode("bad") {
		t.Fatal("code")
	}
}
func TestMetricDeltaMerge(t *testing.T) {
	delta := domain.MetricDelta(map[string]float64{"cod": 1}, map[string]float64{"cod": 3, "ammonia": 2})
	if delta["cod"] != 2 || delta["ammonia"] != 2 {
		t.Fatal(delta)
	}
	merged := domain.MergeMetrics(map[string]float64{"cod": 1}, map[string]float64{"cod": 3})
	if merged["cod"] != 2 {
		t.Fatal(merged)
	}
}
func TestSanitize(t *testing.T) {
	if domain.SanitizeNotes(" x ") != "x" {
		t.Fatal("notes")
	}
	if !domain.RequireFields("a", "b") || domain.RequireFields("a", "") {
		t.Fatal("fields")
	}
}
func TestCalendarBusiness(t *testing.T) {
	cal := domain.NewCalendar("UTC")
	at := time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)
	if !cal.IsBusinessHour(at) || !domain.WithinWindow(at, at.Add(time.Hour)) {
		t.Fatal("window")
	}
}

func TestAdditionalValidationCases(t *testing.T) {
	if domain.ValidQualityClass("bad") {
		t.Fatal("bad quality")
	}
	if domain.ValidRole("bad") {
		t.Fatal("bad role")
	}
	if domain.ValidStationCode("MON-000") == false {
		t.Fatal("code")
	}
}
func TestAdditionalCalendarCases(t *testing.T) {
	cal := domain.NewCalendar("UTC")
	now := time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)
	if cal.EndOfDay(now).Before(now) {
		t.Fatal("end")
	}
	if cal.AddBusinessHours(now, 2).Before(now) {
		t.Fatal("add")
	}
}
func TestAdditionalSettlementCases(t *testing.T) {
	if domain.SettlementAmount(domain.SettlementInput{BaseCents: 100, Eligible: false}) != 0 {
		t.Fatal("ineligible")
	}
	if domain.SettlementAmount(domain.SettlementInput{BaseCents: 100, Eligible: true, Violations: 10}) != 20 {
		t.Fatal("penalty")
	}
}
func TestAdditionalQueueCases(t *testing.T) {
	now := time.Now()
	item := domain.QueueItem{ID: "q", AvailableAt: now.Add(time.Hour)}
	if item.Ready(now) {
		t.Fatal("future")
	}
	if !item.Next(now).AvailableAt.After(now) {
		t.Fatal("retry")
	}
}
func TestAdditionalChecklistCases(t *testing.T) {
	items := []domain.ChecklistItem{{Code: "a", Required: false}}
	if !domain.ChecklistComplete(items) {
		t.Fatal("optional")
	}
	if len(domain.CompleteChecklist(items, "missing", "x")) != 1 {
		t.Fatal("complete")
	}
}
func TestAdditionalExportCases(t *testing.T) {
	job := domain.ExportJob{ID: "x", Status: "queued"}
	if err := job.Start(); err != nil {
		t.Fatal(err)
	}
	if err := job.Complete(-1); err == nil {
		t.Fatal("negative rows")
	}
}
func TestAdditionalThresholdCases(t *testing.T) {
	if domain.ClassifyMetrics(map[string]float64{"cod": 100}) != domain.QualityV {
		t.Fatal("worst")
	}
	if domain.RequiredInspectionDeadline(domain.QualityV, time.Now()).Before(time.Now()) {
		t.Fatal("deadline")
	}
}
func TestAdditionalRegionCases(t *testing.T) {
	agreement := domain.Agreement{UpstreamRegion: "Lanzhou", DownstreamRegion: "Baiyin", Active: true}
	if !domain.CanExchangeData("Lanzhou", "Baiyin", agreement) {
		t.Fatal("exchange")
	}
}
func TestAdditionalActorCases(t *testing.T) {
	admin := domain.Actor{User: domain.User{Role: domain.RoleAdmin}}
	if !admin.Can("anything") {
		t.Fatal("admin")
	}
}
func TestAdditionalAnalyticsCases(t *testing.T) {
	if domain.Round(domain.Logistic(0), 2) != .5 {
		t.Fatal("logistic")
	}
	if len(domain.NormalizeRange([]float64{1, 2, 3})) != 3 {
		t.Fatal("range")
	}
}
func TestAdditionalMapCases(t *testing.T) {
	if len(domain.RegionSet([]string{"Lanzhou", "Lanzhou"})) != 1 {
		t.Fatal("regions")
	}
	if domain.Redact("abcdef") != "ab****ef" {
		t.Fatal("redact")
	}
}
func TestAdditionalStatusCases(t *testing.T) {
	if domain.NormalizeAction(" Issue ") != "issue" {
		t.Fatal("action")
	}
	if domain.CompensationDirection(domain.CrossBorder) != "bilateral" {
		t.Fatal("direction")
	}
}
func TestAdditionalMathCases(t *testing.T) {
	if domain.PercentChange(0, 1) != 0 {
		t.Fatal("zero baseline")
	}
	if domain.IsFinite(1) == false {
		t.Fatal("finite")
	}
}
func TestAdditionalFinalCase(t *testing.T) {
	if domain.NextRetry(time.Now(), 1).IsZero() {
		t.Fatal("retry")
	}
}
