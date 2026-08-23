package domain_test

import (
	"testing"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
)

func TestLifecycleTransitions(t *testing.T) {
	cases := []struct {
		name     string
		ok       bool
		from, to domain.AlertStatus
	}{{"open-investigating", true, domain.AlertOpen, domain.AlertInvestigating}, {"investigating-resolved", true, domain.AlertInvestigating, domain.AlertResolved}, {"open-resolved", false, domain.AlertOpen, domain.AlertResolved}, {"resolved-open", false, domain.AlertResolved, domain.AlertOpen}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := domain.CanAlertTransition(tc.from, tc.to); got != tc.ok {
				t.Fatalf("got %v", got)
			}
		})
	}
}
func TestInspectionTransitions(t *testing.T) {
	valid := [][2]domain.InspectionStatus{{domain.InspectionPending, domain.InspectionRunning}, {domain.InspectionRunning, domain.InspectionCompleted}, {domain.InspectionRunning, domain.InspectionFailed}}
	for _, pair := range valid {
		if !domain.CanInspectionTransition(pair[0], pair[1]) {
			t.Fatalf("expected %s -> %s", pair[0], pair[1])
		}
	}
	if domain.CanInspectionTransition(domain.InspectionCompleted, domain.InspectionRunning) {
		t.Fatal("completed inspection reopened")
	}
}
func TestManifestAndProjectTransitions(t *testing.T) {
	if !domain.CanManifestTransition(domain.ManifestAccepted, domain.ManifestDisposed) {
		t.Fatal("manifest cannot dispose")
	}
	if domain.CanManifestTransition(domain.ManifestCreated, domain.ManifestDisposed) {
		t.Fatal("manifest skipped custody")
	}
	if !domain.CanProjectTransition(domain.ProjectPlanned, domain.ProjectBuilding) {
		t.Fatal("project not started")
	}
}
func TestPolicyAndWindows(t *testing.T) {
	if !domain.AllowedRole(domain.RoleRegional, "sample") {
		t.Skip("regional sample intentionally forbidden")
	}
	if domain.AllowedRole(domain.RoleInspector, "compensate") {
		t.Fatal("inspector may not settle")
	}
	now := time.Now()
	if !domain.WithinWindow(now, now.Add(time.Second)) {
		t.Fatal("window rejected")
	}
	if domain.WithinWindow(now.Add(time.Hour), now) {
		t.Fatal("expired window accepted")
	}
}
func TestMetricsNormalize(t *testing.T) {
	input := map[string]float64{"cod": 5, "nan": 0}
	output := domain.NormalizeMetrics(input)
	if len(output) != 2 {
		t.Fatal(output)
	}
	if domain.RiskScore(map[string]float64{"ammonia": 2, "cod": 2}) <= 4 {
		t.Fatal("risk score too low")
	}
}
func TestPagination(t *testing.T) {
	for _, tc := range []struct {
		p             domain.Page
		offset, limit int
	}{{domain.Page{Number: 1, Size: 20}, 0, 20}, {domain.Page{Number: 3, Size: 10}, 20, 10}, {domain.Page{Number: 0, Size: 0}, 0, 50}, {domain.Page{Number: 1, Size: 500}, 0, 50}} {
		if tc.p.Offset() != tc.offset || tc.p.Limit() != tc.limit {
			t.Fatalf("%+v", tc)
		}
	}
}
