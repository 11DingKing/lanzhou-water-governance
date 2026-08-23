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

func TestD15CancelledJobLeaseEvidence(t *testing.T) {
    _ = domain.RoleAdmin; db:=testsupport.Open(t);reconciler:=worker.Reconciler{DB:db.SQL,Idempotency:repository.Idempotency{DB:db.SQL}};ctx,cancel:=context.WithCancel(context.Background());cancel();if err:=reconciler.Job()(ctx);!errors.Is(err,context.Canceled){t.Fatalf("cancel was lost: %v",err)}
}

var _ *sql.Tx
var _ = errors.New
var _ = time.Now
var _ = context.Background
var _ = sqlite.WithTx
var _ = service.Auth{}
var _ = repository.Users{}
var _ = worker.Job(nil)
