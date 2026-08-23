package service_test

import (
	"context"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/service"
	"github.com/11DingKing/lanzhou-water-governance/internal/testsupport"
	"testing"
	"time"
)

func TestAuthLifecycle(t *testing.T) {
	db := testsupport.Open(t)
	auth := service.NewAuth(repository.Users{DB: db.SQL}, time.Hour)
	user, err := auth.CreateUser(context.Background(), "operator", "secret", domain.RoleInspector, "Lanzhou")
	if err != nil {
		t.Fatal(err)
	}
	token, logged, err := auth.Login(context.Background(), "operator", "secret")
	if err != nil || logged.ID != user.ID {
		t.Fatalf("%+v %v", logged, err)
	}
	current, err := auth.Current(context.Background(), token)
	if err != nil || current.Username != "operator" {
		t.Fatalf("%+v %v", current, err)
	}
	if err = auth.Logout(context.Background(), token); err != nil {
		t.Fatal(err)
	}
	if _, err = auth.Current(context.Background(), token); err == nil {
		t.Fatal("logged out token accepted")
	}
}
func TestAuthRejectsWrongPassword(t *testing.T) {
	db := testsupport.Open(t)
	auth := service.NewAuth(repository.Users{DB: db.SQL}, time.Hour)
	_, _ = auth.CreateUser(context.Background(), "operator", "secret", domain.RoleInspector, "Lanzhou")
	if _, _, err := auth.Login(context.Background(), "operator", "wrong"); err != domain.ErrForbidden {
		t.Fatalf("err=%v", err)
	}
}
