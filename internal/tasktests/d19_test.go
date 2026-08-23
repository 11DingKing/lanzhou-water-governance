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

func TestD19RegionalExposureIsolationEvidence(t *testing.T) {
    db:=testsupport.Open(t);uid:=testsupport.SeedUser(t,db.SQL,"d19","regional","Lanzhou");_ = uid;_,err:=db.SQL.Exec(`INSERT INTO manifests(manifest_no,producer_region,carrier_region,facility_region,waste_type,weight_kg,status,created_at) VALUES(?,?,?,?,?,?,?,?)`,"D19","Baiyin","Lanzhou","Linxia","sludge",99,"created",time.Now().UTC().Format(time.RFC3339Nano));if err!=nil{t.Fatal(err)};user:=domain.User{ID:uid,Role:domain.RoleRegional,Region:"Lanzhou"};svc:=service.Reporting{Repo:repository.Reporting{DB:db.SQL}};if _,err=svc.Exposure(context.Background(),user,"Baiyin");!errors.Is(err,domain.ErrForbidden){t.Fatalf("cross-region exposure allowed: %v",err)}
}

var _ *sql.Tx
var _ = errors.New
var _ = time.Now
var _ = context.Background
var _ = sqlite.WithTx
var _ = service.Auth{}
var _ = repository.Users{}
var _ = worker.Job(nil)
