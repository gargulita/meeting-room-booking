package unit

import (
	"context"
	"testing"
	"time"

	"booking-service/internal/domain"
	"booking-service/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRoomRepository struct {
	mock.Mock
}

func (m *MockRoomRepository) Create(ctx context.Context, room *domain.Room) error {
	args := m.Called(ctx, room)
	return args.Error(0)
}

func (m *MockRoomRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Room), args.Error(1)
}

func (m *MockRoomRepository) List(ctx context.Context) ([]*domain.Room, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Room), args.Error(1)
}

func (m *MockRoomRepository) Update(ctx context.Context, room *domain.Room) error {
	args := m.Called(ctx, room)
	return args.Error(0)
}

func (m *MockRoomRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockScheduleRepository struct {
	mock.Mock
}

func (m *MockScheduleRepository) Create(ctx context.Context, schedule *domain.Schedule) error {
	args := m.Called(ctx, schedule)
	return args.Error(0)
}

func (m *MockScheduleRepository) GetByRoomID(ctx context.Context, roomID uuid.UUID) ([]*domain.Schedule, error) {
	args := m.Called(ctx, roomID)
	return args.Get(0).([]*domain.Schedule), args.Error(1)
}

func (m *MockScheduleRepository) GetByRoomAndWeekDay(ctx context.Context, roomID uuid.UUID, weekDay domain.WeekDay) (*domain.Schedule, error) {
	args := m.Called(ctx, roomID, weekDay)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Schedule), args.Error(1)
}

func (m *MockScheduleRepository) ExistsForRoom(ctx context.Context, roomID uuid.UUID) (bool, error) {
	args := m.Called(ctx, roomID)
	return args.Bool(0), args.Error(1)
}

type MockSlotRepository struct {
	mock.Mock
}

func (m *MockSlotRepository) Create(ctx context.Context, slot *domain.Slot) error {
	args := m.Called(ctx, slot)
	return args.Error(0)
}

func (m *MockSlotRepository) CreateBatch(ctx context.Context, slots []*domain.Slot) error {
	args := m.Called(ctx, slots)
	return args.Error(0)
}

func (m *MockSlotRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Slot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Slot), args.Error(1)
}

func (m *MockSlotRepository) GetAvailableByRoomAndDate(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*domain.Slot, error) {
	args := m.Called(ctx, roomID, date)
	return args.Get(0).([]*domain.Slot), args.Error(1)
}

func (m *MockSlotRepository) GetByRoomAndTimeRange(ctx context.Context, roomID uuid.UUID, start, end time.Time) ([]*domain.Slot, error) {
	args := m.Called(ctx, roomID, start, end)
	return args.Get(0).([]*domain.Slot), args.Error(1)
}

func (m *MockSlotRepository) UpdateBookedStatus(ctx context.Context, id uuid.UUID, isBooked bool) error {
	args := m.Called(ctx, id, isBooked)
	return args.Error(0)
}

func (m *MockSlotRepository) DeleteOldSlots(ctx context.Context, beforeDate time.Time) error {
	args := m.Called(ctx, beforeDate)
	return args.Error(0)
}

func (m *MockSlotRepository) GetSlotsByRoomAndDateRange(ctx context.Context, roomID uuid.UUID, startDate, endDate time.Time) ([]*domain.Slot, error) {
	args := m.Called(ctx, roomID, startDate, endDate)
	return args.Get(0).([]*domain.Slot), args.Error(1)
}

type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	args := m.Called(ctx, booking)
	return args.Error(0)
}

func (m *MockBookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) GetActiveBySlotID(ctx context.Context, slotID uuid.UUID) (*domain.Booking, error) {
	args := m.Called(ctx, slotID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.BookingWithDetails, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.BookingWithDetails), args.Error(1)
}

func (m *MockBookingRepository) ListAll(ctx context.Context, offset, limit int) ([]*domain.BookingWithDetails, int, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]*domain.BookingWithDetails), args.Int(1), args.Error(2)
}

