package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/service"
)

type ProjectsHandler struct {
	Auth     service.Auth
	Projects service.Projects
}

func (h ProjectsHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := h.Auth.Current(r.Context(), bearer(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	var p domain.Project
	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, r, err)
		return
	}
	created, err := h.Projects.Create(r.Context(), user, p)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}
func (h ProjectsHandler) Transition(w http.ResponseWriter, r *http.Request) {
	user, err := h.Auth.Current(r.Context(), bearer(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	id, _ := strconv.ParseInt(pathTail(r.URL.Path), 10, 64)
	version, _ := strconv.ParseInt(r.URL.Query().Get("version"), 10, 64)
	var p domain.Project
	if r.URL.Query().Get("to") == string(domain.ProjectAccepted) {
		p, err = h.Projects.Accept(r.Context(), user, id, version)
	} else {
		p, err = h.Projects.Start(r.Context(), user, id, version)
	}
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}
