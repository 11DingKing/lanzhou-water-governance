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

func TestB05RevokedSessionCannotSample(t *testing.T) {
    db:=testsupport.Open(t); repo:=repository.Users{DB:db.SQL}; user,_:=repo.Create(context.Background(),"b05","hash",domain.RoleInspector,"Lanzhou"); if err:=repo.CreateSession(context.Background(),user.ID,"B05-special",time.Now().Add(time.Hour));err!=nil{t.Fatal(err)};if err:=repo.RevokeSession(context.Background(),"B05-special");err!=nil{t.Fatal(err)};if _,err:=repo.Session(context.Background(),"B05-special",time.Now());err==nil{t.Fatalf("revoked token remained usable")}
}

var _ *sql.Tx
var _ = errors.New
var _ = time.Now
var _ = context.Background
var _ = sqlite.WithTx
var _ = service.Auth{}
var _ = repository.Users{}
var _ = worker.Job(nil)
