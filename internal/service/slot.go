package service

import (
	"context"
	"time"

	"booking-service/internal/domain"
	"booking-service/internal/repository"

	"github.com/google/uuid"
)

type SlotService struct {
	slotRepo     repository.SlotRepository
	scheduleRepo repository.ScheduleRepository
	roomRepo     repository.RoomRepository
	lookupDays   int
}

func NewSlotService(slotRepo repository.SlotRepository, scheduleRepo repository.ScheduleRepository, roomRepo repository.RoomRepository, lookupDays int) *SlotService {
	return &SlotService{
		slotRepo:     slotRepo,
		scheduleRepo: scheduleRepo,
		roomRepo:     roomRepo,
		lookupDays:   lookupDays,
	}
}

func (s *SlotService) GetAvailableSlots(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*domain.Slot, error) {
	dateUTC := time.Date(date.UTC().Year(), date.UTC().Month(), date.UTC().Day(), 0, 0, 0, 0, time.UTC)
	todayUTC := time.Now().UTC().Truncate(24 * time.Hour)

	if dateUTC.Before(todayUTC) {
		return []*domain.Slot{}, nil
	}

	hasSchedule, err := s.scheduleRepo.ExistsForRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if !hasSchedule {
		return []*domain.Slot{}, nil
	}

	slots, err := s.slotRepo.GetAvailableByRoomAndDate(ctx, roomID, dateUTC)
	if err != nil {
		return nil, err
	}
	if len(slots) > 0 {
		return slots, nil
	}

	scheduleService := &ScheduleService{
		scheduleRepo: s.scheduleRepo,
		slotRepo:     s.slotRepo,
		roomRepo:     s.roomRepo,
	}
	if err := scheduleService.GenerateSlots(ctx, roomID, dateUTC); err != nil {
		return nil, err
	}

	return s.slotRepo.GetAvailableByRoomAndDate(ctx, roomID, dateUTC)
}
