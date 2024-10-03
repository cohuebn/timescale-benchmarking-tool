package queries

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Used to store information about a single query run
type QueryMeasurement struct {
	// The parameters used to run the query
	Params CpuUsageQueryParams
	// Used to understand how long the query took
	QueryTime time.Duration
	// If the query failed, this will contain the error
	Error error
}

func MeasureCpuUsageQuery(connectionPool *pgxpool.Pool, queryParams CpuUsageQueryParams) QueryMeasurement {
	query := ConvertCpuUsageParamsToQuery(queryParams)
	connection, connectionError := connectionPool.Acquire(context.Background())
	if connectionError != nil {
		return QueryMeasurement{Error: connectionError, QueryTime: 0}
	}
	defer connection.Release()

	queryStartTime := time.Now()
	_, queryError := connection.Query(context.Background(), query)
	return QueryMeasurement{
		Params: 	queryParams,
		QueryTime: time.Since(queryStartTime),
		Error: 		queryError,
	}
}