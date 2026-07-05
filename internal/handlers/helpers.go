package handlers

import (
	"net/http"

	"booking-service/internal/utils"
)

func requireRole(w http.ResponseWriter, r *http.Request, expected string) bool {
	role, err := utils.GetUserRoleFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return false
	}
	if role != expected {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return false
	}
	return true
}
