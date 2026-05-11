package service

import (
	"context"
	"errors"

	"event-service/internal/model"
)

var (
	ErrParticipantNotFound = errors.New("участник не найден")
)

// ParticipantService содержит бизнес-логику для работы с участниками.
type ParticipantService struct {
	repo      participantRepo
	eventRepo eventReader
}

// NewParticipantService создаёт новый ParticipantService.
func NewParticipantService(repo participantRepo, eventRepo eventReader) *ParticipantService {
	return &ParticipantService{repo: repo, eventRepo: eventRepo}
}

// ListByEvent возвращает всех участников мероприятия.
func (s *ParticipantService) ListByEvent(ctx context.Context, eventID string) ([]model.Participant, error) {
	exists, err := s.eventRepo.Exists(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrEventNotFound
	}
	return s.repo.GetByEventID(ctx, eventID)
}

// Register добавляет нового участника на мероприятие.
func (s *ParticipantService) Register(ctx context.Context, eventID string, req model.CreateParticipantRequest) (*model.Participant, error) {
	if req.Name == "" || req.Email == "" {
		return nil, ErrInvalidInput
	}

	exists, err := s.eventRepo.Exists(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrEventNotFound
	}

	p := &model.Participant{
		EventID: eventID,
		Name:    req.Name,
		Email:   req.Email,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// Update частично обновляет участника.
func (s *ParticipantService) Update(ctx context.Context, id string, req model.UpdateParticipantRequest) (*model.Participant, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrParticipantNotFound
	}

	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Email != nil {
		p.Email = *req.Email
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// Delete удаляет участника по ID.
func (s *ParticipantService) Delete(ctx context.Context, id string) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return ErrParticipantNotFound
	}
	return s.repo.Delete(ctx, id)
}
