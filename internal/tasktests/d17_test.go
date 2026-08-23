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

func TestD17StationSnapshotRecoveryEvidence(t *testing.T) {
    db:=testsupport.Open(t);station:=domain.Station{ID:1,Region:"Lanzhou",Code:"D17-1"};svc:=service.SnapshotService{Repo:repository.Snapshots{DB:db.SQL}};user:=domain.User{ID:1,Role:domain.RoleAdmin,Region:"Lanzhou"};if err:=svc.SaveStation(context.Background(),user,station,domain.QualityTrend{StationID:1,Samples:1});err!=nil{t.Fatal(err)};if _,err:=svc.LoadStation(context.Background(),user,station);err!=nil{t.Fatalf("snapshot recovery failed: %v",err)}
}

var _ *sql.Tx
var _ = errors.New
var _ = time.Now
var _ = context.Background
var _ = sqlite.WithTx
var _ = service.Auth{}
var _ = repository.Users{}
var _ = worker.Job(nil)
