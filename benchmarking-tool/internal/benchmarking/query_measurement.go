package benchmarking

import (
	"context"
	"time"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Used to store information about a single query run
// While query parameters are not used in the resulting summary outputs, keeping them
// to make debugging easier
type QueryMeasurement struct {
	// The parameters used to run the query
	Params database.CpuUsageQueryParams
	// Used to understand how long the query took
	QueryTime time.Duration
	// If the query failed, this will contain the error
	Error error
}

func MeasureCpuUsageQuery(connectionPool *pgxpool.Pool, queryParams database.CpuUsageQueryParams) QueryMeasurement {
	query := database.ConvertCpuUsageParamsToQuery(queryParams)
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