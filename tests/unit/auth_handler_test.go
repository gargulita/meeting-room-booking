package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"booking-service/internal/handlers"

	"github.com/stretchr/testify/require"
)

type dummyLoginResp struct {
	Token string `json:"token"`
}

func TestDummyLogin_InvalidJSON(t *testing.T) {
	authHandler := handlers.NewAuthHandler(nil, testJWTSecret)

	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBufferString(`{"role":`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	authHandler.DummyLogin(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestDummyLogin_InvalidRole(t *testing.T) {
	authHandler := handlers.NewAuthHandler(nil, testJWTSecret)

	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBufferString(`{"role":"manager"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	authHandler.DummyLogin(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestDummyLogin_MissingRole(t *testing.T) {
	authHandler := handlers.NewAuthHandler(nil, testJWTSecret)

	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	authHandler.DummyLogin(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestRegister_InvalidJSON(t *testing.T) {
	authHandler := handlers.NewAuthHandler(nil, testJWTSecret)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{"email":`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	authHandler.Register(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestRegister_MissingRequiredFields(t *testing.T) {
	authHandler := handlers.NewAuthHandler(nil, testJWTSecret)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{"email":"test@example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	authHandler.Register(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "error")
}

func TestLogin_InvalidJSON(t *testing.T) {
	authHandler := handlers.NewAuthHandler(nil, testJWTSecret)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"email":`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	authHandler.Login(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestLogin_MissingRequiredFields(t *testing.T) {
	authHandler := handlers.NewAuthHandler(nil, testJWTSecret)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"email":"test@example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	authHandler.Login(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "error")
}

func TestDummyLogin_ResponseIsJSONOnInvalidRole(t *testing.T) {
	authHandler := handlers.NewAuthHandler(nil, testJWTSecret)

	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBufferString(`{"role":"invalid"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	authHandler.DummyLogin(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var payload map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &payload)
	require.NoError(t, err)
	require.Contains(t, payload, "error")
}
