package repository

import (
	"context"
	"database/sql"

	"booking-service/internal/domain"

	"github.com/google/uuid"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error)
	GetActiveBySlotID(ctx context.Context, slotID uuid.UUID) (*domain.Booking, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.BookingWithDetails, error)
	ListAll(ctx context.Context, offset, limit int) ([]*domain.BookingWithDetails, int, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BookingStatus) error
}

type bookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	booking.ID = uuid.New()
	query := `
		INSERT INTO bookings (id, user_id, slot_id, status, conference_link, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query, booking.ID, booking.UserID, booking.SlotID, booking.Status, booking.ConferenceLink).
		Scan(&booking.CreatedAt, &booking.UpdatedAt)
}

func (r *bookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	query := `SELECT id, user_id, slot_id, status, conference_link, created_at, updated_at FROM bookings WHERE id = $1`

	var booking domain.Booking
	var conferenceLink sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.SlotID,
		&booking.Status,
		&conferenceLink,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if conferenceLink.Valid {
		booking.ConferenceLink = &conferenceLink.String
	}

	return &booking, nil
}

func (r *bookingRepository) GetActiveBySlotID(ctx context.Context, slotID uuid.UUID) (*domain.Booking, error) {
	query := `
		SELECT id, user_id, slot_id, status, conference_link, created_at, updated_at
		FROM bookings
		WHERE slot_id = $1 AND status = 'active'
	`

	var booking domain.Booking
	var conferenceLink sql.NullString

	err := r.db.QueryRowContext(ctx, query, slotID).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.SlotID,
		&booking.Status,
		&conferenceLink,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if conferenceLink.Valid {
		booking.ConferenceLink = &conferenceLink.String
	}

	return &booking, nil
}

func (r *bookingRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.BookingWithDetails, error) {
	query := `
		SELECT b.id, b.user_id, b.slot_id, b.status, b.conference_link, b.created_at, b.updated_at,
		       r.name, s.start_time, s.end_time
		FROM bookings b
		JOIN slots s ON b.slot_id = s.id
		JOIN rooms r ON s.room_id = r.id
		WHERE b.user_id = $1
		  AND s.start_time >= NOW()
		ORDER BY s.start_time
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.BookingWithDetails
	for rows.Next() {
		var item domain.BookingWithDetails
		var conferenceLink sql.NullString

		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.SlotID,
			&item.Status,
			&conferenceLink,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.RoomName,
			&item.StartTime,
			&item.EndTime,
		); err != nil {
			return nil, err
		}

		if conferenceLink.Valid {
			item.ConferenceLink = &conferenceLink.String
		}

		result = append(result, &item)
	}

	return result, rows.Err()
}

func (r *bookingRepository) ListAll(ctx context.Context, offset, limit int) ([]*domain.BookingWithDetails, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM bookings`).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT b.id, b.user_id, b.slot_id, b.status, b.conference_link, b.created_at, b.updated_at,
		       r.name, s.start_time, s.end_time
		FROM bookings b
		JOIN slots s ON b.slot_id = s.id
		JOIN rooms r ON s.room_id = r.id
		ORDER BY b.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*domain.BookingWithDetails
	for rows.Next() {
		var item domain.BookingWithDetails
		var conferenceLink sql.NullString

		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.SlotID,
			&item.Status,
			&conferenceLink,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.RoomName,
			&item.StartTime,
			&item.EndTime,
		); err != nil {
			return nil, 0, err
		}

		if conferenceLink.Valid {
			item.ConferenceLink = &conferenceLink.String
		}

		result = append(result, &item)
	}

	return result, total, rows.Err()
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BookingStatus) error {
	_, err := r.db.ExecContext(ctx, `UPDATE bookings SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	return err
}
