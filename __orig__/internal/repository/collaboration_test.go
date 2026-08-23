package repository_test

import (
	"context"
	"database/sql"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/storage/sqlite"
	"github.com/11DingKing/lanzhou-water-governance/internal/testsupport"
	"testing"
)

func TestAgreementWarningAndCompensation(t *testing.T) {
	db := testsupport.Open(t)
	repo := repository.Collaboration{DB: db.SQL}
	id := testsupport.SeedAgreement(t, db.SQL, "Lanzhou", "Baiyin")
	stationID := testsupport.SeedStation(t, db.SQL, "COL-1", "Lanzhou")
	agreement, err := repo.Agreement(context.Background(), "Lanzhou", "Baiyin")
	if err != nil || agreement.ID != id {
		t.Fatalf("%+v %v", agreement, err)
	}
	var warning domain.Warning
	err = sqlite.WithTx(context.Background(), db.SQL, func(tx *sql.Tx) error {
		var txErr error
		warning, txErr = repo.CreateWarningTx(context.Background(), tx, agreement, stationID, "downstream", map[string]any{"quality": "III"})
		return txErr
	})
	if err != nil || warning.ID == 0 {
		t.Fatalf("%+v %v", warning, err)
	}
	var compensation domain.Compensation
	err = sqlite.WithTx(context.Background(), db.SQL, func(tx *sql.Tx) error {
		var txErr error
		compensation, txErr = repo.CreateCompensationTx(context.Background(), tx, agreement, "2026-08", "upstream", "quality improvement", 12000)
		return txErr
	})
	if err != nil || compensation.AmountCents != 12000 {
		t.Fatalf("%+v %v", compensation, err)
	}
	settled, err := repo.SettleCompensation(context.Background(), compensation.ID)
	if err != nil || settled.Status != "settled" {
		t.Fatalf("%+v %v", settled, err)
	}
}
