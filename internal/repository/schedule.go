package repository

import (
	"context"
	"database/sql"

	"booking-service/internal/domain"

	"github.com/google/uuid"
)

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *domain.Schedule) error
	GetByRoomID(ctx context.Context, roomID uuid.UUID) ([]*domain.Schedule, error)
	GetByRoomAndWeekDay(ctx context.Context, roomID uuid.UUID, weekDay domain.WeekDay) (*domain.Schedule, error)
	ExistsForRoom(ctx context.Context, roomID uuid.UUID) (bool, error)
}

type scheduleRepository struct {
	db *sql.DB
}

func NewScheduleRepository(db *sql.DB) ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) Create(ctx context.Context, schedule *domain.Schedule) error {
	schedule.ID = uuid.New()
	query := `
		INSERT INTO schedules (id, room_id, week_day, start_time, end_time, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`
	_, err := r.db.ExecContext(ctx, query, schedule.ID, schedule.RoomID, schedule.WeekDay,
		schedule.StartTime, schedule.EndTime)
	return err
}

func (r *scheduleRepository) GetByRoomID(ctx context.Context, roomID uuid.UUID) ([]*domain.Schedule, error) {
	query := `
		SELECT id, room_id, week_day, start_time, end_time, created_at
		FROM schedules
		WHERE room_id = $1
		ORDER BY week_day, start_time
	`
	rows, err := r.db.QueryContext(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*domain.Schedule
	for rows.Next() {
		var schedule domain.Schedule
		err := rows.Scan(&schedule.ID, &schedule.RoomID, &schedule.WeekDay,
			&schedule.StartTime, &schedule.EndTime, &schedule.CreatedAt)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, &schedule)
	}
	return schedules, nil
}

func (r *scheduleRepository) GetByRoomAndWeekDay(ctx context.Context, roomID uuid.UUID, weekDay domain.WeekDay) (*domain.Schedule, error) {
	query := `
		SELECT id, room_id, week_day, start_time, end_time, created_at
		FROM schedules
		WHERE room_id = $1 AND week_day = $2
	`
	var schedule domain.Schedule
	err := r.db.QueryRowContext(ctx, query, roomID, weekDay).Scan(
		&schedule.ID, &schedule.RoomID, &schedule.WeekDay,
		&schedule.StartTime, &schedule.EndTime, &schedule.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepository) ExistsForRoom(ctx context.Context, roomID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM schedules WHERE room_id = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, roomID).Scan(&exists)
	return exists, err
}
