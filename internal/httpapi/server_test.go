package httpapi_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/11DingKing/lanzhou-water-governance/internal/httpapi"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/service"
	"github.com/11DingKing/lanzhou-water-governance/internal/testsupport"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthAndReady(t *testing.T) {
	db := testsupport.Open(t)
	router := routerFor(db.SQL)
	for _, path := range []string{"/healthz", "/readyz"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("%s code=%d", path, response.Code)
		}
		if response.Header().Get("X-Request-ID") == "" {
			t.Fatal("request id missing")
		}
	}
}
func TestRegisterLoginLogout(t *testing.T) {
	db := testsupport.Open(t)
	router := routerFor(db.SQL)
	body, _ := json.Marshal(map[string]string{"username": "operator", "password": "secret", "role": "inspector", "region": "Lanzhou"})
	register := httptest.NewRecorder()
	router.ServeHTTP(register, httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body)))
	if register.Code != http.StatusCreated {
		t.Fatalf("register=%d %s", register.Code, register.Body.String())
	}
	login := httptest.NewRecorder()
	router.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body)))
	if login.Code != http.StatusOK {
		t.Fatalf("login=%d", login.Code)
	}
	var result struct {
		Token string `json:"token"`
	}
	_ = json.NewDecoder(login.Body).Decode(&result)
	if result.Token == "" {
		t.Fatal("token missing")
	}
	logoutRequest := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	logoutRequest.Header.Set("Authorization", "Bearer "+result.Token)
	logout := httptest.NewRecorder()
	router.ServeHTTP(logout, logoutRequest)
	if logout.Code != http.StatusNoContent {
		t.Fatalf("logout=%d", logout.Code)
	}
}
func TestMalformedJSONMapsToError(t *testing.T) {
	db := testsupport.Open(t)
	router := routerFor(db.SQL)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString("{")))
	if response.Code != http.StatusInternalServerError && response.Code != http.StatusBadRequest {
		t.Fatalf("code=%d", response.Code)
	}
	if _, err := io.ReadAll(response.Body); err != nil {
		t.Fatal(err)
	}
}
func routerFor(db *sql.DB) http.Handler {
	users := repository.Users{DB: db}
	audit := repository.Audit{DB: db}
	auth := service.NewAuth(users, time.Hour)
	monitoring := service.Monitoring{DB: db, Repo: repository.Monitoring{DB: db}, Inspections: repository.Inspections{DB: db}, Audit: audit}
	return httpapi.NewRouter(httpapi.Dependencies{DB: db, Auth: auth, Monitoring: monitoring, Collaboration: service.Collaboration{DB: db, Repo: repository.Collaboration{DB: db}, Audit: audit}, Waste: service.Waste{Repo: repository.Waste{DB: db}, Audit: audit}, Projects: service.Projects{Repo: repository.Projects{DB: db}, Audit: audit}, Logger: slog.Default()})
}

var _ = context.Background
