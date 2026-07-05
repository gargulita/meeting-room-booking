package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"booking-service/internal/domain"
	"booking-service/internal/repository"
	"booking-service/internal/utils"

	"github.com/google/uuid"
)

type AuthHandler struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthHandler(userRepo repository.UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, jwtSecret: jwtSecret}
}

type DummyLoginRequest struct {
	Role string `json:"role"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "email and password are required")
		return
	}
	if req.Role != string(domain.RoleAdmin) && req.Role != string(domain.RoleUser) {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid role")
		return
	}

	existing, err := h.userRepo.GetByEmail(r.Context(), req.Email)
	if err != nil {
		internalError(w, err)
		return
	}
	if existing != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "email already exists")
		return
	}

	hash := hashPassword(req.Password)
	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Role:         domain.UserRole(req.Role),
		PasswordHash: &hash,
	}
	if err := h.userRepo.Create(r.Context(), user); err != nil {
		internalError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"user": user})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "email and password are required")
		return
	}

	user, err := h.userRepo.GetByEmail(r.Context(), req.Email)
	if err != nil {
		internalError(w, err)
		return
	}
	if user == nil || user.PasswordHash == nil || *user.PasswordHash != hashPassword(req.Password) {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid credentials")
		return
	}

	token, err := utils.GenerateJWT(user.ID, string(user.Role), 24*time.Hour, h.jwtSecret)
	if err != nil {
		internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *AuthHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var req DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}

	var userID uuid.UUID
	var role string
	switch req.Role {
	case "admin":
		userID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
		role = "admin"
	case "user":
		userID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
		role = "user"
	default:
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid role")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		internalError(w, err)
		return
	}
	if user == nil {
		user = &domain.User{ID: userID, Email: role + "@example.com", Role: domain.UserRole(role)}
		if err := h.userRepo.Create(r.Context(), user); err != nil {
			internalError(w, err)
			return
		}
	}

	token, err := utils.GenerateJWT(userID, role, 24*time.Hour, h.jwtSecret)
	if err != nil {
		internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}
