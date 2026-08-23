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

func TestB02WarningAuditAtomic(t *testing.T) {
    db:=testsupport.Open(t); users:=repository.Users{DB:db.SQL}; user,_:=users.Create(context.Background(), "b02", "pw", domain.RoleRegional, "Lanzhou"); agreement:=testsupport.SeedAgreement(t,db.SQL,"Lanzhou","Baiyin"); station:=testsupport.SeedStation(t,db.SQL,"B02-1","Lanzhou"); _=agreement; svc:=service.Collaboration{DB:db.SQL,Repo:repository.Collaboration{DB:db.SQL},Audit:repository.Audit{DB:db.SQL}}; warning,err:=svc.IssueWarning(context.Background(),user,"Lanzhou","Baiyin",station,map[string]any{"quality":"III"}); if err!=nil { t.Fatalf("warning: %v",err) }; var count int; if err=db.SQL.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE object_type='warning' AND object_id=?`, warning.ID).Scan(&count); err!=nil || count!=1 { t.Fatalf("audit count=%d err=%v",count,err) }
}

var _ *sql.Tx
var _ = errors.New
var _ = time.Now
var _ = context.Background
var _ = sqlite.WithTx
var _ = service.Auth{}
var _ = repository.Users{}
var _ = worker.Job(nil)
