package service

import (
	"context"
	"errors"
	"sort"
	"time"

	"booking-service/internal/domain"
	"booking-service/internal/repository"

	"github.com/google/uuid"
)

type ScheduleService struct {
	scheduleRepo repository.ScheduleRepository
	slotRepo     repository.SlotRepository
	roomRepo     repository.RoomRepository
}

func NewScheduleService(scheduleRepo repository.ScheduleRepository, slotRepo repository.SlotRepository, roomRepo repository.RoomRepository) *ScheduleService {
	return &ScheduleService{scheduleRepo: scheduleRepo, slotRepo: slotRepo, roomRepo: roomRepo}
}

func (s *ScheduleService) CreateSchedule(ctx context.Context, schedule *domain.Schedule) error {
	room, err := s.roomRepo.GetByID(ctx, schedule.RoomID)
	if err != nil {
		return err
	}
	if room == nil {
		return errors.New("room not found")
	}

	exists, err := s.scheduleRepo.ExistsForRoom(ctx, schedule.RoomID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("schedule exists")
	}

	if len(schedule.DaysOfWeek) == 0 {
		return errors.New("daysOfWeek is required")
	}
	seen := map[int]bool{}
	cleanDays := make([]int, 0, len(schedule.DaysOfWeek))
	for _, d := range schedule.DaysOfWeek {
		if d < 1 || d > 7 {
			return errors.New("invalid weekday")
		}
		if !seen[d] {
			seen[d] = true
			cleanDays = append(cleanDays, d)
		}
	}
	sort.Ints(cleanDays)
	schedule.DaysOfWeek = cleanDays

	if !schedule.EndTime.After(schedule.StartTime) {
		return errors.New("end time must be after start time")
	}
	if int(schedule.EndTime.Sub(schedule.StartTime).Minutes())%30 != 0 {
		return errors.New("time range must be multiple of 30 minutes")
	}

	for i, d := range cleanDays {
		item := &domain.Schedule{RoomID: schedule.RoomID, WeekDay: domain.WeekDay(d), StartTime: schedule.StartTime, EndTime: schedule.EndTime}
		if err := s.scheduleRepo.Create(ctx, item); err != nil {
			return err
		}
		if i == 0 {
			schedule.ID = item.ID
		}
	}
	return nil
}

func (s *ScheduleService) GenerateSlots(ctx context.Context, roomID uuid.UUID, date time.Time) error {
	weekDay := getWeekDay(date)
	schedule, err := s.scheduleRepo.GetByRoomAndWeekDay(ctx, roomID, weekDay)
	if err != nil {
		return err
	}
	if schedule == nil {
		return nil
	}

	startTime := time.Date(date.Year(), date.Month(), date.Day(), schedule.StartTime.Hour(), schedule.StartTime.Minute(), 0, 0, time.UTC)
	endTime := time.Date(date.Year(), date.Month(), date.Day(), schedule.EndTime.Hour(), schedule.EndTime.Minute(), 0, 0, time.UTC)
	var slots []*domain.Slot
	for current := startTime; current.Before(endTime); current = current.Add(30 * time.Minute) {
		slotEnd := current.Add(30 * time.Minute)
		if slotEnd.After(endTime) {
			break
		}
		slots = append(slots, &domain.Slot{RoomID: roomID, StartTime: current, EndTime: slotEnd})
	}
	if len(slots) == 0 {
		return nil
	}
	return s.slotRepo.CreateBatch(ctx, slots)
}

func getWeekDay(date time.Time) domain.WeekDay {
	switch date.Weekday() {
	case time.Monday:
		return domain.Monday
	case time.Tuesday:
		return domain.Tuesday
	case time.Wednesday:
		return domain.Wednesday
	case time.Thursday:
		return domain.Thursday
	case time.Friday:
		return domain.Friday
	case time.Saturday:
		return domain.Saturday
	default:
		return domain.Sunday
	}
}
