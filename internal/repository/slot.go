package repository

import (
	"context"
	"database/sql"
	"time"

	"booking-service/internal/domain"

	"github.com/google/uuid"
)

type SlotRepository interface {
	Create(ctx context.Context, slot *domain.Slot) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Slot, error)
	GetAvailableByRoomAndDate(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*domain.Slot, error)
	GetByRoomAndTimeRange(ctx context.Context, roomID uuid.UUID, start, end time.Time) ([]*domain.Slot, error)
	UpdateBookedStatus(ctx context.Context, id uuid.UUID, isBooked bool) error
	DeleteOldSlots(ctx context.Context, beforeDate time.Time) error
	CreateBatch(ctx context.Context, slots []*domain.Slot) error
	GetSlotsByRoomAndDateRange(ctx context.Context, roomID uuid.UUID, startDate, endDate time.Time) ([]*domain.Slot, error)
}

type slotRepository struct {
	db *sql.DB
}

func NewSlotRepository(db *sql.DB) SlotRepository {
	return &slotRepository{db: db}
}

func (r *slotRepository) Create(ctx context.Context, slot *domain.Slot) error {
	slot.ID = uuid.New()
	query := `
		INSERT INTO slots (id, room_id, start_time, end_time, is_booked, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`
	_, err := r.db.ExecContext(ctx, query, slot.ID, slot.RoomID, slot.StartTime, slot.EndTime, slot.IsBooked)
	return err
}

func (r *slotRepository) CreateBatch(ctx context.Context, slots []*domain.Slot) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO slots (id, room_id, start_time, end_time, is_booked, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (room_id, start_time) DO NOTHING
	`

	for _, slot := range slots {
		slot.ID = uuid.New()
		if _, err := tx.ExecContext(ctx, query, slot.ID, slot.RoomID, slot.StartTime, slot.EndTime, slot.IsBooked); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *slotRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Slot, error) {
	query := `
		SELECT id, room_id, start_time, end_time, is_booked, created_at
		FROM slots
		WHERE id = $1
	`
	var slot domain.Slot
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&slot.ID, &slot.RoomID, &slot.StartTime, &slot.EndTime, &slot.IsBooked, &slot.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *slotRepository) GetAvailableByRoomAndDate(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*domain.Slot, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)
	effectiveStart := startOfDay
	now := time.Now().UTC()
	if sameUTCDate(date, now) && now.After(effectiveStart) {
		effectiveStart = now
	}

	query := `
		SELECT id, room_id, start_time, end_time, is_booked, created_at
		FROM slots
		WHERE room_id = $1 AND start_time >= $2 AND start_time < $3 AND is_booked = false
		ORDER BY start_time
	`
	rows, err := r.db.QueryContext(ctx, query, roomID, effectiveStart, endOfDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*domain.Slot
	for rows.Next() {
		var slot domain.Slot
		if err := rows.Scan(&slot.ID, &slot.RoomID, &slot.StartTime, &slot.EndTime, &slot.IsBooked, &slot.CreatedAt); err != nil {
			return nil, err
		}
		slots = append(slots, &slot)
	}
	return slots, rows.Err()
}

func (r *slotRepository) GetByRoomAndTimeRange(ctx context.Context, roomID uuid.UUID, start, end time.Time) ([]*domain.Slot, error) {
	query := `
		SELECT id, room_id, start_time, end_time, is_booked, created_at
		FROM slots
		WHERE room_id = $1 AND start_time >= $2 AND end_time <= $3
		ORDER BY start_time
	`
	rows, err := r.db.QueryContext(ctx, query, roomID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*domain.Slot
	for rows.Next() {
		var slot domain.Slot
		if err := rows.Scan(&slot.ID, &slot.RoomID, &slot.StartTime, &slot.EndTime, &slot.IsBooked, &slot.CreatedAt); err != nil {
			return nil, err
		}
		slots = append(slots, &slot)
	}
	return slots, rows.Err()
}

func (r *slotRepository) UpdateBookedStatus(ctx context.Context, id uuid.UUID, isBooked bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE slots SET is_booked = $1 WHERE id = $2`, isBooked, id)
	return err
}

func (r *slotRepository) DeleteOldSlots(ctx context.Context, beforeDate time.Time) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM slots WHERE start_time < $1`, beforeDate)
	return err
}

func (r *slotRepository) GetSlotsByRoomAndDateRange(ctx context.Context, roomID uuid.UUID, startDate, endDate time.Time) ([]*domain.Slot, error) {
	query := `
		SELECT id, room_id, start_time, end_time, is_booked, created_at
		FROM slots
		WHERE room_id = $1 AND start_time >= $2 AND start_time < $3
		ORDER BY start_time
	`
	rows, err := r.db.QueryContext(ctx, query, roomID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*domain.Slot
	for rows.Next() {
		var slot domain.Slot
		if err := rows.Scan(&slot.ID, &slot.RoomID, &slot.StartTime, &slot.EndTime, &slot.IsBooked, &slot.CreatedAt); err != nil {
			return nil, err
		}
		slots = append(slots, &slot)
	}
	return slots, rows.Err()
}

func sameUTCDate(a, b time.Time) bool {
	au := a.UTC()
	bu := b.UTC()
	return au.Year() == bu.Year() && au.Month() == bu.Month() && au.Day() == bu.Day()
}
