package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"booking-service/internal/domain"
	"booking-service/internal/repository"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type BookingService struct {
	bookingRepo repository.BookingRepository
	slotRepo    repository.SlotRepository
	userRepo    repository.UserRepository
	db          *sql.DB
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	slotRepo repository.SlotRepository,
	userRepo repository.UserRepository,
	db ...*sql.DB,
) *BookingService {
	svc := &BookingService{
		bookingRepo: bookingRepo,
		slotRepo:    slotRepo,
		userRepo:    userRepo,
	}
	if len(db) > 0 {
		svc.db = db[0]
	}
	return svc
}

func (s *BookingService) CreateBooking(ctx context.Context, userID, slotID uuid.UUID, createConferenceLink bool) (*domain.Booking, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	slot, err := s.slotRepo.GetByID(ctx, slotID)
	if err != nil {
		return nil, err
	}
	if slot == nil {
		return nil, errors.New("slot not found")
	}
	if slot.StartTime.Before(time.Now().UTC()) {
		return nil, errors.New("cannot book slot in the past")
	}

	existingBooking, err := s.bookingRepo.GetActiveBySlotID(ctx, slotID)
	if err != nil {
		return nil, err
	}
	if existingBooking != nil {
		return nil, errors.New("slot is already booked")
	}

	booking := &domain.Booking{
		UserID: userID,
		SlotID: slotID,
		Status: domain.BookingStatusActive,
	}
	if createConferenceLink {
		link := fmt.Sprintf("https://meet.example.com/%s", slotID.String())
		booking.ConferenceLink = &link
	}

	if s.db == nil {
		if err := s.bookingRepo.Create(ctx, booking); err != nil {
			if isUniqueViolation(err) {
				return nil, errors.New("slot is already booked")
			}
			return nil, err
		}
		if err := s.slotRepo.UpdateBookedStatus(ctx, slotID, true); err != nil {
			return nil, err
		}
		return booking, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	booking.ID = uuid.New()
	insertQuery := `
		INSERT INTO bookings (id, user_id, slot_id, status, conference_link, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	if err := tx.QueryRowContext(
		ctx,
		insertQuery,
		booking.ID,
		booking.UserID,
		booking.SlotID,
		booking.Status,
		booking.ConferenceLink,
	).Scan(&booking.CreatedAt, &booking.UpdatedAt); err != nil {
		if isUniqueViolation(err) {
			return nil, errors.New("slot is already booked")
		}
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `UPDATE slots SET is_booked = true WHERE id = $1`, slotID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		if isUniqueViolation(err) {
			return nil, errors.New("slot is already booked")
		}
		return nil, err
	}

	return booking, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, bookingID, userID uuid.UUID) (*domain.Booking, error) {
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}
	if booking.UserID != userID {
		return nil, errors.New("forbidden")
	}
	if booking.Status == domain.BookingStatusCancelled {
		return booking, nil
	}

	if s.db == nil {
		if err := s.bookingRepo.UpdateStatus(ctx, bookingID, domain.BookingStatusCancelled); err != nil {
			return nil, err
		}
		if err := s.slotRepo.UpdateBookedStatus(ctx, booking.SlotID, false); err != nil {
			return nil, err
		}
		booking.Status = domain.BookingStatusCancelled
		return booking, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(
		ctx,
		`UPDATE bookings SET status = $1, updated_at = NOW() WHERE id = $2`,
		domain.BookingStatusCancelled,
		bookingID,
	); err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `UPDATE slots SET is_booked = false WHERE id = $1`, booking.SlotID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	booking.Status = domain.BookingStatusCancelled
	return booking, nil
}

func (s *BookingService) GetUserBookings(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error) {
	rows, err := s.bookingRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Booking, 0, len(rows))
	for _, row := range rows {
		b := row.Booking
		result = append(result, &b)
	}
	return result, nil
}

func (s *BookingService) ListAllBookings(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error) {
	offset := (page - 1) * pageSize
	rows, total, err := s.bookingRepo.ListAll(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*domain.Booking, 0, len(rows))
	for _, row := range rows {
		b := row.Booking
		result = append(result, &b)
	}
	return result, total, nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}
