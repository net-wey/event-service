package repository

import (
	"context"
	"database/sql"

	"event-service/internal/model"
)

// EventRepository отвечает за операции с мероприятиями в базе данных.
type EventRepository struct {
	db *sql.DB
}

// NewEventRepository создаёт новый EventRepository.
func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

// GetAll возвращает все мероприятия, отсортированные по дате.
func (r *EventRepository) GetAll(ctx context.Context) ([]model.Event, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, description, location, event_date, created_at, updated_at
		 FROM events ORDER BY event_date`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var e model.Event
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Location,
			&e.EventDate, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

// GetByID возвращает мероприятие по ID или nil, если не найдено.
func (r *EventRepository) GetByID(ctx context.Context, id string) (*model.Event, error) {
	var e model.Event
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, description, location, event_date, created_at, updated_at
		 FROM events WHERE id = $1`, id,
	).Scan(&e.ID, &e.Title, &e.Description, &e.Location,
		&e.EventDate, &e.CreatedAt, &e.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// Create вставляет новое мероприятие и заполняет сгенерированные поля.
func (r *EventRepository) Create(ctx context.Context, e *model.Event) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO events (title, description, location, event_date)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at, updated_at`,
		e.Title, e.Description, e.Location, e.EventDate,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

// Update обновляет существующее мероприятие.
func (r *EventRepository) Update(ctx context.Context, e *model.Event) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE events SET title = $1, description = $2, location = $3,
		 event_date = $4, updated_at = NOW()
		 WHERE id = $5`,
		e.Title, e.Description, e.Location, e.EventDate, e.ID,
	)
	return err
}

// Delete удаляет мероприятие по ID.
func (r *EventRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM events WHERE id = $1`, id)
	return err
}

// Exists проверяет существование мероприятия с данным ID.
func (r *EventRepository) Exists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)`, id,
	).Scan(&exists)
	return exists, err
}
