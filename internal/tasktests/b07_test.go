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

func TestB07StableSampleCursor(t *testing.T) {
    db:=testsupport.Open(t);repo:=repository.Monitoring{DB:db.SQL};station:=testsupport.SeedStation(t,db.SQL,"B07-1","Lanzhou");uid:=testsupport.SeedUser(t,db.SQL,"b07","inspector","Lanzhou");for i:=1;i<=3;i++{_,err:=repo.AddSample(context.Background(),domain.Sample{StationID:station,SampledAt:time.Date(2026,8,i,0,0,0,0,time.UTC),Quality:domain.QualityII,Metrics:map[string]float64{"cod":float64(i)},CreatedBy:uid});if err!=nil{t.Fatal(err)}};rows,err:=repo.ListSamples(context.Background(),station,domain.Page{Number:1,Size:3});if err!=nil{t.Fatal(err)};if len(rows)!=3||!rows[0].SampledAt.After(rows[1].SampledAt)||!rows[1].SampledAt.After(rows[2].SampledAt){t.Fatalf("order=%v",rows)}
}

var _ *sql.Tx
var _ = errors.New
var _ = time.Now
var _ = context.Background
var _ = sqlite.WithTx
var _ = service.Auth{}
var _ = repository.Users{}
var _ = worker.Job(nil)
