package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/service"
)

type WasteHandler struct {
	Auth  service.Auth
	Waste service.Waste
}

func (h WasteHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := h.Auth.Current(r.Context(), bearer(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	var m domain.Manifest
	if err = json.NewDecoder(r.Body).Decode(&m); err != nil {
		writeError(w, r, err)
		return
	}
	created, err := h.Waste.Create(r.Context(), user, m)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}
func (h WasteHandler) Advance(w http.ResponseWriter, r *http.Request) {
	user, err := h.Auth.Current(r.Context(), bearer(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	id, _ := strconv.ParseInt(pathTail(r.URL.Path), 10, 64)
	from := domain.ManifestStatus(r.URL.Query().Get("from"))
	to := domain.ManifestStatus(r.URL.Query().Get("to"))
	version, _ := strconv.ParseInt(r.URL.Query().Get("version"), 10, 64)
	m, err := h.Waste.Advance(r.Context(), user, id, from, to, version)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, m)
}
