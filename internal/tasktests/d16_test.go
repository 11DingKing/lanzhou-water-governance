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

func TestD16StationSummaryFilterConsistency(t *testing.T) {
    db:=testsupport.Open(t);station:=testsupport.SeedStation(t,db.SQL,"D16-1","Lanzhou");uid:=testsupport.SeedUser(t,db.SQL,"d16","inspector","Lanzhou");repo:=repository.Monitoring{DB:db.SQL};for i:=1;i<=2;i++{_,err:=repo.AddSample(context.Background(),domain.Sample{StationID:station,SampledAt:time.Now().Add(time.Duration(i)*time.Minute),Quality:domain.QualityII,Metrics:map[string]float64{"cod":1},CreatedBy:uid});if err!=nil{t.Fatal(err)}};_,err:=db.SQL.Exec(`INSERT INTO alerts(station_id,status,severity,opened_at) VALUES(?,?,?,?)`,station,"open","high",time.Now().UTC().Format(time.RFC3339Nano));if err!=nil{t.Fatal(err)};rows,err:=repository.Reporting{DB:db.SQL}.StationSummary(context.Background(),"Lanzhou",domain.Page{Number:1,Size:10});if err!=nil{t.Fatal(err)};if len(rows)!=1||rows[0].OpenAlerts!=1{t.Fatalf("summary=%v",rows)}
}

var _ *sql.Tx
var _ = errors.New
var _ = time.Now
var _ = context.Background
var _ = sqlite.WithTx
var _ = service.Auth{}
var _ = repository.Users{}
var _ = worker.Job(nil)
