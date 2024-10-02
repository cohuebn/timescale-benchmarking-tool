package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type connectivityResult struct {
	SuccessfulConnection bool
	CpuUsageCount int
	Error error
}

// Run a query to prove that the connection to the database is working.
func RunConnectivityCheck(connectionPool *pgxpool.Pool) connectivityResult {
	connection, err := connectionPool.Acquire(context.Background())
	if (err != nil) {
		return connectivityResult{
			SuccessfulConnection: false,
			Error: err,
		}
	}
	defer connection.Release()
	var cpuUsageCount int
	err = connection.
		QueryRow(context.Background(), "select count(1) as cpu_usage_count from public.cpu_usage").
		Scan(&cpuUsageCount)
	
	if (err != nil) {
		return connectivityResult{
			SuccessfulConnection: false,
			Error: err,
		}
	}

	return connectivityResult{
		SuccessfulConnection: true,
		CpuUsageCount: cpuUsageCount,
	}
}