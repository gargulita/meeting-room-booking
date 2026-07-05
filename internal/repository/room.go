package repository

import (
	"context"
	"database/sql"

	"booking-service/internal/domain"

	"github.com/google/uuid"
)

type RoomRepository interface {
	Create(ctx context.Context, room *domain.Room) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Room, error)
	List(ctx context.Context) ([]*domain.Room, error)
	Update(ctx context.Context, room *domain.Room) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type roomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *domain.Room) error {
	room.ID = uuid.New()
	query := `
		INSERT INTO rooms (id, name, description, capacity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query, room.ID, room.Name, room.Description, room.Capacity).
		Scan(&room.CreatedAt, &room.UpdatedAt)
}

func (r *roomRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	query := `
		SELECT id, name, description, capacity, created_at, updated_at
		FROM rooms
		WHERE id = $1
	`

	var room domain.Room
	var description sql.NullString
	var capacity sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&room.ID,
		&room.Name,
		&description,
		&capacity,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if description.Valid {
		room.Description = &description.String
	}
	if capacity.Valid {
		v := int(capacity.Int64)
		room.Capacity = &v
	}

	return &room, nil
}

func (r *roomRepository) List(ctx context.Context) ([]*domain.Room, error) {
	query := `
		SELECT id, name, description, capacity, created_at, updated_at
		FROM rooms
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*domain.Room
	for rows.Next() {
		var room domain.Room
		var description sql.NullString
		var capacity sql.NullInt64

		err := rows.Scan(
			&room.ID,
			&room.Name,
			&description,
			&capacity,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			room.Description = &description.String
		}
		if capacity.Valid {
			v := int(capacity.Int64)
			room.Capacity = &v
		}

		rooms = append(rooms, &room)
	}
	return rooms, rows.Err()
}

func (r *roomRepository) Update(ctx context.Context, room *domain.Room) error {
	query := `
		UPDATE rooms
		SET name = $1, description = $2, capacity = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, room.Name, room.Description, room.Capacity, room.ID)
	return err
}

func (r *roomRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM rooms WHERE id = $1`, id)
	return err
}
