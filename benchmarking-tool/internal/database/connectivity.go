package database

import (
	"context"

	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type connectivityResult struct {
	SuccessfulConnection bool
	LatestCpuUsageTimestamp time.Time
	Error error
}

// Run a query against the cpu_usage_count table to prove that the connection to the database is working.
func RunConnectivityCheck(connectionPool *pgxpool.Pool) connectivityResult {
	connection, err := connectionPool.Acquire(context.Background())
	if (err != nil) {
		return connectivityResult{
			SuccessfulConnection: false,
			Error: err,
		}
	}
	defer connection.Release()
	var latestTimestamp time.Time
	err = connection.
		QueryRow(context.Background(), "select * from public.cpu_usage order by ts desc limit 1;").
		Scan(&latestTimestamp)
	
	if (err != nil) {
		return connectivityResult{
			SuccessfulConnection: false,
			Error: err,
		}
	}

	return connectivityResult{
		SuccessfulConnection: true,
		LatestCpuUsageTimestamp: latestTimestamp,
	}
}