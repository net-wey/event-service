package service

import (
	"context"

	"event-service/internal/model"
)

type eventReader interface {
	GetByID(ctx context.Context, id string) (*model.Event, error)
	Exists(ctx context.Context, id string) (bool, error)
}

type eventRepo interface {
	eventReader
	GetAll(ctx context.Context) ([]model.Event, error)
	Create(ctx context.Context, event *model.Event) error
	Update(ctx context.Context, event *model.Event) error
	Delete(ctx context.Context, id string) error
}

type participantRepo interface {
	GetByEventID(ctx context.Context, eventID string) ([]model.Participant, error)
	Create(ctx context.Context, participant *model.Participant) error
	GetByID(ctx context.Context, id string) (*model.Participant, error)
	Update(ctx context.Context, participant *model.Participant) error
	Delete(ctx context.Context, id string) error
}

type venueRepo interface {
	GetAll(ctx context.Context) ([]model.Venue, error)
	GetByID(ctx context.Context, id string) (*model.Venue, error)
	Create(ctx context.Context, venue *model.Venue) error
	Update(ctx context.Context, venue *model.Venue) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}
