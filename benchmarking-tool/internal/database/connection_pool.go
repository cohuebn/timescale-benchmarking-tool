package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Create a new connection pool to the database using the provided connection string.
func CreateConnectionPool(connectionString string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), connectionString)
}