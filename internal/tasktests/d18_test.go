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

func TestD18EventDeliveryDuplicateEvidence(t *testing.T) {
    db:=testsupport.Open(t);bus:=service.EventBus{Repo:repository.Events{DB:db.SQL}};user:=domain.User{ID:1,Role:domain.RoleAdmin,Region:"Lanzhou"};event:=domain.Event{ID:"D18-event",Type:domain.EventWarningSent,ObjectType:"warning",ObjectID:"42",OccurredAt:time.Now(),Payload:map[string]any{"x":1}};if err:=bus.Publish(context.Background(),user,event);err!=nil{t.Fatal(err)};if err:=bus.DeliverPending(context.Background(),10);err!=nil{t.Fatal(err)};count,err:=bus.PendingCount(context.Background());if err!=nil||count!=0{t.Fatalf("pending=%d err=%v",count,err)}
}

var _ *sql.Tx
var _ = errors.New
var _ = time.Now
var _ = context.Background
var _ = sqlite.WithTx
var _ = service.Auth{}
var _ = repository.Users{}
var _ = worker.Job(nil)
