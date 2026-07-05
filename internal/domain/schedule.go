package domain

import (
	"time"

	"github.com/google/uuid"
)

type WeekDay int

const (
	Monday WeekDay = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

type Schedule struct {
	ID         uuid.UUID `json:"id"`
	RoomID     uuid.UUID `json:"roomId"`
	DaysOfWeek []int     `json:"daysOfWeek,omitempty"`
	WeekDay    WeekDay   `json:"-"`
	StartTime  time.Time `json:"-"`
	EndTime    time.Time `json:"-"`
	CreatedAt  time.Time `json:"createdAt"`
}
