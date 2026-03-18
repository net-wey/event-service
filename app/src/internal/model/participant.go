package model

import "time"

// Participant представляет участника, зарегистрированного на мероприятие.
type Participant struct {
	ID           string    `json:"id"`
	EventID      string    `json:"event_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	RegisteredAt time.Time `json:"registered_at"`
}

// CreateParticipantRequest — тело запроса для регистрации участника.
type CreateParticipantRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateParticipantRequest — тело запроса для частичного обновления участника.
type UpdateParticipantRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
