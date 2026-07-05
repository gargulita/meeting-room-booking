package domain

import (
	"time"

	"github.com/google/uuid"
)

type Slot struct {
	ID        uuid.UUID `json:"id"`
	RoomID    uuid.UUID `json:"roomId"`
	StartTime time.Time `json:"-"`
	EndTime   time.Time `json:"-"`
	IsBooked  bool      `json:"-"`
	CreatedAt time.Time `json:"-"`
}
