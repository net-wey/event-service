package service

import (
	"context"
	"errors"
	"time"

	"event-service/internal/model"
)

var (
	ErrEventNotFound = errors.New("мероприятие не найдено")
	ErrInvalidInput  = errors.New("некорректные входные данные")
)

// EventService содержит бизнес-логику для работы с мероприятиями.
type EventService struct {
	repo eventRepo
}

// NewEventService создаёт новый EventService.
func NewEventService(repo eventRepo) *EventService {
	return &EventService{repo: repo}
}

// List возвращает все мероприятия.
func (s *EventService) List(ctx context.Context) ([]model.Event, error) {
	return s.repo.GetAll(ctx)
}

// Get возвращает мероприятие по ID.
func (s *EventService) Get(ctx context.Context, id string) (*model.Event, error) {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, ErrEventNotFound
	}
	return event, nil
}

// Create валидирует входные данные и создаёт новое мероприятие.
func (s *EventService) Create(ctx context.Context, req model.CreateEventRequest) (*model.Event, error) {
	if req.Title == "" || req.EventDate == "" {
		return nil, ErrInvalidInput
	}

	eventDate, err := time.Parse(time.RFC3339, req.EventDate)
	if err != nil {
		return nil, ErrInvalidInput
	}

	event := &model.Event{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		EventDate:   eventDate,
	}

	if err := s.repo.Create(ctx, event); err != nil {
		return nil, err
	}
	return event, nil
}

// Update частично обновляет существующее мероприятие.
func (s *EventService) Update(ctx context.Context, id string, req model.UpdateEventRequest) (*model.Event, error) {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, ErrEventNotFound
	}

	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.Location != nil {
		event.Location = *req.Location
	}
	if req.EventDate != nil {
		eventDate, err := time.Parse(time.RFC3339, *req.EventDate)
		if err != nil {
			return nil, ErrInvalidInput
		}
		event.EventDate = eventDate
	}

	if err := s.repo.Update(ctx, event); err != nil {
		return nil, err
	}
	return event, nil
}

// Delete удаляет мероприятие по ID.
func (s *EventService) Delete(ctx context.Context, id string) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrEventNotFound
	}
	return s.repo.Delete(ctx, id)
}
