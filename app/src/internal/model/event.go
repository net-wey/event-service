package model

import "time"

// Event представляет запланированное мероприятие.
type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	EventDate   time.Time `json:"event_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateEventRequest — тело запроса для создания мероприятия.
type CreateEventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	EventDate   string `json:"event_date"`
}

// UpdateEventRequest — тело запроса для частичного обновления мероприятия.
type UpdateEventRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
	EventDate   *string `json:"event_date"`
}
