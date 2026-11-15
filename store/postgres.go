package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	DB *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, connString string) (*PostgresStore, error) {
	db, error := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("Create Postgres Pool: %w", error)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("Ping Postgres: %w", err)
	}

	return &PostgresStore{
		DB: db,
	}, nil
}

func (s *PostgresStore) Close() {
	s.DB.Close()
}