func (m *MockBookingRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BookingStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestRoomService_CreateRoom(t *testing.T) {
	tests := []struct {
		name         string
		room         *domain.Room
		setup        func(*MockRoomRepository)
		wantErr      bool
		errMsg       string
		shouldCreate bool
	}{
		{
			name: "success",
			room: &domain.Room{
				Name:        "Test Room",
				Description: strPtr("Test Description"),
				Capacity:    intPtr(10),
			},
			setup: func(m *MockRoomRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.Room")).Return(nil)
			},
			wantErr:      false,
			shouldCreate: true,
		},
		{
			name: "missing_name",
			room: &domain.Room{
				Name:     "",
				Capacity: intPtr(10),
			},
			setup:        func(m *MockRoomRepository) {},
			wantErr:      true,
			errMsg:       "name is required",
			shouldCreate: false,
		},
		{
			name: "invalid_capacity",
			room: &domain.Room{
				Name:     "Test Room",
				Capacity: intPtr(-1),
			},
			setup:        func(m *MockRoomRepository) {},
			wantErr:      true,
			errMsg:       "capacity must be non-negative",
			shouldCreate: false,
		},
		{
			name: "nil_capacity_is_ok",
			room: &domain.Room{
				Name:        "Test Room",
				Description: strPtr("No capacity"),
				Capacity:    nil,
			},
			setup: func(m *MockRoomRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.Room")).Return(nil)
			},
			wantErr:      false,
			shouldCreate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRoomRepository)
			tt.setup(mockRepo)

			svc := service.NewRoomService(mockRepo)
			err := svc.CreateRoom(context.Background(), tt.room)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}

			if tt.shouldCreate {
				mockRepo.AssertCalled(t, "Create", mock.Anything, mock.AnythingOfType("*domain.Room"))
				mockRepo.AssertExpectations(t)
			} else {
				mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
			}
		})
	}
}

