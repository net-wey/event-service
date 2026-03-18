package model

import "time"

// Venue представляет площадку для проведения мероприятий.
type Venue struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateVenueRequest — тело запроса для создания площадки.
type CreateVenueRequest struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity int    `json:"capacity"`
}

// UpdateVenueRequest — тело запроса для частичного обновления площадки.
type UpdateVenueRequest struct {
	Name     *string `json:"name"`
	Address  *string `json:"address"`
	Capacity *int    `json:"capacity"`
}
