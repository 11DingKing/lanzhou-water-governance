package repository_test

import (
	"context"
	"database/sql"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/storage/sqlite"
	"github.com/11DingKing/lanzhou-water-governance/internal/testsupport"
	"testing"
	"time"
)

func TestMonitoringSampleAndAlert(t *testing.T) {
	db := testsupport.Open(t)
	repo := repository.Monitoring{DB: db.SQL}
	userID := testsupport.SeedUser(t, db.SQL, "sample-user", "inspector", "Lanzhou")
	stationID := testsupport.SeedStation(t, db.SQL, "YL-1", "Lanzhou")
	station, _ := repo.Station(context.Background(), stationID)
	if station.Code != "YL-1" {
		t.Fatal(station)
	}
	sample, err := repo.AddSample(context.Background(), domain.Sample{StationID: stationID, SampledAt: time.Now(), Quality: domain.QualityII, Metrics: map[string]float64{"cod": 3}, CreatedBy: userID})
	if err != nil {
		t.Fatal(err)
	}
	if sample.ID == 0 {
		t.Fatal("sample id missing")
	}
	rows, err := repo.ListSamples(context.Background(), stationID, domain.Page{Number: 1, Size: 10})
	if err != nil || len(rows) != 1 {
		t.Fatalf("%d %v", len(rows), err)
	}
}
func TestMonitoringAlertTransition(t *testing.T) {
	db := testsupport.Open(t)
	repo := repository.Monitoring{DB: db.SQL}
	userID := testsupport.SeedUser(t, db.SQL, "alert-user", "inspector", "Lanzhou")
	station := testsupport.SeedStation(t, db.SQL, "YL-2", "Lanzhou")
	sample, err := repo.AddSample(context.Background(), domain.Sample{StationID: station, SampledAt: time.Now(), Quality: domain.QualityIII, Metrics: map[string]float64{"cod": 12}, CreatedBy: userID})
	if err != nil {
		t.Fatal(err)
	}
	var alert domain.Alert
	err = sqlite.WithTx(context.Background(), db.SQL, func(tx *sql.Tx) error {
		var txErr error
		alert, txErr = repo.OpenAlertTx(context.Background(), tx, station, sample.ID, "high")
		return txErr
	})
	if err != nil || alert.ID == 0 {
		t.Fatalf("%+v %v", alert, err)
	}
	opened, _ := repo.Alert(context.Background(), alert.ID)
	if opened.Status != domain.AlertOpen {
		t.Fatal(opened.Status)
	}
}
