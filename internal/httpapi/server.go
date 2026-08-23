package httpapi

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/middleware"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/service"
)

type Dependencies struct {
	DB            *sql.DB
	Auth          service.Auth
	Monitoring    service.Monitoring
	Collaboration service.Collaboration
	Waste         service.Waste
	Projects      service.Projects
	Logger        *slog.Logger
}

func NewRouter(d Dependencies) http.Handler {
	mux := http.NewServeMux()
	auth := AuthHandler{Auth: d.Auth}
	monitor := MonitoringHandler{Auth: d.Auth, Monitoring: d.Monitoring}
	collab := CollaborationHandler{Auth: d.Auth, Collaboration: d.Collaboration}
	waste := WasteHandler{Auth: d.Auth, Waste: d.Waste}
	projects := ProjectsHandler{Auth: d.Auth, Projects: d.Projects}
	reporting := ReportingHandler{Auth: d.Auth, Reporting: service.Reporting{Repo: repository.Reporting{DB: d.DB}}}
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := d.DB.PingContext(r.Context()); err != nil {
			writeError(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})
	mux.HandleFunc("POST /api/auth/register", auth.Register)
	mux.HandleFunc("POST /api/auth/login", auth.Login)
	mux.HandleFunc("POST /api/auth/logout", auth.Logout)
	mux.HandleFunc("POST /api/monitoring/samples", monitor.Sample)
	mux.HandleFunc("POST /api/inspections/", monitor.StartInspection)
	mux.HandleFunc("PATCH /api/inspections/", monitor.CompleteInspection)
	mux.HandleFunc("POST /api/collaboration/warnings", collab.Warning)
	mux.HandleFunc("POST /api/collaboration/compensations", collab.Compensation)
	mux.HandleFunc("POST /api/waste/manifests", waste.Create)
	mux.HandleFunc("PATCH /api/waste/manifests/", waste.Advance)
	mux.HandleFunc("POST /api/projects", projects.Create)
	mux.HandleFunc("PATCH /api/projects/", projects.Transition)
	mux.HandleFunc("GET /api/reports/stations", reporting.Stations)
	mux.HandleFunc("GET /api/reports/projects/", reporting.Progress)
	handler := middleware.RequestID(middleware.Recovery(d.Logger, middleware.Timeout(mux, 15*time.Second)))
	return handler
}
