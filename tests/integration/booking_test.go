package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestContext struct {
	BaseURL    string
	AdminToken string
	UserToken  string
	RoomID     string
	SlotID     string
}

func setupTest(t *testing.T) *TestContext {
	ctx := &TestContext{
		BaseURL: "http://localhost:8080",
	}

	adminResp := makeRequest(t, ctx.BaseURL+"/dummyLogin", "POST", map[string]string{
		"role": "admin",
	})
	ctx.AdminToken = adminResp["token"].(string)

	userResp := makeRequest(t, ctx.BaseURL+"/dummyLogin", "POST", map[string]string{
		"role": "user",
	})
	ctx.UserToken = userResp["token"].(string)

	return ctx
}

func makeRequest(t *testing.T, url, method string, body interface{}) map[string]interface{} {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
	}

	return result
}

func makeAuthRequest(t *testing.T, url, method, token string, body interface{}) map[string]interface{} {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
	}

	return result
}

func makeAuthRequestWithStatus(t *testing.T, url, method, token string, body interface{}) (map[string]interface{}, int) {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
	}

	return result, resp.StatusCode
}

func TestCreateRoomScheduleAndBooking(t *testing.T) {
	ctx := setupTest(t)

	t.Log("Creating room...")
	roomData := map[string]interface{}{
		"name":        "Integration Test Room",
		"description": "Room for integration testing",
		"capacity":    10,
	}

	roomResp := makeAuthRequest(t, ctx.BaseURL+"/rooms/create", "POST", ctx.AdminToken, roomData)
	require.NotNil(t, roomResp)

	roomObj, ok := roomResp["room"].(map[string]interface{})
	require.True(t, ok, "expected room object in response")

	roomID, ok := roomObj["id"].(string)
	require.True(t, ok, "expected room.id in response")
	require.NotEmpty(t, roomID)
	ctx.RoomID = roomID

	t.Log("Creating schedule...")
	scheduleData := map[string]interface{}{
		"roomId":     roomID,
		"daysOfWeek": []int{1},
		"startTime":  "09:00",
		"endTime":    "17:00",
	}

	scheduleResp := makeAuthRequest(t, ctx.BaseURL+"/rooms/"+roomID+"/schedule/create", "POST", ctx.AdminToken, scheduleData)
	require.NotNil(t, scheduleResp)

	_, ok = scheduleResp["schedule"].(map[string]interface{})
	require.True(t, ok, "expected schedule object in response")

	now := time.Now().UTC()
	daysUntilMonday := (8 - int(now.Weekday())) % 7
	if daysUntilMonday == 0 {
		daysUntilMonday = 7
	}
	futureMonday := now.AddDate(0, 0, daysUntilMonday)
	dateStr := futureMonday.Format("2006-01-02")

	t.Logf("Getting available slots for %s...", dateStr)
	slotsResp := makeAuthRequest(
		t,
		ctx.BaseURL+"/rooms/"+roomID+"/slots/list?date="+dateStr,
		"GET",
		ctx.UserToken,
		nil,
	)

	slots, ok := slotsResp["slots"].([]interface{})
	require.True(t, ok, "expected slots field in response")
	if len(slots) == 0 {
		t.Skip("No slots available for testing, skipping booking test")
		return
	}

	slot, ok := slots[0].(map[string]interface{})
	require.True(t, ok, "expected slot object")

	slotID, ok := slot["id"].(string)
	require.True(t, ok, "expected slot.id")
	require.NotEmpty(t, slotID)
	ctx.SlotID = slotID

	t.Log("Creating booking...")
	bookingData := map[string]interface{}{
		"slotId":               slotID,
		"createConferenceLink": false,
	}

	bookingResp := makeAuthRequest(t, ctx.BaseURL+"/bookings/create", "POST", ctx.UserToken, bookingData)
	require.NotNil(t, bookingResp)

	bookingObj, ok := bookingResp["booking"].(map[string]interface{})
	require.True(t, ok, "expected booking object in response")

	bookingID, ok := bookingObj["id"].(string)
	require.True(t, ok, "expected booking.id")
	require.NotEmpty(t, bookingID)

	t.Log("Verifying booking in user's list...")
	myBookingsResp := makeAuthRequest(t, ctx.BaseURL+"/bookings/my", "GET", ctx.UserToken, nil)
	require.NotNil(t, myBookingsResp)

	bookings, ok := myBookingsResp["bookings"].([]interface{})
	require.True(t, ok, "expected bookings field in response")

	found := false
	for _, b := range bookings {
		bMap, ok := b.(map[string]interface{})
		if !ok {
			continue
		}
		id, ok := bMap["id"].(string)
		if ok && id == bookingID {
			found = true
			break
		}
	}
	assert.True(t, found, "booking not found in user's bookings list")

	t.Log("Verifying slot is booked...")
	_, statusCode := makeAuthRequestWithStatus(t, ctx.BaseURL+"/bookings/create", "POST", ctx.UserToken, bookingData)
	assert.Equal(t, http.StatusConflict, statusCode, "should not be able to book same slot twice")
}
