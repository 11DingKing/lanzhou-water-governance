package sqlite_test

import (
	"context"
	"github.com/11DingKing/lanzhou-water-governance/internal/storage/sqlite"
	"testing"
)

func TestMigrationsCreateRelations(t *testing.T) {
	db, err := sqlite.Open(context.Background(), ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var count int
	if err = db.SQL.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name IN ('users','stations','samples','alerts','inspections','manifests','projects','audit_events')`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 8 {
		t.Fatalf("tables=%d", count)
	}
}
func TestMigrationIsIdempotent(t *testing.T) {
	db, err := sqlite.Open(context.Background(), ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}
	var versions int
	if err = db.SQL.QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&versions); err != nil {
		t.Fatal(err)
	}
	if versions != 1 {
		t.Fatalf("versions=%d", versions)
	}
}
