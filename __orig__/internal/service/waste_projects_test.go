package service_test

import (
	"context"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/service"
	"github.com/11DingKing/lanzhou-water-governance/internal/testsupport"
	"testing"
)

func TestWasteAndProjects(t *testing.T) {
	db := testsupport.Open(t)
	users := repository.Users{DB: db.SQL}
	admin, _ := users.Create(context.Background(), "admin", "pw", domain.RoleAdmin, "Lanzhou")
	waste := service.Waste{Repo: repository.Waste{DB: db.SQL}, Audit: repository.Audit{DB: db.SQL}}
	m, err := waste.Create(context.Background(), admin, domain.Manifest{Number: "WM-1", ProducerRegion: "Lanzhou", CarrierRegion: "Baiyin", FacilityRegion: "Linxia", WasteType: "sludge", WeightKG: 50})
	if err != nil {
		t.Fatal(err)
	}
	advanced, err := waste.Advance(context.Background(), admin, m.ID, domain.ManifestCreated, domain.ManifestInTransit, m.Version)
	if err != nil || advanced.Status != domain.ManifestInTransit {
		t.Fatalf("%+v %v", advanced, err)
	}
	projects := service.Projects{Repo: repository.Projects{DB: db.SQL}, Audit: repository.Audit{DB: db.SQL}}
	p, err := projects.Create(context.Background(), admin, domain.Project{Name: "riverbank", Region: "Lanzhou", TargetHectares: 3, BudgetCents: 1000})
	if err != nil {
		t.Fatal(err)
	}
	started, err := projects.Start(context.Background(), admin, p.ID, p.Version)
	if err != nil || started.Status != domain.ProjectBuilding {
		t.Fatalf("%+v %v", started, err)
	}
}
