package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"booking-service/internal/domain"
)

type apiErrorBody struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type roomResponse struct {
	Room interface{} `json:"room"`
}

type roomsResponse struct {
	Rooms interface{} `json:"rooms"`
}

type scheduleDTO struct {
	ID         string `json:"id,omitempty"`
	RoomID     string `json:"roomId"`
	DaysOfWeek []int  `json:"daysOfWeek"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

type scheduleResponse struct {
	Schedule scheduleDTO `json:"schedule"`
}

type slotDTO struct {
	ID     string    `json:"id"`
	RoomID string    `json:"roomId"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}

type slotsResponse struct {
	Slots []slotDTO `json:"slots"`
}

type bookingResponse struct {
	Booking *domain.Booking `json:"booking"`
}

type bookingsResponse struct {
	Bookings []*domain.Booking `json:"bookings"`
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	resp := apiErrorBody{}
	resp.Error.Code = code
	resp.Error.Message = message
	writeJSON(w, status, resp)
}

func WriteAPIError(w http.ResponseWriter, status int, code, message string) {
	writeError(w, status, code, message)
}

func internalError(w http.ResponseWriter, _ error) {
	writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
}

func toScheduleDTO(s *domain.Schedule) scheduleDTO {
	return scheduleDTO{
		ID:         s.ID.String(),
		RoomID:     s.RoomID.String(),
		DaysOfWeek: s.DaysOfWeek,
		StartTime:  s.StartTime.Format("15:04"),
		EndTime:    s.EndTime.Format("15:04"),
	}
}

func toSlotDTOs(slots []*domain.Slot) []slotDTO {
	result := make([]slotDTO, 0, len(slots))
	for _, s := range slots {
		result = append(result, slotDTO{
			ID:     s.ID.String(),
			RoomID: s.RoomID.String(),
			Start:  s.StartTime.UTC(),
			End:    s.EndTime.UTC(),
		})
	}
	return result
}
