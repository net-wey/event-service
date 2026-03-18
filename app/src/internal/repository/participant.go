package repository

import (
	"context"
	"database/sql"

	"event-service/internal/model"
)

// ParticipantRepository отвечает за операции с участниками в базе данных.
type ParticipantRepository struct {
	db *sql.DB
}

// NewParticipantRepository создаёт новый ParticipantRepository.
func NewParticipantRepository(db *sql.DB) *ParticipantRepository {
	return &ParticipantRepository{db: db}
}

// GetByEventID возвращает всех участников для данного мероприятия.
func (r *ParticipantRepository) GetByEventID(ctx context.Context, eventID string) ([]model.Participant, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, event_id, name, email, registered_at
		 FROM participants WHERE event_id = $1 ORDER BY registered_at`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []model.Participant
	for rows.Next() {
		var p model.Participant
		if err := rows.Scan(&p.ID, &p.EventID, &p.Name, &p.Email, &p.RegisteredAt); err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, rows.Err()
}

// GetByID возвращает участника по ID или nil, если не найден.
func (r *ParticipantRepository) GetByID(ctx context.Context, id string) (*model.Participant, error) {
	var p model.Participant
	err := r.db.QueryRowContext(ctx,
		`SELECT id, event_id, name, email, registered_at
		 FROM participants WHERE id = $1`, id,
	).Scan(&p.ID, &p.EventID, &p.Name, &p.Email, &p.RegisteredAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Create вставляет нового участника и заполняет сгенерированные поля.
func (r *ParticipantRepository) Create(ctx context.Context, p *model.Participant) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO participants (event_id, name, email)
		 VALUES ($1, $2, $3)
		 RETURNING id, registered_at`,
		p.EventID, p.Name, p.Email,
	).Scan(&p.ID, &p.RegisteredAt)
}

// Update обновляет существующего участника.
func (r *ParticipantRepository) Update(ctx context.Context, p *model.Participant) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE participants SET name = $1, email = $2 WHERE id = $3`,
		p.Name, p.Email, p.ID,
	)
	return err
}

// Delete удаляет участника по ID.
func (r *ParticipantRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM participants WHERE id = $1`, id)
	return err
}
