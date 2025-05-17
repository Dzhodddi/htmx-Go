package store

import (
	"context"
	"database/sql"
	"log"
)

type RolesStorage struct {
	db *sql.DB
}

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       string `json:"level"`
}

func (s *RolesStorage) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `SELECT id, name, description, level FROM roles WHERE name = $1`
	log.Printf(name)
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDelay)
	defer cancel()
	role := &Role{}
	err := s.db.QueryRowContext(ctx, query, name).Scan(&role.ID,
		&role.Name, &role.Description, &role.Level)
	if err != nil {
		return nil, err
	}
	log.Printf(role.Level)
	return role, nil
}
