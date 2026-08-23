package httpapi

import (
	"encoding/json"
	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/service"
	"net/http"
)

type AdminHandler struct {
	Auth     service.Auth
	Projects service.Projects
}

func (h AdminHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	user, err := h.Auth.Current(r.Context(), bearer(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	var project domain.Project
	if err = json.NewDecoder(r.Body).Decode(&project); err != nil {
		writeError(w, r, err)
		return
	}
	created, err := h.Projects.Create(r.Context(), user, project)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}
func (h AdminHandler) StartProject(w http.ResponseWriter, r *http.Request) {
	user, err := h.Auth.Current(r.Context(), bearer(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	id := int64(atoi(pathTail(r.URL.Path)))
	version := int64(atoi(r.URL.Query().Get("version")))
	project, err := h.Projects.Start(r.Context(), user, id, version)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, project)
}
func (h AdminHandler) AcceptProject(w http.ResponseWriter, r *http.Request) {
	user, err := h.Auth.Current(r.Context(), bearer(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	id := int64(atoi(pathTail(r.URL.Path)))
	version := int64(atoi(r.URL.Query().Get("version")))
	project, err := h.Projects.Accept(r.Context(), user, id, version)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, project)
}
