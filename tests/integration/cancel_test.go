package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCancelBooking(t *testing.T) {
	ctx := setupTest(t)

	t.Log("Creating room...")
	roomData := map[string]interface{}{
		"name":        "Cancel Test Room",
		"description": "Room for cancel integration test",
		"capacity":    8,
	}

	roomResp := makeAuthRequest(t, ctx.BaseURL+"/rooms/create", "POST", ctx.AdminToken, roomData)
	require.NotNil(t, roomResp)

	roomObj, ok := roomResp["room"].(map[string]interface{})
	require.True(t, ok, "expected room object in response")

	roomID, ok := roomObj["id"].(string)
	require.True(t, ok, "expected room.id in response")
	require.NotEmpty(t, roomID)

	t.Log("Creating schedule...")
	scheduleData := map[string]interface{}{
		"roomId":     roomID,
		"daysOfWeek": []int{1},
		"startTime":  "09:00",
		"endTime":    "17:00",
	}

	scheduleResp := makeAuthRequest(
		t,
		ctx.BaseURL+"/rooms/"+roomID+"/schedule/create",
		"POST",
		ctx.AdminToken,
		scheduleData,
	)
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
	if !ok || len(slots) == 0 {
		t.Skip("No slots available for testing, skipping cancel test")
		return
	}

	slot, ok := slots[0].(map[string]interface{})
	require.True(t, ok, "expected slot object")

	slotID, ok := slot["id"].(string)
	require.True(t, ok, "expected slot.id in response")
	require.NotEmpty(t, slotID)

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
	require.True(t, ok, "expected booking.id in response")
	require.NotEmpty(t, bookingID)

	t.Log("Cancelling booking...")
	cancelResp, statusCode := makeAuthRequestWithStatus(
		t,
		ctx.BaseURL+"/bookings/"+bookingID+"/cancel",
		"POST",
		ctx.UserToken,
		nil,
	)
	assert.Equal(t, http.StatusOK, statusCode, "cancel should return 200")
	require.NotNil(t, cancelResp)

	cancelBookingObj, ok := cancelResp["booking"].(map[string]interface{})
	require.True(t, ok, "expected booking object in cancel response")

	status, ok := cancelBookingObj["status"].(string)
	require.True(t, ok, "expected booking.status in cancel response")
	assert.Equal(t, "cancelled", status)

	t.Log("Verifying booking is NOT present in /bookings/my...")
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
	assert.False(t, found, "cancelled booking should not appear in user's future bookings")

	t.Log("Verifying slot is available again...")
	slotsAfterCancelResp := makeAuthRequest(
		t,
		ctx.BaseURL+"/rooms/"+roomID+"/slots/list?date="+dateStr,
		"GET",
		ctx.UserToken,
		nil,
	)

	slotsAfter, ok := slotsAfterCancelResp["slots"].([]interface{})
	require.True(t, ok, "expected slots field after cancellation")

	found = false
	for _, s := range slotsAfter {
		sMap, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		id, ok := sMap["id"].(string)
		if ok && id == slotID {
			found = true
			break
		}
	}
	assert.True(t, found, "slot should be available after cancellation")

	t.Log("Testing idempotency - cancelling again...")
	cancelResp2, statusCode2 := makeAuthRequestWithStatus(
		t,
		ctx.BaseURL+"/bookings/"+bookingID+"/cancel",
		"POST",
		ctx.UserToken,
		nil,
	)
	assert.Equal(t, http.StatusOK, statusCode2, "second cancel should also return 200")
	require.NotNil(t, cancelResp2)

	cancelBookingObj2, ok := cancelResp2["booking"].(map[string]interface{})
	require.True(t, ok, "expected booking object in second cancel response")

	status2, ok := cancelBookingObj2["status"].(string)
	require.True(t, ok, "expected booking.status in second cancel response")
	assert.Equal(t, "cancelled", status2)

	t.Log("Testing cancellation by admin...")
	_, statusCode3 := makeAuthRequestWithStatus(
		t,
		ctx.BaseURL+"/bookings/"+bookingID+"/cancel",
		"POST",
		ctx.AdminToken,
		nil,
	)
	assert.True(
		t,
		statusCode3 == http.StatusForbidden || statusCode3 == http.StatusUnauthorized,
		"admin should not be able to cancel user booking",
	)
}

func TestCancelNonExistentBooking(t *testing.T) {
	ctx := setupTest(t)

	nonExistentID := "99999999-9999-9999-9999-999999999999"

	t.Log("Cancelling non-existent booking...")
	_, statusCode := makeAuthRequestWithStatus(
		t,
		ctx.BaseURL+"/bookings/"+nonExistentID+"/cancel",
		"POST",
		ctx.UserToken,
		nil,
	)

	assert.Equal(t, http.StatusNotFound, statusCode, "non-existent booking should return 404")
}

func TestCancelBookingInPast(t *testing.T) {
	t.Skip("Skipping past booking test - would require database seeding with past slots")
}
