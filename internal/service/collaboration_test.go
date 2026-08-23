package service_test

import (
	"context"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/service"
	"github.com/11DingKing/lanzhou-water-governance/internal/testsupport"
	"testing"
)

func TestCollaborationWritesWarningAndCompensationAtomically(t *testing.T) {
	db := testsupport.Open(t)
	users := repository.Users{DB: db.SQL}
	regional, _ := users.Create(context.Background(), "regional", "pw", domain.RoleRegional, "Lanzhou")
	testsupport.SeedAgreement(t, db.SQL, "Lanzhou", "Baiyin")
	stationID := testsupport.SeedStation(t, db.SQL, "SVC-COL-1", "Lanzhou")
	svc := service.Collaboration{DB: db.SQL, Repo: repository.Collaboration{DB: db.SQL}, Audit: repository.Audit{DB: db.SQL}}
	warning, err := svc.IssueWarning(context.Background(), regional, "Lanzhou", "Baiyin", stationID, map[string]any{"class": "III"})
	if err != nil || warning.ID == 0 {
		t.Fatalf("%+v %v", warning, err)
	}
	compensation, err := svc.CalculateCompensation(context.Background(), regional, "Lanzhou", "Baiyin", "2026-08", "upstream", "改善断面", 10000)
	if err != nil || compensation.AmountCents != 10000 {
		t.Fatalf("%+v %v", compensation, err)
	}
}
