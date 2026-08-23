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

func TestD20SampleArchiveReferenceEvidence(t *testing.T) {
    _ = domain.RoleAdmin; db:=testsupport.Open(t);station:=testsupport.SeedStation(t,db.SQL,"D20-1","Lanzhou");uid:=testsupport.SeedUser(t,db.SQL,"d20","inspector","Lanzhou");res,err:=db.SQL.Exec(`INSERT INTO samples(station_id,sampled_at,quality_class,metrics_json,created_by,created_at) VALUES(?,?,?,?,?,?)`,station,time.Now().Add(-48*time.Hour).Format(time.RFC3339Nano),"II","{}",uid,time.Now().UTC().Format(time.RFC3339Nano));if err!=nil{t.Fatal(err)};sid,_:=res.LastInsertId();_,err=db.SQL.Exec(`INSERT INTO alerts(station_id,sample_id,status,severity,opened_at) VALUES(?,?,?,?,?)`,station,sid,"open","high",time.Now().UTC().Format(time.RFC3339Nano));if err!=nil{t.Fatal(err)};removed,err:=repository.Maintenance{DB:db.SQL}.DeleteOldSamples(context.Background(),time.Now().Add(-24*time.Hour));if err!=nil{t.Fatal(err)};if removed!=0{t.Fatalf("referenced sample removed: %d",removed)}
}

var _ *sql.Tx
var _ = errors.New
var _ = time.Now
var _ = context.Background
var _ = sqlite.WithTx
var _ = service.Auth{}
var _ = repository.Users{}
var _ = worker.Job(nil)
