package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"booking-service/internal/domain"
	"booking-service/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ScheduleHandler struct {
	scheduleService *service.ScheduleService
}

func NewScheduleHandler(scheduleService *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{scheduleService: scheduleService}
}

type CreateScheduleRequest struct {
	ID         string `json:"id,omitempty"`
	RoomID     string `json:"roomId"`
	DaysOfWeek []int  `json:"daysOfWeek"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

func (h *ScheduleHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, "admin") {
		return
	}

	pathRoomID, err := uuid.Parse(mux.Vars(r)["roomId"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid roomId")
		return
	}

	var req CreateScheduleRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}

	if req.RoomID == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "roomId is required")
		return
	}

	bodyRoomID, err := uuid.Parse(req.RoomID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid roomId")
		return
	}

	if bodyRoomID != pathRoomID {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "roomId in body must match roomId in path")
		return
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid startTime")
		return
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid endTime")
		return
	}

	baseDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	schedule := &domain.Schedule{
		RoomID:     pathRoomID,
		DaysOfWeek: req.DaysOfWeek,
		StartTime:  time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.UTC),
		EndTime:    time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.UTC),
	}

	if err := h.scheduleService.CreateSchedule(r.Context(), schedule); err != nil {
		switch err.Error() {
		case "room not found":
			writeError(w, http.StatusNotFound, "ROOM_NOT_FOUND", "room not found")
		case "schedule exists":
			writeError(w, http.StatusConflict, "SCHEDULE_EXISTS", "schedule for this room already exists and cannot be changed")
		case "daysOfWeek is required",
			"invalid weekday",
			"end time must be after start time",
			"time range must be multiple of 30 minutes":
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		default:
			internalError(w, err)
		}
		return
	}

	writeJSON(w, http.StatusCreated, scheduleResponse{Schedule: toScheduleDTO(schedule)})
}