func TestBookingService_CreateBooking(t *testing.T) {
	userID := uuid.New()
	slotID := uuid.New()
	futureTime := time.Now().UTC().Add(24 * time.Hour)

	tests := []struct {
		name    string
		setup   func(*MockBookingRepository, *MockSlotRepository, *MockUserRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			setup: func(b *MockBookingRepository, s *MockSlotRepository, u *MockUserRepository) {
				u.On("GetByID", mock.Anything, userID).Return(&domain.User{ID: userID, Role: domain.RoleUser}, nil)
				s.On("GetByID", mock.Anything, slotID).Return(&domain.Slot{
					ID:        slotID,
					StartTime: futureTime,
					IsBooked:  false,
				}, nil)
				b.On("GetActiveBySlotID", mock.Anything, slotID).Return(nil, nil)
				b.On("Create", mock.Anything, mock.AnythingOfType("*domain.Booking")).Return(nil)
				s.On("UpdateBookedStatus", mock.Anything, slotID, true).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "user not found",
			setup: func(b *MockBookingRepository, s *MockSlotRepository, u *MockUserRepository) {
				u.On("GetByID", mock.Anything, userID).Return(nil, nil)
			},
			wantErr: true,
			errMsg:  "user not found",
		},
		{
			name: "slot not found",
			setup: func(b *MockBookingRepository, s *MockSlotRepository, u *MockUserRepository) {
				u.On("GetByID", mock.Anything, userID).Return(&domain.User{ID: userID}, nil)
				s.On("GetByID", mock.Anything, slotID).Return(nil, nil)
			},
			wantErr: true,
			errMsg:  "slot not found",
		},
		{
			name: "slot in the past",
			setup: func(b *MockBookingRepository, s *MockSlotRepository, u *MockUserRepository) {
				u.On("GetByID", mock.Anything, userID).Return(&domain.User{ID: userID}, nil)
				s.On("GetByID", mock.Anything, slotID).Return(&domain.Slot{
					ID:        slotID,
					StartTime: time.Now().UTC().Add(-24 * time.Hour),
				}, nil)
			},
			wantErr: true,
			errMsg:  "cannot book slot in the past",
		},
		{
			name: "slot already booked",
			setup: func(b *MockBookingRepository, s *MockSlotRepository, u *MockUserRepository) {
				u.On("GetByID", mock.Anything, userID).Return(&domain.User{ID: userID}, nil)
				s.On("GetByID", mock.Anything, slotID).Return(&domain.Slot{
					ID:        slotID,
					StartTime: futureTime,
				}, nil)
				b.On("GetActiveBySlotID", mock.Anything, slotID).Return(&domain.Booking{
					ID:     uuid.New(),
					SlotID: slotID,
					Status: domain.BookingStatusActive,
				}, nil)
			},
			wantErr: true,
			errMsg:  "slot is already booked",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBookingRepo := new(MockBookingRepository)
			mockSlotRepo := new(MockSlotRepository)
			mockUserRepo := new(MockUserRepository)

			tt.setup(mockBookingRepo, mockSlotRepo, mockUserRepo)

			svc := service.NewBookingService(mockBookingRepo, mockSlotRepo, mockUserRepo)
			_, err := svc.CreateBooking(context.Background(), userID, slotID, false)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBookingService_CancelBooking(t *testing.T) {
	userID := uuid.New()
	bookingID := uuid.New()
	slotID := uuid.New()

	tests := []struct {
		name    string
		setup   func(*MockBookingRepository, *MockSlotRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			setup: func(b *MockBookingRepository, s *MockSlotRepository) {
				b.On("GetByID", mock.Anything, bookingID).Return(&domain.Booking{
					ID:     bookingID,
					UserID: userID,
					SlotID: slotID,
					Status: domain.BookingStatusActive,
				}, nil)
				b.On("UpdateStatus", mock.Anything, bookingID, domain.BookingStatusCancelled).Return(nil)
				s.On("UpdateBookedStatus", mock.Anything, slotID, false).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "already cancelled",
			setup: func(b *MockBookingRepository, s *MockSlotRepository) {
				b.On("GetByID", mock.Anything, bookingID).Return(&domain.Booking{
					ID:     bookingID,
					UserID: userID,
					SlotID: slotID,
					Status: domain.BookingStatusCancelled,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "booking not found",
			setup: func(b *MockBookingRepository, s *MockSlotRepository) {
				b.On("GetByID", mock.Anything, bookingID).Return(nil, nil)
			},
			wantErr: true,
			errMsg:  "booking not found",
		},
		{
			name: "not the owner",
			setup: func(b *MockBookingRepository, s *MockSlotRepository) {
				b.On("GetByID", mock.Anything, bookingID).Return(&domain.Booking{
					ID:     bookingID,
					UserID: uuid.New(),
					SlotID: slotID,
					Status: domain.BookingStatusActive,
				}, nil)
			},
			wantErr: true,
			errMsg:  "forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBookingRepo := new(MockBookingRepository)
			mockSlotRepo := new(MockSlotRepository)
			mockUserRepo := new(MockUserRepository)

			tt.setup(mockBookingRepo, mockSlotRepo)

			svc := service.NewBookingService(mockBookingRepo, mockSlotRepo, mockUserRepo)
			booking, err := svc.CancelBooking(context.Background(), bookingID, userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, booking)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, booking)
			}
		})
	}
}

func TestScheduleService_CreateSchedule(t *testing.T) {
	roomID := uuid.New()
	startTime, _ := time.Parse("15:04", "09:00")
	endTime, _ := time.Parse("15:04", "17:00")

	tests := []struct {
		name     string
		schedule *domain.Schedule
		setup    func(*MockScheduleRepository, *MockRoomRepository)
		wantErr  bool
		errMsg   string
	}{
		{
			name: "success",
			schedule: &domain.Schedule{
				RoomID:     roomID,
				DaysOfWeek: []int{1},
				StartTime:  startTime,
				EndTime:    endTime,
			},
			setup: func(s *MockScheduleRepository, r *MockRoomRepository) {
				r.On("GetByID", mock.Anything, roomID).Return(&domain.Room{ID: roomID}, nil)
				s.On("ExistsForRoom", mock.Anything, roomID).Return(false, nil)
				s.On("Create", mock.Anything, mock.AnythingOfType("*domain.Schedule")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "room not found",
			schedule: &domain.Schedule{
				RoomID:     roomID,
				DaysOfWeek: []int{1},
				StartTime:  startTime,
				EndTime:    endTime,
			},
			setup: func(s *MockScheduleRepository, r *MockRoomRepository) {
				r.On("GetByID", mock.Anything, roomID).Return(nil, nil)
			},
			wantErr: true,
			errMsg:  "room not found",
		},
		{
			name: "schedule exists",
			schedule: &domain.Schedule{
				RoomID:     roomID,
				DaysOfWeek: []int{1},
				StartTime:  startTime,
				EndTime:    endTime,
			},
			setup: func(s *MockScheduleRepository, r *MockRoomRepository) {
				r.On("GetByID", mock.Anything, roomID).Return(&domain.Room{ID: roomID}, nil)
				s.On("ExistsForRoom", mock.Anything, roomID).Return(true, nil)
			},
			wantErr: true,
			errMsg:  "schedule exists",
		},
		{
			name: "invalid time range",
			schedule: &domain.Schedule{
				RoomID:     roomID,
				DaysOfWeek: []int{1},
				StartTime:  endTime,
				EndTime:    startTime,
			},
			setup: func(s *MockScheduleRepository, r *MockRoomRepository) {
				r.On("GetByID", mock.Anything, roomID).Return(&domain.Room{ID: roomID}, nil)
				s.On("ExistsForRoom", mock.Anything, roomID).Return(false, nil)
			},
			wantErr: true,
			errMsg:  "end time must be after start time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockScheduleRepo := new(MockScheduleRepository)
			mockSlotRepo := new(MockSlotRepository)
			mockRoomRepo := new(MockRoomRepository)

			tt.setup(mockScheduleRepo, mockRoomRepo)

			svc := service.NewScheduleService(mockScheduleRepo, mockSlotRepo, mockRoomRepo)
			err := svc.CreateSchedule(context.Background(), tt.schedule)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestScheduleService_CreateSchedule_ValidationCases(t *testing.T) {
	roomID := uuid.New()
	baseStart, _ := time.Parse("15:04", "09:00")
	baseEnd, _ := time.Parse("15:04", "17:00")
	shortEnd, _ := time.Parse("15:04", "09:45")

	tests := []struct {
		name     string
		schedule *domain.Schedule
		setup    func(*MockScheduleRepository, *MockRoomRepository)
		wantErr  string
	}{
		{
			name: "days_of_week_required",
			schedule: &domain.Schedule{
				RoomID:     roomID,
				DaysOfWeek: []int{},
				StartTime:  baseStart,
				EndTime:    baseEnd,
			},
			setup: func(s *MockScheduleRepository, r *MockRoomRepository) {
				r.On("GetByID", mock.Anything, roomID).Return(&domain.Room{ID: roomID}, nil)
				s.On("ExistsForRoom", mock.Anything, roomID).Return(false, nil)
			},
			wantErr: "daysOfWeek is required",
		},
		{
			name: "invalid_weekday_zero",
			schedule: &domain.Schedule{
				RoomID:     roomID,
				DaysOfWeek: []int{0},
				StartTime:  baseStart,
				EndTime:    baseEnd,
			},
			setup: func(s *MockScheduleRepository, r *MockRoomRepository) {
				r.On("GetByID", mock.Anything, roomID).Return(&domain.Room{ID: roomID}, nil)
				s.On("ExistsForRoom", mock.Anything, roomID).Return(false, nil)
			},
			wantErr: "invalid weekday",
		},
		{
			name: "invalid_weekday_eight",
			schedule: &domain.Schedule{
				RoomID:     roomID,
				DaysOfWeek: []int{8},
				StartTime:  baseStart,
				EndTime:    baseEnd,
			},
			setup: func(s *MockScheduleRepository, r *MockRoomRepository) {
				r.On("GetByID", mock.Anything, roomID).Return(&domain.Room{ID: roomID}, nil)
				s.On("ExistsForRoom", mock.Anything, roomID).Return(false, nil)
			},
			wantErr: "invalid weekday",
		},
		{
			name: "time_range_must_be_multiple_of_30_minutes",
			schedule: &domain.Schedule{
				RoomID:     roomID,
				DaysOfWeek: []int{1},
				StartTime:  baseStart,
				EndTime:    shortEnd,
			},
			setup: func(s *MockScheduleRepository, r *MockRoomRepository) {
				r.On("GetByID", mock.Anything, roomID).Return(&domain.Room{ID: roomID}, nil)
				s.On("ExistsForRoom", mock.Anything, roomID).Return(false, nil)
			},
			wantErr: "time range must be multiple of 30 minutes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockScheduleRepo := new(MockScheduleRepository)
			mockSlotRepo := new(MockSlotRepository)
			mockRoomRepo := new(MockRoomRepository)

			tt.setup(mockScheduleRepo, mockRoomRepo)

			svc := service.NewScheduleService(mockScheduleRepo, mockSlotRepo, mockRoomRepo)
			err := svc.CreateSchedule(context.Background(), tt.schedule)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestSlotService_GetAvailableSlots_ReturnsEmptyForPastDate(t *testing.T) {
	mockSlotRepo := new(MockSlotRepository)
	mockScheduleRepo := new(MockScheduleRepository)
	mockRoomRepo := new(MockRoomRepository)

	svc := service.NewSlotService(mockSlotRepo, mockScheduleRepo, mockRoomRepo, 7)

	roomID := uuid.New()
	pastDate := time.Now().UTC().AddDate(0, 0, -1)

	slots, err := svc.GetAvailableSlots(context.Background(), roomID, pastDate)
	assert.NoError(t, err)
	assert.Empty(t, slots)

	mockSlotRepo.AssertNotCalled(t, "GetAvailableByRoomAndDate", mock.Anything, mock.Anything, mock.Anything)
	mockSlotRepo.AssertNotCalled(t, "CreateBatch", mock.Anything, mock.Anything)
}

func TestSlotService_GetAvailableSlots_ReturnsEmptyWhenNoScheduleExists(t *testing.T) {
	mockSlotRepo := new(MockSlotRepository)
	mockScheduleRepo := new(MockScheduleRepository)
	mockRoomRepo := new(MockRoomRepository)

	roomID := uuid.New()
	date := time.Date(2026, 3, 30, 0, 0, 0, 0, time.UTC)

	mockScheduleRepo.
		On("ExistsForRoom", mock.Anything, roomID).
		Return(false, nil)

	svc := service.NewSlotService(mockSlotRepo, mockScheduleRepo, mockRoomRepo, 7)

	slots, err := svc.GetAvailableSlots(context.Background(), roomID, date)
	assert.NoError(t, err)
	assert.Empty(t, slots)

	mockScheduleRepo.AssertExpectations(t)
	mockSlotRepo.AssertNotCalled(t, "GetAvailableByRoomAndDate", mock.Anything, mock.Anything, mock.Anything)
	mockSlotRepo.AssertNotCalled(t, "CreateBatch", mock.Anything, mock.Anything)
}

func TestSlotService_GetAvailableSlots_ReturnsExistingSlots(t *testing.T) {
	mockSlotRepo := new(MockSlotRepository)
	mockScheduleRepo := new(MockScheduleRepository)
	mockRoomRepo := new(MockRoomRepository)

	roomID := uuid.New()
	date := time.Date(2026, 3, 30, 0, 0, 0, 0, time.UTC)

	expectedSlots := []*domain.Slot{
		{
			ID:        uuid.New(),
			RoomID:    roomID,
			StartTime: time.Date(2026, 3, 30, 9, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2026, 3, 30, 9, 30, 0, 0, time.UTC),
		},
	}

	mockScheduleRepo.
		On("ExistsForRoom", mock.Anything, roomID).
		Return(true, nil)

	mockSlotRepo.
		On("GetAvailableByRoomAndDate", mock.Anything, roomID, date).
		Return(expectedSlots, nil)

	svc := service.NewSlotService(mockSlotRepo, mockScheduleRepo, mockRoomRepo, 7)

	slots, err := svc.GetAvailableSlots(context.Background(), roomID, date)
	assert.NoError(t, err)
	assert.Equal(t, expectedSlots, slots)

	mockScheduleRepo.AssertExpectations(t)
	mockSlotRepo.AssertExpectations(t)
	mockSlotRepo.AssertNotCalled(t, "CreateBatch", mock.Anything, mock.Anything)
}

func TestScheduleService_CreateSchedule_EndTimeMustBeAfterStartTime(t *testing.T) {
	mockScheduleRepo := new(MockScheduleRepository)
	mockSlotRepo := new(MockSlotRepository)
	mockRoomRepo := new(MockRoomRepository)

	roomID := uuid.New()
	startTime, _ := time.Parse("15:04", "17:00")
	endTime, _ := time.Parse("15:04", "09:00")

	schedule := &domain.Schedule{
		RoomID:     roomID,
		DaysOfWeek: []int{1},
		StartTime:  startTime,
		EndTime:    endTime,
	}

	mockRoomRepo.
		On("GetByID", mock.Anything, roomID).
		Return(&domain.Room{ID: roomID, Name: "Room A"}, nil)

	mockScheduleRepo.
		On("ExistsForRoom", mock.Anything, roomID).
		Return(false, nil)

	svc := service.NewScheduleService(mockScheduleRepo, mockSlotRepo, mockRoomRepo)

	err := svc.CreateSchedule(context.Background(), schedule)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "end time must be after start time")

	mockRoomRepo.AssertExpectations(t)
	mockScheduleRepo.AssertExpectations(t)
}
