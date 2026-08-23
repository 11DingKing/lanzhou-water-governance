package tasktests_test

import (
    "context"
    "database/sql"
    "errors"
    "testing"
    "time"
    "github.com/11DingKing/lanzhou-water-governance/internal/domain"
    "github.com/11DingKing/lanzhou-water-governance/internal/repository"
    "github.com/11DingKing/lanzhou-water-governance/internal/service"
    "github.com/11DingKing/lanzhou-water-governance/internal/storage/sqlite"
    "github.com/11DingKing/lanzhou-water-governance/internal/testsupport"
    "github.com/11DingKing/lanzhou-water-governance/internal/worker"
)

func TestB01SampleAlertInspectionAtomic(t *testing.T) {
    db := testsupport.Open(t); users := repository.Users{DB: db.SQL}; user,_ := users.Create(context.Background(), "b01", "pw", domain.RoleInspector, "Lanzhou"); station := testsupport.SeedStation(t, db.SQL, "B01-1", "Lanzhou"); svc := service.Monitoring{DB:db.SQL, Repo:repository.Monitoring{DB:db.SQL}, Inspections:repository.Inspections{DB:db.SQL}}; _, alert, err := svc.RecordSample(context.Background(), user, domain.Sample{StationID:station, Quality:domain.QualityIII, SampledAt:time.Now(), Metrics:map[string]float64{"cod":40}}); if err != nil { t.Fatalf("sample failed: %v", err) }; if alert.ID==0 { t.Fatalf("alert missing") }; var count int; if err=db.SQL.QueryRow(`SELECT COUNT(*) FROM inspections WHERE alert_id=?`, alert.ID).Scan(&count); err!=nil || count!=1 { t.Fatalf("inspection count=%d err=%v", count,err) }
}

var _ *sql.Tx
var _ = errors.New
var _ = time.Now
var _ = context.Background
var _ = sqlite.WithTx
var _ = service.Auth{}
var _ = repository.Users{}
var _ = worker.Job(nil)
