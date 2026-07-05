package unit

import (
	"bytes"
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

func getTestToken(t *testing.T, role string) string {
	t.Helper()

	token, err := utils.GenerateJWT(uuid.New(), role, time.Hour, testJWTSecret)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	return token
}

func performAuthedRequest(
	t *testing.T,
	token string,
	method string,
	path string,
	body *bytes.Buffer,
	register func(r *mux.Router),
) *httptest.ResponseRecorder {
	t.Helper()

	router := mux.NewRouter()
	api := router.PathPrefix("/").Subrouter()
	api.Use(middleware.JWTAuth(testJWTSecret))
	register(api)

	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, body)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}

func TestCreateRoom_ForbiddenForUser(t *testing.T) {
	token := getTestToken(t, "user")
	roomHandler := handlers.NewRoomHandler(nil)

	rec := performAuthedRequest(
		t,
		token,
		http.MethodPost,
		"/rooms/create",
		bytes.NewBufferString(`{"name":"A"}`),
		func(r *mux.Router) {
			r.HandleFunc("/rooms/create", roomHandler.CreateRoom).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestCreateRoom_InvalidBody(t *testing.T) {
	token := getTestToken(t, "admin")
	roomHandler := handlers.NewRoomHandler(nil)

	rec := performAuthedRequest(
		t,
		token,
		http.MethodPost,
		"/rooms/create",
		bytes.NewBufferString(`{"name":`),
		func(r *mux.Router) {
			r.HandleFunc("/rooms/create", roomHandler.CreateRoom).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestCreateSchedule_ForbiddenForUser(t *testing.T) {
	token := getTestToken(t, "user")
	scheduleHandler := handlers.NewScheduleHandler(nil)
	roomID := uuid.New()

	rec := performAuthedRequest(
		t,
		token,
		http.MethodPost,
		"/rooms/"+roomID.String()+"/schedule/create",
		bytes.NewBufferString(`{}`),
		func(r *mux.Router) {
			r.HandleFunc("/rooms/{roomId}/schedule/create", scheduleHandler.CreateSchedule).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestCreateSchedule_InvalidRoomID(t *testing.T) {
	token := getTestToken(t, "admin")
	scheduleHandler := handlers.NewScheduleHandler(nil)

	rec := performAuthedRequest(
		t,
		token,
		http.MethodPost,
		"/rooms/not-a-uuid/schedule/create",
		bytes.NewBufferString(`{}`),
		func(r *mux.Router) {
			r.HandleFunc("/rooms/{roomId}/schedule/create", scheduleHandler.CreateSchedule).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestCreateSchedule_InvalidBody(t *testing.T) {
	token := getTestToken(t, "admin")
	scheduleHandler := handlers.NewScheduleHandler(nil)
	roomID := uuid.New()

	rec := performAuthedRequest(
		t,
		token,
		http.MethodPost,
		"/rooms/"+roomID.String()+"/schedule/create",
		bytes.NewBufferString(`{"roomId":`),
		func(r *mux.Router) {
			r.HandleFunc("/rooms/{roomId}/schedule/create", scheduleHandler.CreateSchedule).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestCreateSchedule_RoomIDMismatch(t *testing.T) {
	token := getTestToken(t, "admin")
	scheduleHandler := handlers.NewScheduleHandler(nil)

	pathRoomID := uuid.New()
	bodyRoomID := uuid.New()

	body := `{
		"roomId":"` + bodyRoomID.String() + `",
		"daysOfWeek":[1,2],
		"startTime":"09:00",
		"endTime":"10:00"
	}`

	rec := performAuthedRequest(
		t,
		token,
		http.MethodPost,
		"/rooms/"+pathRoomID.String()+"/schedule/create",
		bytes.NewBufferString(body),
		func(r *mux.Router) {
			r.HandleFunc("/rooms/{roomId}/schedule/create", scheduleHandler.CreateSchedule).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestListAvailableSlots_MissingDate(t *testing.T) {
	token := getTestToken(t, "user")
	slotHandler := handlers.NewSlotHandler(nil, nil)
	roomID := uuid.New()

	rec := performAuthedRequest(
		t,
		token,
		http.MethodGet,
		"/rooms/"+roomID.String()+"/slots/list",
		nil,
		func(r *mux.Router) {
			r.HandleFunc("/rooms/{roomId}/slots/list", slotHandler.ListAvailableSlots).Methods(http.MethodGet)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestListAvailableSlots_InvalidDate(t *testing.T) {
	token := getTestToken(t, "user")
	slotHandler := handlers.NewSlotHandler(nil, nil)
	roomID := uuid.New()

	rec := performAuthedRequest(
		t,
		token,
		http.MethodGet,
		"/rooms/"+roomID.String()+"/slots/list?date=not-a-date",
		nil,
		func(r *mux.Router) {
			r.HandleFunc("/rooms/{roomId}/slots/list", slotHandler.ListAvailableSlots).Methods(http.MethodGet)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestListAvailableSlots_InvalidRoomID(t *testing.T) {
	token := getTestToken(t, "user")
	slotHandler := handlers.NewSlotHandler(nil, nil)

	rec := performAuthedRequest(
		t,
		token,
		http.MethodGet,
		"/rooms/not-a-uuid/slots/list?date=2026-03-30",
		nil,
		func(r *mux.Router) {
			r.HandleFunc("/rooms/{roomId}/slots/list", slotHandler.ListAvailableSlots).Methods(http.MethodGet)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestCreateBooking_ForbiddenForAdmin(t *testing.T) {
	token := getTestToken(t, "admin")
	bookingHandler := handlers.NewBookingHandler(nil)

	rec := performAuthedRequest(
		t,
		token,
		http.MethodPost,
		"/bookings/create",
		bytes.NewBufferString(`{"slotId":"`+uuid.New().String()+`"}`),
		func(r *mux.Router) {
			r.HandleFunc("/bookings/create", bookingHandler.CreateBooking).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestCreateBooking_InvalidBody(t *testing.T) {
	token := getTestToken(t, "user")
	bookingHandler := handlers.NewBookingHandler(nil)

	rec := performAuthedRequest(
		t,
		token,
		http.MethodPost,
		"/bookings/create",
		bytes.NewBufferString(`{"slotId":`),
		func(r *mux.Router) {
			r.HandleFunc("/bookings/create", bookingHandler.CreateBooking).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestCreateBooking_InvalidSlotID(t *testing.T) {
	token := getTestToken(t, "user")
	bookingHandler := handlers.NewBookingHandler(nil)

	rec := performAuthedRequest(
		t,
		token,
		http.MethodPost,
		"/bookings/create",
		bytes.NewBufferString(`{"slotId":"not-a-uuid"}`),
		func(r *mux.Router) {
			r.HandleFunc("/bookings/create", bookingHandler.CreateBooking).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestCancelBooking_ForbiddenForAdmin(t *testing.T) {
	token := getTestToken(t, "admin")
	bookingHandler := handlers.NewBookingHandler(nil)

	rec := performAuthedRequest(
		t,
		token,
		http.MethodPost,
		"/bookings/"+uuid.New().String()+"/cancel",
		nil,
		func(r *mux.Router) {
			r.HandleFunc("/bookings/{bookingId}/cancel", bookingHandler.CancelBooking).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestCancelBooking_InvalidBookingID(t *testing.T) {
	token := getTestToken(t, "user")
	bookingHandler := handlers.NewBookingHandler(nil)

	rec := performAuthedRequest(
		t,
		token,
		http.MethodPost,
		"/bookings/not-a-uuid/cancel",
		nil,
		func(r *mux.Router) {
			r.HandleFunc("/bookings/{bookingId}/cancel", bookingHandler.CancelBooking).Methods(http.MethodPost)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestListMyBookings_ForbiddenForAdmin(t *testing.T) {
	token := getTestToken(t, "admin")
	bookingHandler := handlers.NewBookingHandler(nil)

	rec := performAuthedRequest(
		t,
		token,
		http.MethodGet,
		"/bookings/my",
		nil,
		func(r *mux.Router) {
			r.HandleFunc("/bookings/my", bookingHandler.ListMyBookings).Methods(http.MethodGet)
		},
	)

	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestListAllBookings_InvalidPage(t *testing.T) {
	token := getTestToken(t, "admin")
	bookingHandler := handlers.NewBookingHandler(nil)

	rec := performAuthedRequest(
		t,
		token,
		http.MethodGet,
		"/bookings/list?page=abc",
		nil,
		func(r *mux.Router) {
			r.HandleFunc("/bookings/list", bookingHandler.ListAllBookings).Methods(http.MethodGet)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}

func TestListAllBookings_InvalidPageSize(t *testing.T) {
	token := getTestToken(t, "admin")
	bookingHandler := handlers.NewBookingHandler(nil)

	rec := performAuthedRequest(
		t,
		token,
		http.MethodGet,
		"/bookings/list?page=1&pageSize=101",
		nil,
		func(r *mux.Router) {
			r.HandleFunc("/bookings/list", bookingHandler.ListAllBookings).Methods(http.MethodGet)
		},
	)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REQUEST")
}
