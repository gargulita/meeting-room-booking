package unit

import (
	"bytes"
	"net/http"
	"testing"

	"booking-service/internal/handlers"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestRoomsCreate_RequiresAuthentication(t *testing.T) {
	roomHandler := handlers.NewRoomHandler(nil)

	rec := performAuthedRequest(
		t,
		"",
		http.MethodPost,
		"/rooms/create",
		bytes.NewBufferString(`{"name":"A"}`),
		func(r *mux.Router) {
			r.HandleFunc("/rooms/create", roomHandler.CreateRoom).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestSchedulesCreate_RequiresAuthentication(t *testing.T) {
	scheduleHandler := handlers.NewScheduleHandler(nil)
	roomID := uuid.New()

	rec := performAuthedRequest(
		t,
		"",
		http.MethodPost,
		"/rooms/"+roomID.String()+"/schedule/create",
		bytes.NewBufferString(`{}`),
		func(r *mux.Router) {
			r.HandleFunc("/rooms/{roomId}/schedule/create", scheduleHandler.CreateSchedule).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestSchedulesCreate_RejectsMissingRoomIDInBody(t *testing.T) {
	token := getTestToken(t, "admin")
	scheduleHandler := handlers.NewScheduleHandler(nil)
	roomID := uuid.New()

	rec := performAuthedRequest(
		t,
		token,
		http.MethodPost,
		"/rooms/"+roomID.String()+"/schedule/create",
		bytes.NewBufferString(`{
			"daysOfWeek":[1,2],
			"startTime":"09:00",
			"endTime":"10:00"
		}`),
		func(r *mux.Router) {
			r.HandleFunc("/rooms/{roomId}/schedule/create", scheduleHandler.CreateSchedule).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestSchedulesCreate_RejectsInvalidStartTime(t *testing.T) {
	token := getTestToken(t, "admin")
	scheduleHandler := handlers.NewScheduleHandler(nil)
	roomID := uuid.New()

	rec := performAuthedRequest(
		t,
		token,
		http.MethodPost,
		"/rooms/"+roomID.String()+"/schedule/create",
		bytes.NewBufferString(`{
			"roomId":"`+roomID.String()+`",
			"daysOfWeek":[1,2],
			"startTime":"99:99",
			"endTime":"10:00"
		}`),
		func(r *mux.Router) {
			r.HandleFunc("/rooms/{roomId}/schedule/create", scheduleHandler.CreateSchedule).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestSchedulesCreate_RejectsInvalidEndTime(t *testing.T) {
	token := getTestToken(t, "admin")
	scheduleHandler := handlers.NewScheduleHandler(nil)
	roomID := uuid.New()

	rec := performAuthedRequest(
		t,
		token,
		http.MethodPost,
		"/rooms/"+roomID.String()+"/schedule/create",
		bytes.NewBufferString(`{
			"roomId":"`+roomID.String()+`",
			"daysOfWeek":[1,2],
			"startTime":"09:00",
			"endTime":"25:61"
		}`),
		func(r *mux.Router) {
			r.HandleFunc("/rooms/{roomId}/schedule/create", scheduleHandler.CreateSchedule).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestSlotsList_RequiresAuthentication(t *testing.T) {
	slotHandler := handlers.NewSlotHandler(nil, nil)
	roomID := uuid.New()

	rec := performAuthedRequest(
		t,
		"",
		http.MethodGet,
		"/rooms/"+roomID.String()+"/slots/list?date=2026-03-30",
		nil,
		func(r *mux.Router) {
			r.HandleFunc("/rooms/{roomId}/slots/list", slotHandler.ListAvailableSlots).Methods(http.MethodGet)
		},
	)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestBookingsCreate_RequiresAuthentication(t *testing.T) {
	bookingHandler := handlers.NewBookingHandler(nil)

	rec := performAuthedRequest(
		t,
		"",
		http.MethodPost,
		"/bookings/create",
		bytes.NewBufferString(`{"slotId":"`+uuid.New().String()+`"}`),
		func(r *mux.Router) {
			r.HandleFunc("/bookings/create", bookingHandler.CreateBooking).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestBookingsCancel_RequiresAuthentication(t *testing.T) {
	bookingHandler := handlers.NewBookingHandler(nil)

	rec := performAuthedRequest(
		t,
		"",
		http.MethodPost,
		"/bookings/"+uuid.New().String()+"/cancel",
		nil,
		func(r *mux.Router) {
			r.HandleFunc("/bookings/{bookingId}/cancel", bookingHandler.CancelBooking).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestBookingsMy_RequiresAuthentication(t *testing.T) {
	bookingHandler := handlers.NewBookingHandler(nil)

	rec := performAuthedRequest(
		t,
		"",
		http.MethodGet,
		"/bookings/my",
		nil,
		func(r *mux.Router) {
			r.HandleFunc("/bookings/my", bookingHandler.ListMyBookings).Methods(http.MethodGet)
		},
	)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestBookingsList_RequiresAuthentication(t *testing.T) {
	bookingHandler := handlers.NewBookingHandler(nil)

	rec := performAuthedRequest(
		t,
		"",
		http.MethodGet,
		"/bookings/list",
		nil,
		func(r *mux.Router) {
			r.HandleFunc("/bookings/list", bookingHandler.ListAllBookings).Methods(http.MethodGet)
		},
	)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestBookingsList_RejectsPageLessThanOne(t *testing.T) {
	token := getTestToken(t, "admin")
	bookingHandler := handlers.NewBookingHandler(nil)

	rec := performAuthedRequest(
		t,
		token,
		http.MethodGet,
		"/bookings/list?page=0",
		nil,
		func(r *mux.Router) {
			r.HandleFunc("/bookings/list", bookingHandler.ListAllBookings).Methods(http.MethodGet)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestBookingsList_RejectsPageSizeLessThanOne(t *testing.T) {
	token := getTestToken(t, "admin")
	bookingHandler := handlers.NewBookingHandler(nil)

	rec := performAuthedRequest(
		t,
		token,
		http.MethodGet,
		"/bookings/list?page=1&pageSize=0",
		nil,
		func(r *mux.Router) {
			r.HandleFunc("/bookings/list", bookingHandler.ListAllBookings).Methods(http.MethodGet)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}
