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

func TestB03ManifestTimeline(t *testing.T) {
    db:=testsupport.Open(t); repo:=repository.Waste{DB:db.SQL}; m,err:=repo.Create(context.Background(),domain.Manifest{Number:"B03",ProducerRegion:"Lanzhou",CarrierRegion:"Baiyin",FacilityRegion:"Linxia",WasteType:"sludge",WeightKG:10}); if err!=nil{t.Fatal(err)}; m,err=repo.Transition(context.Background(),m.ID,domain.ManifestCreated,domain.ManifestInTransit,m.Version); if err!=nil{t.Fatal(err)}; m,err=repo.Transition(context.Background(),m.ID,domain.ManifestInTransit,domain.ManifestAccepted,m.Version); if err!=nil{t.Fatal(err)}; if m.AcceptedAt==nil || m.DisposedAt!=nil { t.Fatalf("timeline accepted=%v disposed=%v",m.AcceptedAt,m.DisposedAt) }
}

var _ *sql.Tx
var _ = errors.New
var _ = time.Now
var _ = context.Background
var _ = sqlite.WithTx
var _ = service.Auth{}
var _ = repository.Users{}
var _ = worker.Job(nil)
