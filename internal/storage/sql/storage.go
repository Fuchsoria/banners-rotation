package sqlstorage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(ctx context.Context, connectionString string) (*Storage, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("cannot open db, %w", err)
	}

	return &Storage{db}, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("cannot connect to db, %w", err)
	}

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
