package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/11DingKing/lanzhou-water-governance/internal/domain"
	"github.com/11DingKing/lanzhou-water-governance/internal/middleware"
)

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusInternalServerError
	code := "internal_error"
	message := "internal server error"
	switch {
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
		code = "not_found"
		message = "resource not found"
	case errors.Is(err, domain.ErrForbidden):
		status = http.StatusForbidden
		code = "forbidden"
		message = "permission denied"
	case errors.Is(err, domain.ErrConflict):
		status = http.StatusConflict
		code = "conflict"
		message = "state conflict"
	case errors.Is(err, domain.ErrInvalidState):
		status = http.StatusUnprocessableEntity
		code = "invalid_state"
		message = "invalid lifecycle transition"
	case errors.Is(err, domain.ErrExpired):
		status = http.StatusUnauthorized
		code = "expired"
		message = "session expired"
	}
	writeJSON(w, status, map[string]any{"code": code, "message": message, "request_id": middleware.RequestIDFrom(r.Context())})
}
