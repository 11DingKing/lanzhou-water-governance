package testsupport

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/storage/sqlite"
)

func Open(t testing.TB) *sqlite.DB {
	t.Helper()
	db, err := sqlite.Open(context.Background(), ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}
func SeedUser(t testing.TB, db *sql.DB, username, role, region string) int64 {
	t.Helper()
	res, err := db.Exec(`INSERT INTO users(username,password_hash,role,region,created_at) VALUES(?,?,?,?,?)`, username, "pw", role, region, time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		t.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	return id
}
func SeedStation(t testing.TB, db *sql.DB, code, region string) int64 {
	t.Helper()
	res, err := db.Exec(`INSERT INTO stations(code,name,region,river,level,created_at) VALUES(?,?,?,?,?,?)`, code, "断面"+code, region, "黄河", "国控", time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		t.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	return id
}
func SeedAgreement(t testing.TB, db *sql.DB, up, down string) int64 {
	t.Helper()
	res, err := db.Exec(`INSERT INTO agreements(upstream_region,downstream_region,threshold_class,created_at) VALUES(?,?,?,?)`, up, down, "II", time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		t.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	return id
}
func MustExec(t testing.TB, db *sql.DB, query string, args ...any) {
	t.Helper()
	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("%s: %v", fmt.Sprintf(query, args...), err)
	}
}

var _ = repository.Audit{}
