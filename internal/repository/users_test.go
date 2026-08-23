package repository_test

import (
	"context"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/testsupport"
	"testing"
	"time"
)

func TestUsersAndSessions(t *testing.T) {
	db := testsupport.Open(t)
	repo := repository.Users{DB: db.SQL}
	user, err := repo.Create(context.Background(), "water-admin", "hash", domain.RoleAdmin, "Lanzhou")
	if err != nil {
		t.Fatal(err)
	}
	found, err := repo.Find(context.Background(), user.Username)
	if err != nil || found.Role != domain.RoleAdmin {
		t.Fatalf("%+v %v", found, err)
	}
	if err = repo.CreateSession(context.Background(), user.ID, "tok", time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	session, err := repo.Session(context.Background(), "tok", time.Now())
	if err != nil || session.ID != user.ID {
		t.Fatalf("%+v %v", session, err)
	}
	if err = repo.RevokeSession(context.Background(), "tok"); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.Session(context.Background(), "tok", time.Now()); err == nil {
		t.Fatal("revoked session accepted")
	}
}
func TestExpiredSession(t *testing.T) {
	db := testsupport.Open(t)
	repo := repository.Users{DB: db.SQL}
	user, _ := repo.Create(context.Background(), "expired", "hash", domain.RoleInspector, "Lanzhou")
	_ = repo.CreateSession(context.Background(), user.ID, "old", time.Now().Add(-time.Minute))
	if _, err := repo.Session(context.Background(), "old", time.Now()); err != domain.ErrExpired {
		t.Fatalf("err=%v", err)
	}
}
func TestDisableUser(t *testing.T) {
	db := testsupport.Open(t)
	repo := repository.Users{DB: db.SQL}
	user, _ := repo.Create(context.Background(), "disabled", "hash", domain.RoleInspector, "Lanzhou")
	if err := repo.Disable(context.Background(), user.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Password(context.Background(), user.Username); err == nil {
		t.Fatal("disabled password returned")
	}
}
