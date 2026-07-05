package unit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"booking-service/internal/middleware"
	"booking-service/internal/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestJWTValidation_RejectsWrongSecret(t *testing.T) {
	token, err := utils.GenerateJWT(uuid.New(), "user", time.Hour, testJWTSecret)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := utils.ValidateJWT(token, "definitely-wrong-secret")
	require.Error(t, err)
	require.Nil(t, claims)
}

func TestJWTValidation_RejectsExpiredToken(t *testing.T) {
	token, err := utils.GenerateJWT(uuid.New(), "user", -time.Minute, testJWTSecret)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := utils.ValidateJWT(token, testJWTSecret)
	require.Error(t, err)
	require.Nil(t, claims)
}

func TestGetUserIDFromContext_RejectsWrongType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), utils.CtxUserIDKey, "not-a-uuid"))

	userID, err := utils.GetUserIDFromContext(req.Context())
	require.Error(t, err)
	require.Equal(t, uuid.Nil, userID)
}

func TestGetUserRoleFromContext_RejectsWrongType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), utils.CtxUserRoleKey, 123))

	role, err := utils.GetUserRoleFromContext(req.Context())
	require.Error(t, err)
	require.Empty(t, role)
}

func TestJWTMiddleware_PropagatesUserIDAndRoleToContext(t *testing.T) {
	expectedUserID := uuid.New()
	token, err := utils.GenerateJWT(expectedUserID, "admin", time.Hour, testJWTSecret)
	require.NoError(t, err)

	router := mux.NewRouter()
	api := router.PathPrefix("/").Subrouter()
	api.Use(middleware.JWTAuth(testJWTSecret))

	api.HandleFunc("/whoami", func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.GetUserIDFromContext(r.Context())
		require.NoError(t, err)
		require.Equal(t, expectedUserID, userID)

		role, err := utils.GetUserRoleFromContext(r.Context())
		require.NoError(t, err)
		require.Equal(t, "admin", role)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "ok", rec.Body.String())
}

func TestJWTMiddleware_RejectsMalformedAuthorizationHeader(t *testing.T) {
	router := mux.NewRouter()
	api := router.PathPrefix("/").Subrouter()
	api.Use(middleware.JWTAuth(testJWTSecret))

	api.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token something")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTMiddleware_RejectsEmptyBearerToken(t *testing.T) {
	router := mux.NewRouter()
	api := router.PathPrefix("/").Subrouter()
	api.Use(middleware.JWTAuth(testJWTSecret))

	api.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer ")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
