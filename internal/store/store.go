package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct{ DB *sql.DB }

func Open(ctx context.Context, dsn string) (*Store, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil { return nil, err }
	db.SetMaxOpenConns(10); db.SetMaxIdleConns(5); db.SetConnMaxLifetime(30*time.Minute)
	if err := db.PingContext(ctx); err != nil { _ = db.Close(); return nil, fmt.Errorf("database ping: %w", err) }
	return &Store{DB: db}, nil
}

func (s *Store) Close() error { return s.DB.Close() }
