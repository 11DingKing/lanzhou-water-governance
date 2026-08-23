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

func TestB09ProjectAcceptanceRequiresMilestones(t *testing.T) {
    db:=testsupport.Open(t);users:=repository.Users{DB:db.SQL};admin,_:=users.Create(context.Background(),"b09","pw",domain.RoleAdmin,"Lanzhou");repo:=repository.Projects{DB:db.SQL};p,_:=repo.Create(context.Background(),domain.Project{Name:"B09",Region:"Lanzhou",TargetHectares:2,BudgetCents:100});_,err:=repo.AddMilestone(context.Background(),p.ID,"buffer",time.Now().Add(time.Hour));if err!=nil{t.Fatal(err)};svc:=service.Projects{Repo:repo,Audit:repository.Audit{DB:db.SQL}};p,err=svc.Start(context.Background(),admin,p.ID,p.Version);if err!=nil{t.Fatal(err)};if _,err=svc.Accept(context.Background(),admin,p.ID,p.Version);!errors.Is(err,domain.ErrConflict){t.Fatalf("accepted with pending milestone: %v",err)}
}

var _ *sql.Tx
var _ = errors.New
var _ = time.Now
var _ = context.Background
var _ = sqlite.WithTx
var _ = service.Auth{}
var _ = repository.Users{}
var _ = worker.Job(nil)
