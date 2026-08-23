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

func TestB04ConcurrentInspectionResolution(t *testing.T) {
    db:=testsupport.Open(t); users:=repository.Users{DB:db.SQL}; owner,_:=users.Create(context.Background(),"b04","pw",domain.RoleInspector,"Lanzhou"); station:=testsupport.SeedStation(t,db.SQL,"B04-1","Lanzhou"); var alertID,inspectionID int64; err:=sqlite.WithTx(context.Background(),db.SQL,func(tx *sql.Tx)error{res,e:=tx.Exec(`INSERT INTO alerts(station_id,status,severity,opened_at) VALUES(?,?,?,?)`,station,"investigating","high",time.Now().UTC().Format(time.RFC3339Nano));if e!=nil{return e};alertID,_=res.LastInsertId();res,e=tx.Exec(`INSERT INTO inspections(alert_id,owner_id,status,due_at) VALUES(?,?,?,?,?)`,alertID,owner.ID,"running",time.Now().Add(time.Hour).Format(time.RFC3339Nano));if e!=nil{return e};inspectionID,_=res.LastInsertId();return nil});if err!=nil{t.Fatal(err)}; repo:=repository.Inspections{DB:db.SQL}; start:=make(chan struct{}); results:=make(chan error,2); for i:=0;i<2;i++{go func(){<-start; _,e:=repo.Transition(context.Background(),inspectionID,domain.InspectionRunning,domain.InspectionCompleted,1,"done");results<-e}()};close(start); success:=0;for i:=0;i<2;i++{if <-results==nil{success++}};if success!=1{t.Fatalf("successful transitions=%d",success)}
}

var _ *sql.Tx
var _ = errors.New
var _ = time.Now
var _ = context.Background
var _ = sqlite.WithTx
var _ = service.Auth{}
var _ = repository.Users{}
var _ = worker.Job(nil)
