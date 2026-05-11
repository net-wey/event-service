package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"event-service/internal/model"
)

type mockEventRepo struct {
	getAllFn  func(context.Context) ([]model.Event, error)
	getByIDFn func(context.Context, string) (*model.Event, error)
	existsFn  func(context.Context, string) (bool, error)
	createFn  func(context.Context, *model.Event) error
	updateFn  func(context.Context, *model.Event) error
	deleteFn  func(context.Context, string) error
}

func (m mockEventRepo) GetAll(ctx context.Context) ([]model.Event, error) { return m.getAllFn(ctx) }
func (m mockEventRepo) GetByID(ctx context.Context, id string) (*model.Event, error) {
	return m.getByIDFn(ctx, id)
}
func (m mockEventRepo) Exists(ctx context.Context, id string) (bool, error) {
	return m.existsFn(ctx, id)
}
func (m mockEventRepo) Create(ctx context.Context, event *model.Event) error {
	return m.createFn(ctx, event)
}
func (m mockEventRepo) Update(ctx context.Context, event *model.Event) error {
	return m.updateFn(ctx, event)
}
func (m mockEventRepo) Delete(ctx context.Context, id string) error { return m.deleteFn(ctx, id) }

type mockParticipantRepo struct {
	getByEventIDFn func(context.Context, string) ([]model.Participant, error)
	createFn       func(context.Context, *model.Participant) error
	getByIDFn      func(context.Context, string) (*model.Participant, error)
	updateFn       func(context.Context, *model.Participant) error
	deleteFn       func(context.Context, string) error
}

func (m mockParticipantRepo) GetByEventID(ctx context.Context, eventID string) ([]model.Participant, error) {
	return m.getByEventIDFn(ctx, eventID)
}
func (m mockParticipantRepo) Create(ctx context.Context, p *model.Participant) error {
	return m.createFn(ctx, p)
}
func (m mockParticipantRepo) GetByID(ctx context.Context, id string) (*model.Participant, error) {
	return m.getByIDFn(ctx, id)
}
func (m mockParticipantRepo) Update(ctx context.Context, p *model.Participant) error {
	return m.updateFn(ctx, p)
}
func (m mockParticipantRepo) Delete(ctx context.Context, id string) error { return m.deleteFn(ctx, id) }

type mockVenueRepo struct {
	getAllFn  func(context.Context) ([]model.Venue, error)
	getByIDFn func(context.Context, string) (*model.Venue, error)
	createFn  func(context.Context, *model.Venue) error
	updateFn  func(context.Context, *model.Venue) error
	deleteFn  func(context.Context, string) error
	existsFn  func(context.Context, string) (bool, error)
}

func (m mockVenueRepo) GetAll(ctx context.Context) ([]model.Venue, error) { return m.getAllFn(ctx) }
func (m mockVenueRepo) GetByID(ctx context.Context, id string) (*model.Venue, error) {
	return m.getByIDFn(ctx, id)
}
func (m mockVenueRepo) Create(ctx context.Context, v *model.Venue) error { return m.createFn(ctx, v) }
func (m mockVenueRepo) Update(ctx context.Context, v *model.Venue) error { return m.updateFn(ctx, v) }
func (m mockVenueRepo) Delete(ctx context.Context, id string) error      { return m.deleteFn(ctx, id) }
func (m mockVenueRepo) Exists(ctx context.Context, id string) (bool, error) {
	return m.existsFn(ctx, id)
}

