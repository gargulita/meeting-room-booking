package service

import "errors"

var (
	ErrRoomNotFound      = errors.New("room not found")
	ErrScheduleExists    = errors.New("schedule exists")
	ErrSlotNotFound      = errors.New("slot not found")
	ErrSlotAlreadyBooked = errors.New("slot already booked")
	ErrBookingNotFound   = errors.New("booking not found")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidRequest    = errors.New("invalid request")
	ErrPastSlot          = errors.New("cannot book slot in the past")
)
