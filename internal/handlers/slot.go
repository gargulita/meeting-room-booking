package handlers

import (
	"errors"
	"net/http"
	"time"

	"booking-service/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SlotHandler struct {
	slotService *service.SlotService
	roomService *service.RoomService
}

func NewSlotHandler(slotService *service.SlotService, roomService *service.RoomService) *SlotHandler {
	return &SlotHandler{
		slotService: slotService,
		roomService: roomService,
	}
}

func (h *SlotHandler) ListAvailableSlots(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "date is required")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid date")
		return
	}

	roomID, err := uuid.Parse(mux.Vars(r)["roomId"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid roomId")
		return
	}

	_, err = h.roomService.GetRoom(r.Context(), roomID)
	if err != nil {
		if errors.Is(err, service.ErrRoomNotFound) {
			writeError(w, http.StatusNotFound, "ROOM_NOT_FOUND", "room not found")
			return
		}
		internalError(w, err)
		return
	}

	slots, err := h.slotService.GetAvailableSlots(r.Context(), roomID, date)
	if err != nil {
		internalError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, slotsResponse{
		Slots: toSlotDTOs(slots),
	})
}
