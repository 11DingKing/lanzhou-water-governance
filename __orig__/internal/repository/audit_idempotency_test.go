package repository_test

import (
	"context"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/testsupport"
	"testing"
	"time"
)

func TestAuditAndIdempotency(t *testing.T) {
	db := testsupport.Open(t)
	audit := repository.Audit{DB: db.SQL}
	actor := testsupport.SeedUser(t, db.SQL, "audit-user", "admin", "Lanzhou")
	if err := audit.Record(context.Background(), actor, "req-1", "station", "1", "sample", "ok", map[string]any{"quality": "II"}); err != nil {
		t.Fatal(err)
	}
	count, err := audit.Count(context.Background(), "station", "1")
	if err != nil || count != 1 {
		t.Fatalf("%d %v", count, err)
	}
	idem := repository.Idempotency{DB: db.SQL}
	if err = idem.Save(context.Background(), "key", "hash", "{\"ok\":true}", time.Hour); err != nil {
		t.Fatal(err)
	}
	response, err := idem.Load(context.Background(), "key", "hash")
	if err != nil || response != "{\"ok\":true}" {
		t.Fatalf("%s %v", response, err)
	}
	if _, err = idem.Load(context.Background(), "key", "other"); err == nil {
		t.Fatal("hash mismatch accepted")
	}
}
