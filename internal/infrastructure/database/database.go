package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, connString string) (*DB, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	db := &DB{pool: pool}
	if err := db.ensureWorkingTable(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ensure working table: %w", err)
	}
	return db, nil
}

func (db *DB) ensureWorkingTable(ctx context.Context) error {
	if _, err := db.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS secrets (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			secret_hash TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`); err != nil {
		return err
	}
	if _, err := db.pool.Exec(ctx, `
		ALTER TABLE secrets ADD COLUMN IF NOT EXISTS name TEXT
	`); err != nil {
		return err
	}
	if _, err := db.pool.Exec(ctx, `
		UPDATE secrets SET name = 'legacy-' || id::text WHERE name IS NULL
	`); err != nil {
		return err
	}
	if _, err := db.pool.Exec(ctx, `
		ALTER TABLE secrets ALTER COLUMN name SET NOT NULL
	`); err != nil {
		return err
	}
	if _, err := db.pool.Exec(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS secrets_name_unique ON secrets (name)
	`); err != nil {
		return err
	}
	return nil
}

func (db *DB) InsertSecret(ctx context.Context, name, secretHash string) error {
	_, err := db.pool.Exec(ctx, `
		INSERT INTO secrets (name, secret_hash) VALUES ($1, $2)
	`, name, secretHash)
	return err
}

func (db *DB) Close() {
	db.pool.Close()
}
