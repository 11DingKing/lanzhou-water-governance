package repository_test

import (
	"context"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/testsupport"
	"testing"
)

func TestManifestCustody(t *testing.T) {
	db := testsupport.Open(t)
	repo := repository.Waste{DB: db.SQL}
	manifest, err := repo.Create(context.Background(), domain.Manifest{Number: "M-1", ProducerRegion: "Lanzhou", CarrierRegion: "Baiyin", FacilityRegion: "Linxia", WasteType: "sludge", WeightKG: 300})
	if err != nil {
		t.Fatal(err)
	}
	next, err := repo.Transition(context.Background(), manifest.ID, domain.ManifestCreated, domain.ManifestInTransit, manifest.Version)
	if err != nil || next.Status != domain.ManifestInTransit {
		t.Fatalf("%+v %v", next, err)
	}
	next, err = repo.Transition(context.Background(), manifest.ID, domain.ManifestInTransit, domain.ManifestAccepted, next.Version)
	if err != nil || next.AcceptedAt == nil {
		t.Fatalf("%+v %v", next, err)
	}
	next, err = repo.Transition(context.Background(), manifest.ID, domain.ManifestAccepted, domain.ManifestDisposed, next.Version)
	if err != nil || next.DisposedAt == nil {
		t.Fatalf("%+v %v", next, err)
	}
}
func TestManifestConflict(t *testing.T) {
	db := testsupport.Open(t)
	repo := repository.Waste{DB: db.SQL}
	manifest, _ := repo.Create(context.Background(), domain.Manifest{Number: "M-2", ProducerRegion: "Lanzhou", CarrierRegion: "Baiyin", FacilityRegion: "Linxia", WasteType: "ash", WeightKG: 100})
	if _, err := repo.Transition(context.Background(), manifest.ID, domain.ManifestCreated, domain.ManifestDisposed, manifest.Version); err == nil {
		t.Fatal("illegal transition accepted")
	}
}
