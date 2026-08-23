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

func TestB11ManifestRegionInvariant(t *testing.T) {
    db:=testsupport.Open(t);users:=repository.Users{DB:db.SQL};user,_:=users.Create(context.Background(),"b11","pw",domain.RoleRegional,"Lanzhou");svc:=service.Waste{Repo:repository.Waste{DB:db.SQL},Audit:repository.Audit{DB:db.SQL}};m,err:=svc.Create(context.Background(),user,domain.Manifest{Number:"B11",ProducerRegion:"Lanzhou",CarrierRegion:"Baiyin",FacilityRegion:"Baiyin",WasteType:"sludge",WeightKG:10});if err!=nil{t.Fatal(err)};if _,err=svc.Advance(context.Background(),user,m.ID,domain.ManifestCreated,domain.ManifestInTransit,m.Version);!errors.Is(err,domain.ErrConflict){t.Fatalf("illegal region transition accepted: %v",err)}
}

var _ *sql.Tx
var _ = errors.New
var _ = time.Now
var _ = context.Background
var _ = sqlite.WithTx
var _ = service.Auth{}
var _ = repository.Users{}
var _ = worker.Job(nil)
