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

func TestB06CompensationIdempotency(t *testing.T) {
    db:=testsupport.Open(t); repo:=repository.Idempotency{DB:db.SQL};if err:=repo.Save(context.Background(),"B06-key","hash-a","old",time.Hour);err!=nil{t.Fatal(err)};if _,err:=repo.Load(context.Background(),"B06-key","hash-b");!errors.Is(err,domain.ErrDuplicate){t.Fatalf("mismatched request accepted: %v",err)}
}

var _ *sql.Tx
var _ = errors.New
var _ = time.Now
var _ = context.Background
var _ = sqlite.WithTx
var _ = service.Auth{}
var _ = repository.Users{}
var _ = worker.Job(nil)
