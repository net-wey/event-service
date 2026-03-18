package service

import (
	"context"
	"errors"

	"event-service/internal/model"
	"event-service/internal/repository"
)

var ErrVenueNotFound = errors.New("площадка не найдена")

// VenueService содержит бизнес-логику для работы с площадками.
type VenueService struct {
	repo *repository.VenueRepository
}

// NewVenueService создаёт новый VenueService.
func NewVenueService(repo *repository.VenueRepository) *VenueService {
	return &VenueService{repo: repo}
}

// List возвращает все площадки.
func (s *VenueService) List(ctx context.Context) ([]model.Venue, error) {
	return s.repo.GetAll(ctx)
}

// Get возвращает площадку по ID.
func (s *VenueService) Get(ctx context.Context, id string) (*model.Venue, error) {
	venue, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if venue == nil {
		return nil, ErrVenueNotFound
	}
	return venue, nil
}

// Create валидирует входные данные и создаёт новую площадку.
func (s *VenueService) Create(ctx context.Context, req model.CreateVenueRequest) (*model.Venue, error) {
	if req.Name == "" {
		return nil, ErrInvalidInput
	}

	venue := &model.Venue{
		Name:     req.Name,
		Address:  req.Address,
		Capacity: req.Capacity,
	}

	if err := s.repo.Create(ctx, venue); err != nil {
		return nil, err
	}
	return venue, nil
}

// Update частично обновляет существующую площадку.
func (s *VenueService) Update(ctx context.Context, id string, req model.UpdateVenueRequest) (*model.Venue, error) {
	venue, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if venue == nil {
		return nil, ErrVenueNotFound
	}

	if req.Name != nil {
		venue.Name = *req.Name
	}
	if req.Address != nil {
		venue.Address = *req.Address
	}
	if req.Capacity != nil {
		venue.Capacity = *req.Capacity
	}

	if err := s.repo.Update(ctx, venue); err != nil {
		return nil, err
	}
	return venue, nil
}

// Delete удаляет площадку по ID.
func (s *VenueService) Delete(ctx context.Context, id string) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrVenueNotFound
	}
	return s.repo.Delete(ctx, id)
}