func TestEventServiceCreate(t *testing.T) {
	t.Run("invalid date", func(t *testing.T) {
		svc := NewEventService(mockEventRepo{})
		_, err := svc.Create(context.Background(), model.CreateEventRequest{Title: "x", EventDate: "bad"})
		if !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		called := false
		svc := NewEventService(mockEventRepo{
			createFn: func(_ context.Context, event *model.Event) error {
				called = true
				if event.Title != "GoConf" {
					t.Fatalf("unexpected title: %s", event.Title)
				}
				return nil
			},
		})

		res, err := svc.Create(context.Background(), model.CreateEventRequest{
			Title:       "GoConf",
			Description: "d",
			Location:    "Moscow",
			EventDate:   "2026-05-11T10:00:00Z",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !called || res == nil {
			t.Fatalf("expected event to be created")
		}
	})
}

func TestEventServiceDelete_NotFound(t *testing.T) {
	svc := NewEventService(mockEventRepo{
		existsFn: func(context.Context, string) (bool, error) { return false, nil },
	})

	err := svc.Delete(context.Background(), "id")
	if !errors.Is(err, ErrEventNotFound) {
		t.Fatalf("expected ErrEventNotFound, got %v", err)
	}
}

func TestParticipantServiceRegister(t *testing.T) {
	t.Run("event not found", func(t *testing.T) {
		svc := NewParticipantService(
			mockParticipantRepo{},
			mockEventRepo{existsFn: func(context.Context, string) (bool, error) { return false, nil }},
		)
		_, err := svc.Register(context.Background(), "e1", model.CreateParticipantRequest{Name: "n", Email: "e@x.com"})
		if !errors.Is(err, ErrEventNotFound) {
			t.Fatalf("expected ErrEventNotFound, got %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		created := false
		svc := NewParticipantService(
			mockParticipantRepo{
				createFn: func(_ context.Context, p *model.Participant) error {
					created = true
					if p.EventID != "e1" {
						t.Fatalf("unexpected event id: %s", p.EventID)
					}
					return nil
				},
			},
			mockEventRepo{existsFn: func(context.Context, string) (bool, error) { return true, nil }},
		)

		_, err := svc.Register(context.Background(), "e1", model.CreateParticipantRequest{Name: "n", Email: "e@x.com"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !created {
			t.Fatalf("expected create call")
		}
	})
}

func TestVenueServiceGetAndUpdate(t *testing.T) {
	t.Run("get not found", func(t *testing.T) {
		svc := NewVenueService(mockVenueRepo{getByIDFn: func(context.Context, string) (*model.Venue, error) { return nil, nil }})
		_, err := svc.Get(context.Background(), "v1")
		if !errors.Is(err, ErrVenueNotFound) {
			t.Fatalf("expected ErrVenueNotFound, got %v", err)
		}
	})

	t.Run("update success", func(t *testing.T) {
		name := "new"
		updated := false
		svc := NewVenueService(mockVenueRepo{
			getByIDFn: func(context.Context, string) (*model.Venue, error) {
				return &model.Venue{ID: "v1", Name: "old"}, nil
			},
			updateFn: func(_ context.Context, v *model.Venue) error {
				updated = true
				if v.Name != "new" {
					t.Fatalf("unexpected name: %s", v.Name)
				}
				return nil
			},
		})
		_, err := svc.Update(context.Background(), "v1", model.UpdateVenueRequest{Name: &name})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !updated {
			t.Fatalf("expected update call")
		}
	})
}

func TestEventServiceUpdate_InvalidDate(t *testing.T) {
	badDate := "2026/10/01"
	svc := NewEventService(mockEventRepo{
		getByIDFn: func(context.Context, string) (*model.Event, error) {
			return &model.Event{ID: "e1", EventDate: time.Now()}, nil
		},
	})

	_, err := svc.Update(context.Background(), "e1", model.UpdateEventRequest{EventDate: &badDate})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestEventServiceGetAndDeleteSuccess(t *testing.T) {
	svc := NewEventService(mockEventRepo{
		getByIDFn: func(context.Context, string) (*model.Event, error) {
			return &model.Event{ID: "e1", Title: "T"}, nil
		},
		existsFn: func(context.Context, string) (bool, error) { return true, nil },
		deleteFn: func(context.Context, string) error { return nil },
	})

	_, err := svc.Get(context.Background(), "e1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := svc.Delete(context.Background(), "e1"); err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}
}

func TestParticipantServiceListUpdateDelete(t *testing.T) {
	newName := "Neo"
	newEmail := "neo@example.com"
	svc := NewParticipantService(
		mockParticipantRepo{
			getByEventIDFn: func(context.Context, string) ([]model.Participant, error) {
				return []model.Participant{{ID: "p1"}}, nil
			},
			getByIDFn: func(context.Context, string) (*model.Participant, error) {
				return &model.Participant{ID: "p1", Name: "old", Email: "old@example.com"}, nil
			},
			updateFn: func(context.Context, *model.Participant) error { return nil },
			deleteFn: func(context.Context, string) error { return nil },
		},
		mockEventRepo{existsFn: func(context.Context, string) (bool, error) { return true, nil }},
	)

	list, err := svc.ListByEvent(context.Background(), "e1")
	if err != nil || len(list) != 1 {
		t.Fatalf("unexpected list result, err=%v len=%d", err, len(list))
	}
	_, err = svc.Update(context.Background(), "p1", model.UpdateParticipantRequest{Name: &newName, Email: &newEmail})
	if err != nil {
		t.Fatalf("unexpected update error: %v", err)
	}
	if err := svc.Delete(context.Background(), "p1"); err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}
}

func TestVenueServiceListCreateDelete(t *testing.T) {
	svc := NewVenueService(mockVenueRepo{
		getAllFn: func(context.Context) ([]model.Venue, error) { return []model.Venue{{ID: "v1"}}, nil },
		createFn: func(context.Context, *model.Venue) error { return nil },
		existsFn: func(context.Context, string) (bool, error) { return true, nil },
		deleteFn: func(context.Context, string) error { return nil },
	})

	venues, err := svc.List(context.Background())
	if err != nil || len(venues) != 1 {
		t.Fatalf("unexpected list result, err=%v len=%d", err, len(venues))
	}
	_, err = svc.Create(context.Background(), model.CreateVenueRequest{Name: "Hall", Capacity: 100})
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	if err := svc.Delete(context.Background(), "v1"); err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}
}
