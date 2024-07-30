package database

import (
	"context"
	"fmt"
	"log/slog"
	"message_processing-service/internal/config"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Db     *pgxpool.Pool
	log    *slog.Logger
	Config *config.Config
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, connString string, log *slog.Logger, cfg *config.Config) (*Postgres, error) {
	var err error

	pgOnce.Do(func() {
		var db *pgxpool.Pool
		db, err = pgxpool.New(ctx, connString)

		if err != nil {
			err = fmt.Errorf("unable to create connection pool: %w", err)
			return
		}

		pgInstance = &Postgres{db, log, cfg}

		if err = CreateTable(ctx, db, log, cfg); err != nil {
			return
		}
	})

	if err != nil {
		return nil, err
	}
	return pgInstance, nil
}

func CreateTable(ctx context.Context, db *pgxpool.Pool, log *slog.Logger, cfg *config.Config) error {
	_, err := db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS messages (
	id SERIAL PRIMARY KEY,
	content TEXT NOT NULL,
	status VARCHAR(10) NOT NULL, 
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	processed_at TIMESTAMPTZ)`)
	if err != nil {
		return fmt.Errorf("failed to create mesagges table")
	}

	_, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS Users (
	    id SERIAL PRIMARY KEY, 
	    email VARCHAR(100) UNIQUE NOT NULL, 
	    password VARCHAR(255) NOT NULL,
	    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	return nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.Db.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.Db.Close()
}
