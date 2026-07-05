package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"booking-service/internal/domain"
	"booking-service/internal/service"
)

type RoomHandler struct {
	roomService *service.RoomService
}

func NewRoomHandler(roomService *service.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
}

type CreateRoomRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity"`
}

func roomToResponse(room *domain.Room) map[string]interface{} {
	var createdAt interface{}
	if !room.CreatedAt.IsZero() {
		createdAt = room.CreatedAt.UTC().Format(time.RFC3339)
	}

	return map[string]interface{}{
		"id":          room.ID.String(),
		"name":        room.Name,
		"description": room.Description,
		"capacity":    room.Capacity,
		"createdAt":   createdAt,
	}
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, "admin") {
		return
	}

	var req CreateRoomRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "name is required")
		return
	}

	room := &domain.Room{
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
	}

	if err := h.roomService.CreateRoom(r.Context(), room); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"room": roomToResponse(room),
	})
}

func (h *RoomHandler) ListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.roomService.ListRooms(r.Context())
	if err != nil {
		internalError(w, err)
		return
	}

	result := make([]map[string]interface{}, 0, len(rooms))
	for _, room := range rooms {
		result = append(result, roomToResponse(room))
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"rooms": result,
	})
}
