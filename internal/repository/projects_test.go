package repository_test

import (
	"context"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/testsupport"
	"testing"
	"time"
)

func TestProjectMilestones(t *testing.T) {
	db := testsupport.Open(t)
	repo := repository.Projects{DB: db.SQL}
	p, err := repo.Create(context.Background(), domain.Project{Name: "buffer zone", Region: "Lanzhou", TargetHectares: 12.5, BudgetCents: 500000})
	if err != nil {
		t.Fatal(err)
	}
	m, err := repo.AddMilestone(context.Background(), p.ID, "planting", time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.CompleteMilestone(context.Background(), m.ID, "photos"); err != nil {
		t.Fatal(err)
	}
	started, err := repo.Transition(context.Background(), p.ID, domain.ProjectPlanned, domain.ProjectBuilding, p.Version)
	if err != nil || started.Status != domain.ProjectBuilding {
		t.Fatalf("%+v %v", started, err)
	}
	accepted, err := repo.Transition(context.Background(), p.ID, domain.ProjectBuilding, domain.ProjectAccepted, started.Version)
	if err != nil || accepted.Status != domain.ProjectAccepted {
		t.Fatalf("%+v %v", accepted, err)
	}
}
func TestProjectVersionConflict(t *testing.T) {
	db := testsupport.Open(t)
	repo := repository.Projects{DB: db.SQL}
	p, _ := repo.Create(context.Background(), domain.Project{Name: "wetland", Region: "Lanzhou", TargetHectares: 4, BudgetCents: 100})
	if _, err := repo.Transition(context.Background(), p.ID, domain.ProjectPlanned, domain.ProjectBuilding, p.Version+2); err != domain.ErrConflict {
		t.Fatalf("err=%v", err)
	}
}
