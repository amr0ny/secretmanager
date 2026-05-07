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
	_, err := db.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS secrets (
			id SERIAL PRIMARY KEY,
			secret_hash TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	return err
}

func (db *DB) InsertSecret(ctx context.Context, secretHash string) error {
	_, err := db.pool.Exec(ctx, `
		INSERT INTO secrets (secret_hash) VALUES ($1)
	`, secretHash)
	return err
}

func (db *DB) Close() {
	db.pool.Close()
}
