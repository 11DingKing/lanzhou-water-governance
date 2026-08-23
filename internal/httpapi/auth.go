package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/service"
)

type AuthHandler struct{ Auth service.Auth }

func (h AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct{ Username, Password, Role, Region string }
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, r, err)
		return
	}
	user, err := h.Auth.CreateUser(r.Context(), input.Username, input.Password, domain.Role(input.Role), input.Region)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": user.ID, "username": user.Username, "role": user.Role, "region": user.Region})
}
func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct{ Username, Password string }
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, r, err)
		return
	}
	token, user, err := h.Auth.Login(r.Context(), input.Username, input.Password)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "user": user})
}
func (h AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if len(token) > 7 {
		token = token[7:]
	}
	if err := h.Auth.Logout(r.Context(), token); err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
