package repository

import (
	"context"
	"database/sql"

	"event-service/internal/model"
)

// VenueRepository отвечает за операции с площадками в базе данных.
type VenueRepository struct {
	db *sql.DB
}

// NewVenueRepository создаёт новый VenueRepository.
func NewVenueRepository(db *sql.DB) *VenueRepository {
	return &VenueRepository{db: db}
}

// GetAll возвращает все площадки, отсортированные по названию.
func (r *VenueRepository) GetAll(ctx context.Context) ([]model.Venue, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, address, capacity, created_at
		 FROM venues ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var venues []model.Venue
	for rows.Next() {
		var v model.Venue
		if err := rows.Scan(&v.ID, &v.Name, &v.Address, &v.Capacity, &v.CreatedAt); err != nil {
			return nil, err
		}
		venues = append(venues, v)
	}
	return venues, rows.Err()
}

// GetByID возвращает площадку по ID или nil, если не найдена.
func (r *VenueRepository) GetByID(ctx context.Context, id string) (*model.Venue, error) {
	var v model.Venue
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, address, capacity, created_at
		 FROM venues WHERE id = $1`, id,
	).Scan(&v.ID, &v.Name, &v.Address, &v.Capacity, &v.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// Create вставляет новую площадку и заполняет сгенерированные поля.
func (r *VenueRepository) Create(ctx context.Context, v *model.Venue) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO venues (name, address, capacity)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at`,
		v.Name, v.Address, v.Capacity,
	).Scan(&v.ID, &v.CreatedAt)
}

// Update обновляет существующую площадку.
func (r *VenueRepository) Update(ctx context.Context, v *model.Venue) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE venues SET name = $1, address = $2, capacity = $3 WHERE id = $4`,
		v.Name, v.Address, v.Capacity, v.ID,
	)
	return err
}

// Delete удаляет площадку по ID.
func (r *VenueRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM venues WHERE id = $1`, id)
	return err
}

// Exists проверяет существование площадки с данным ID.
func (r *VenueRepository) Exists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM venues WHERE id = $1)`, id,
	).Scan(&exists)
	return exists, err
}
