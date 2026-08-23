package service_test

import (
	"context"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/service"
	"github.com/11DingKing/lanzhou-water-governance/internal/testsupport"
	"testing"
	"time"
)

func TestMonitoringCreatesInspectionForPoorQuality(t *testing.T) {
	db := testsupport.Open(t)
	users := repository.Users{DB: db.SQL}
	user, _ := users.Create(context.Background(), "inspector", "pw", domain.RoleInspector, "Lanzhou")
	station := testsupport.SeedStation(t, db.SQL, "MON-1", "Lanzhou")
	svc := service.Monitoring{DB: db.SQL, Repo: repository.Monitoring{DB: db.SQL}, Inspections: repository.Inspections{DB: db.SQL}}
	sample, alert, err := svc.RecordSample(context.Background(), user, domain.Sample{StationID: station, Quality: domain.QualityIII, Metrics: map[string]float64{"cod": 20}, SampledAt: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	if sample.ID == 0 || alert.ID == 0 {
		t.Fatalf("sample=%+v alert=%+v", sample, alert)
	}
	var status string
	if err = db.SQL.QueryRow(`SELECT status FROM inspections WHERE alert_id=?`, alert.ID).Scan(&status); err != nil || status != "pending" {
		t.Fatalf("%s %v", status, err)
	}
}
func TestMonitoringRejectsCrossRegionSample(t *testing.T) {
	db := testsupport.Open(t)
	users := repository.Users{DB: db.SQL}
	user, _ := users.Create(context.Background(), "inspector", "pw", domain.RoleInspector, "Lanzhou")
	station := testsupport.SeedStation(t, db.SQL, "MON-2", "Baiyin")
	svc := service.Monitoring{DB: db.SQL, Repo: repository.Monitoring{DB: db.SQL}, Inspections: repository.Inspections{DB: db.SQL}}
	if _, _, err := svc.RecordSample(context.Background(), user, domain.Sample{StationID: station, Quality: domain.QualityII}); err != domain.ErrForbidden {
		t.Fatalf("err=%v", err)
	}
}
