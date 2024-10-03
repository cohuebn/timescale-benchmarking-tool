package workers

import (
	"log/slog"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/queries"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO: In a production application, I'd add integration tests to test the interaction
// between the main thread, the workers, and the actual database. In the interest
// of simplicity, I'm leaving that out of scope for this assignment.

// Create a "CPU usage worker" that will run CPU usage queries and measure important metrics on each query
func createCpuUsageWorker(workerId int, connectionPool *pgxpool.Pool, requests <-chan queries.CpuUsageQueryParams, responses chan<- queries.QueryMeasurement) {
	slog.Info("Worker started", "workerId", workerId)
	defer close(responses)

	for request := range requests {
		measurement := queries.MeasureCpuUsageQuery(connectionPool, request)
		responses <- measurement
	}
}

// Measure CPU usage queries using a pool of workers to run those queries
// Return all recorded query measurements
func MeasureCpuUsageQueries(numberOfWorkers int, connectionPool *pgxpool.Pool, incomingQueryParameters <-chan queries.CpuUsageQueryParams) <-chan queries.QueryMeasurement {
	// Setup channels for workers to receive requests and send responses
	// Each worker gets its own channel to receive requests on
	requestChannels := make([]chan queries.CpuUsageQueryParams, numberOfWorkers)
	responses := make(chan queries.QueryMeasurement, 100)
	for _, requestChannel := range requestChannels {
		defer close(requestChannel)
	}
	defer close(responses)

	// Setup worker pools
	for i := 0; i < numberOfWorkers; i++ {
		go createCpuUsageWorker(i, connectionPool, requestChannels[i], responses)
	}

	// Assign incoming query parameters to workers and send them for processing
	// Aggregate measures as they are returned to avoid storing each measure in memory
	// Wait for all workers to finish processing before moving on
	workerAssigner := WorkerAssigner{
		numberOfWorkers: numberOfWorkers,
		assignedWorkers: make(map[string]int),
	}
	for queryParameters := range incomingQueryParameters {
		assignedWorker := workerAssigner.AssignHostToWorker(queryParameters.Hostname)
		requestChannels[assignedWorker] <- queryParameters
	}

	return responses
}