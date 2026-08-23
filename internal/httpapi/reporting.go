package httpapi

import (
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/service"
	"net/http"
	"strconv"
)

type ReportingHandler struct {
	Auth      service.Auth
	Reporting service.Reporting
}

func (h ReportingHandler) Stations(w http.ResponseWriter, r *http.Request) {
	user, err := h.Auth.Current(r.Context(), bearer(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	page := domain.Page{Number: atoi(r.URL.Query().Get("page")), Size: atoi(r.URL.Query().Get("size"))}
	rows, err := h.Reporting.StationSummary(r.Context(), user, r.URL.Query().Get("region"), page)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": rows, "page": page})
}
func (h ReportingHandler) Progress(w http.ResponseWriter, r *http.Request) {
	user, err := h.Auth.Current(r.Context(), bearer(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	id := int64(atoi(pathTail(r.URL.Path)))
	progress, err := h.Reporting.Progress(r.Context(), user, id)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"project_id": id, "progress": progress})
}
func atoi(value string) int { n, _ := strconv.Atoi(value); return n }
