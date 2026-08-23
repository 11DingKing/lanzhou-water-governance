package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/service"
)

type MonitoringHandler struct {
	Auth       service.Auth
	Monitoring service.Monitoring
}

func (h MonitoringHandler) Sample(w http.ResponseWriter, r *http.Request) {
	user, err := h.current(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	var in struct {
		StationID int64              `json:"station_id"`
		SampledAt string             `json:"sampled_at"`
		Quality   string             `json:"quality_class"`
		Metrics   map[string]float64 `json:"metrics"`
	}
	if err = json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, r, err)
		return
	}
	sampled := time.Now().UTC()
	if in.SampledAt != "" {
		sampled, _ = time.Parse(time.RFC3339, in.SampledAt)
	}
	sample, alert, err := h.Monitoring.RecordSample(r.Context(), user, domain.Sample{StationID: in.StationID, SampledAt: sampled, Quality: domain.QualityClass(in.Quality), Metrics: in.Metrics})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"sample": sample, "alert": alert})
}
func (h MonitoringHandler) StartInspection(w http.ResponseWriter, r *http.Request) {
	user, err := h.current(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	id, _ := strconv.ParseInt(pathTail(r.URL.Path), 10, 64)
	version, _ := strconv.ParseInt(r.URL.Query().Get("version"), 10, 64)
	inspection, err := h.Monitoring.StartInspection(r.Context(), user, id, version)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, inspection)
}
func (h MonitoringHandler) CompleteInspection(w http.ResponseWriter, r *http.Request) {
	user, err := h.current(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	id, _ := strconv.ParseInt(pathTail(r.URL.Path), 10, 64)
	version, _ := strconv.ParseInt(r.URL.Query().Get("version"), 10, 64)
	var in struct {
		Notes string `json:"notes"`
	}
	_ = json.NewDecoder(r.Body).Decode(&in)
	inspection, err := h.Monitoring.CompleteInspection(r.Context(), user, id, version, in.Notes)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, inspection)
}
func (h MonitoringHandler) current(r *http.Request) (domain.User, error) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	return h.Auth.Current(r.Context(), token)
}
func pathTail(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return "0"
	}
	return parts[len(parts)-1]
}
