package service

import (
	"context"
	"errors"

	"booking-service/internal/domain"
	"booking-service/internal/repository"

	"github.com/google/uuid"
)

type RoomService struct {
	roomRepo repository.RoomRepository
}

func NewRoomService(roomRepo repository.RoomRepository) *RoomService {
	return &RoomService{roomRepo: roomRepo}
}

func (s *RoomService) CreateRoom(ctx context.Context, room *domain.Room) error {
	if room.Name == "" {
		return errors.New("name is required")
	}
	if room.Capacity != nil && *room.Capacity < 0 {
		return errors.New("capacity must be non-negative")
	}
	return s.roomRepo.Create(ctx, room)
}

func (s *RoomService) GetRoom(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, ErrRoomNotFound
	}
	return room, nil
}

func (s *RoomService) ListRooms(ctx context.Context) ([]*domain.Room, error) {
	return s.roomRepo.List(ctx)
}
