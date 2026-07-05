package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"booking-service/internal/service"
	"booking-service/internal/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

type CreateBookingRequest struct {
	SlotID               string `json:"slotId"`
	CreateConferenceLink bool   `json:"createConferenceLink"`
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, "user") {
		return
	}

	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}

	slotID, err := uuid.Parse(req.SlotID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid slotId")
		return
	}

	booking, err := h.bookingService.CreateBooking(r.Context(), userID, slotID, req.CreateConferenceLink)
	if err != nil {
		switch err.Error() {
		case "slot not found":
			writeError(w, http.StatusNotFound, "SLOT_NOT_FOUND", "slot not found")
		case "slot is already booked":
			writeError(w, http.StatusConflict, "SLOT_ALREADY_BOOKED", "slot is already booked")
		case "slot already booked":
			writeError(w, http.StatusConflict, "SLOT_ALREADY_BOOKED", "slot is already booked")
		case "cannot book slot in the past":
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "cannot book slot in the past")
		case "forbidden":
			writeError(w, http.StatusForbidden, "FORBIDDEN", "booking is allowed only for user role")
		default:
			internalError(w, err)
		}
		return
	}

	writeJSON(w, http.StatusCreated, bookingResponse{Booking: booking})
}

func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, "user") {
		return
	}

	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	bookingID, err := uuid.Parse(mux.Vars(r)["bookingId"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid bookingId")
		return
	}

	booking, err := h.bookingService.CancelBooking(r.Context(), bookingID, userID)
	if err != nil {
		switch err.Error() {
		case "booking not found":
			writeError(w, http.StatusNotFound, "BOOKING_NOT_FOUND", "booking not found")
		case "forbidden":
			writeError(w, http.StatusForbidden, "FORBIDDEN", "cannot cancel another user's booking")
		default:
			internalError(w, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, bookingResponse{Booking: booking})
}

func (h *BookingHandler) ListMyBookings(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, "user") {
		return
	}

	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	bookings, err := h.bookingService.GetUserBookings(r.Context(), userID)
	if err != nil {
		internalError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, bookingsResponse{Bookings: bookings})
}

func (h *BookingHandler) ListAllBookings(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, "admin") {
		return
	}

	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil || parsedPage < 1 {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid page")
			return
		}
		page = parsedPage
	}

	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize := 20
	if pageSizeStr != "" {
		parsedPageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || parsedPageSize < 1 || parsedPageSize > 100 {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid pageSize")
			return
		}
		pageSize = parsedPageSize
	}

	bookings, total, err := h.bookingService.ListAllBookings(r.Context(), page, pageSize)
	if err != nil {
		internalError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"bookings": bookings,
		"pagination": map[string]int{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
		},
	})
}
