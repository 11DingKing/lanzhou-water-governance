package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/11DingKing/lanzhou-water-governance/internal/service"
)

type CollaborationHandler struct {
	Auth          service.Auth
	Collaboration service.Collaboration
}

func (h CollaborationHandler) Warning(w http.ResponseWriter, r *http.Request) {
	user, err := h.Auth.Current(r.Context(), bearer(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	var in struct {
		Upstream, Downstream string
		StationID            int64
		Payload              map[string]any
	}
	if err = json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, r, err)
		return
	}
	warning, err := h.Collaboration.IssueWarning(r.Context(), user, in.Upstream, in.Downstream, in.StationID, in.Payload)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, warning)
}
func (h CollaborationHandler) Compensation(w http.ResponseWriter, r *http.Request) {
	user, err := h.Auth.Current(r.Context(), bearer(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	var in struct {
		Upstream, Downstream, Period, Direction, Reason string
		AmountCents                                     int64
	}
	if err = json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, r, err)
		return
	}
	c, err := h.Collaboration.CalculateCompensation(r.Context(), user, in.Upstream, in.Downstream, in.Period, in.Direction, in.Reason, in.AmountCents)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, c)
}
func bearer(r *http.Request) string {
	value := r.Header.Get("Authorization")
	if len(value) > 7 {
		return value[7:]
	}
	return value
}
