package unit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"booking-service/internal/handlers"
	"booking-service/internal/middleware"
	"booking-service/internal/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

const testJWTSecret = "unit-test-secret"

func makeTestToken(t *testing.T, role string) string {
	t.Helper()

	token, err := utils.GenerateJWT(uuid.New(), role, time.Hour, testJWTSecret)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	return token
}

func TestGenerateAndValidateJWT_User(t *testing.T) {
	userID := uuid.New()

	token, err := utils.GenerateJWT(userID, "user", time.Hour, testJWTSecret)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := utils.ValidateJWT(token, testJWTSecret)
	require.NoError(t, err)
	require.NotNil(t, claims)
	require.Equal(t, userID, claims.UserID)
	require.Equal(t, "user", claims.Role)
}

func TestGenerateAndValidateJWT_Admin(t *testing.T) {
	userID := uuid.New()

	token, err := utils.GenerateJWT(userID, "admin", time.Hour, testJWTSecret)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := utils.ValidateJWT(token, testJWTSecret)
	require.NoError(t, err)
	require.NotNil(t, claims)
	require.Equal(t, userID, claims.UserID)
	require.Equal(t, "admin", claims.Role)
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	claims, err := utils.ValidateJWT("definitely-not-a-valid-token", testJWTSecret)
	require.Error(t, err)
	require.Nil(t, claims)
}

func TestGetUserIDFromContext_Success(t *testing.T) {
	expectedID := uuid.New()
	ctxKey := utils.CtxUserIDKey

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), ctxKey, expectedID))

	userID, err := utils.GetUserIDFromContext(req.Context())
	require.NoError(t, err)
	require.Equal(t, expectedID, userID)
}

func TestGetUserIDFromContext_Missing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	userID, err := utils.GetUserIDFromContext(req.Context())
	require.Error(t, err)
	require.Equal(t, uuid.Nil, userID)
}

func TestGetUserRoleFromContext_Success(t *testing.T) {
	expectedRole := "user"
	ctxKey := utils.CtxUserRoleKey

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), ctxKey, expectedRole))

	role, err := utils.GetUserRoleFromContext(req.Context())
	require.NoError(t, err)
	require.Equal(t, expectedRole, role)
}

func TestGetUserRoleFromContext_Missing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	role, err := utils.GetUserRoleFromContext(req.Context())
	require.Error(t, err)
	require.Empty(t, role)
}

func TestJWTMiddleware_MissingAuthorizationHeader(t *testing.T) {
	router := mux.NewRouter()
	api := router.PathPrefix("/").Subrouter()
	api.Use(middleware.JWTAuth(testJWTSecret))
	api.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	router := mux.NewRouter()
	api := router.PathPrefix("/").Subrouter()
	api.Use(middleware.JWTAuth(testJWTSecret))
	api.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer definitely-not-a-valid-token")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTMiddleware_ValidUserToken(t *testing.T) {
	token := makeTestToken(t, "user")

	router := mux.NewRouter()
	api := router.PathPrefix("/").Subrouter()
	api.Use(middleware.JWTAuth(testJWTSecret))
	api.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "ok", rec.Body.String())
}

func TestJWTMiddleware_ValidAdminToken(t *testing.T) {
	token := makeTestToken(t, "admin")

	router := mux.NewRouter()
	api := router.PathPrefix("/").Subrouter()
	api.Use(middleware.JWTAuth(testJWTSecret))
	api.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "ok", rec.Body.String())
}

func TestInfoHandler(t *testing.T) {
	infoHandler := handlers.NewInfoHandler()

	req := httptest.NewRequest(http.MethodGet, "/_info", nil)
	rec := httptest.NewRecorder()

	infoHandler.Info(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}
