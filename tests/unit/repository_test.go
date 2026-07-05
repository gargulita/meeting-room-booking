package unit

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"booking-service/internal/domain"
	"booking-service/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func intPtr(v int) *int {
	return &v
}

func strPtr(v string) *string {
	return &v
}

func TestRoomRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewRoomRepository(db)
	room := &domain.Room{
		Name:        "Test Room",
		Description: strPtr("Test Description"),
		Capacity:    intPtr(10),
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO rooms (id, name, description, capacity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING created_at, updated_at
	`)).
		WithArgs(sqlmock.AnyArg(), room.Name, room.Description, room.Capacity).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
			AddRow(time.Now(), time.Now()))

	err = repo.Create(context.Background(), room)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewRoomRepository(db)
	roomID := uuid.New()
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, description, capacity, created_at, updated_at
		FROM rooms
		WHERE id = $1
	`)).
		WithArgs(roomID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "capacity", "created_at", "updated_at"}).
			AddRow(roomID, "Test Room", "Test Description", 10, now, now))

	room, err := repo.GetByID(context.Background(), roomID)
	assert.NoError(t, err)
	require.NotNil(t, room)
	assert.Equal(t, roomID, room.ID)
	assert.Equal(t, "Test Room", room.Name)
	require.NotNil(t, room.Description)
	assert.Equal(t, "Test Description", *room.Description)
	require.NotNil(t, room.Capacity)
	assert.Equal(t, 10, *room.Capacity)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewRoomRepository(db)
	roomID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, description, capacity, created_at, updated_at
		FROM rooms
		WHERE id = $1
	`)).
		WithArgs(roomID).
		WillReturnError(sql.ErrNoRows)

	room, err := repo.GetByID(context.Background(), roomID)
	assert.NoError(t, err)
	assert.Nil(t, room)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewRoomRepository(db)
	now := time.Now()
	roomID1 := uuid.New()
	roomID2 := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, description, capacity, created_at, updated_at
		FROM rooms
		ORDER BY created_at DESC
	`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "capacity", "created_at", "updated_at"}).
			AddRow(roomID1, "Room 1", "Desc 1", 5, now, now).
			AddRow(roomID2, "Room 2", nil, nil, now, now))

	rooms, err := repo.List(context.Background())
	assert.NoError(t, err)
	require.Len(t, rooms, 2)

	assert.Equal(t, roomID1, rooms[0].ID)
	assert.Equal(t, "Room 1", rooms[0].Name)
	require.NotNil(t, rooms[0].Description)
	assert.Equal(t, "Desc 1", *rooms[0].Description)
	require.NotNil(t, rooms[0].Capacity)
	assert.Equal(t, 5, *rooms[0].Capacity)

	assert.Equal(t, roomID2, rooms[1].ID)
	assert.Equal(t, "Room 2", rooms[1].Name)
	assert.Nil(t, rooms[1].Description)
	assert.Nil(t, rooms[1].Capacity)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewRoomRepository(db)

	roomID := uuid.New()
	room := &domain.Room{
		ID:          roomID,
		Name:        "Updated Room",
		Description: strPtr("Updated Description"),
		Capacity:    intPtr(20),
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE rooms
		SET name = $1, description = $2, capacity = $3, updated_at = NOW()
		WHERE id = $4
	`)).
		WithArgs(room.Name, room.Description, room.Capacity, room.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(context.Background(), room)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewRoomRepository(db)
	roomID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM rooms WHERE id = $1`)).
		WithArgs(roomID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), roomID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
