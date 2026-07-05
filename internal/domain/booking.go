package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookingStatus string

const (
	BookingStatusActive    BookingStatus = "active"
	BookingStatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID             uuid.UUID     `json:"id"`
	UserID         uuid.UUID     `json:"userId"`
	SlotID         uuid.UUID     `json:"slotId"`
	Status         BookingStatus `json:"status"`
	ConferenceLink *string       `json:"conferenceLink,omitempty"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"-"`
}

type BookingWithDetails struct {
	Booking
	RoomName  string    `json:"-"`
	StartTime time.Time `json:"-"`
	EndTime   time.Time `json:"-"`
}
