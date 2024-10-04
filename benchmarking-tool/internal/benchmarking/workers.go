package benchmarking

import (
	"log/slog"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO: In a production application, I'd add integration tests to test the interaction
// between the main thread, the workers, and the actual database. In the interest
// of simplicity, I'm leaving that out of scope for this assignment.

// Create a "CPU usage worker" that will run CPU usage queries and measure important metrics on each query
func createCpuUsageWorker(workerId int, connectionPool *pgxpool.Pool, requests <-chan database.CpuUsageQueryParams, responses chan<- QueryMeasurement) {
	slog.Debug("Worker started", "workerId", workerId)

	for request := range requests {
		measurement := MeasureCpuUsageQuery(connectionPool, request)
		responses <- measurement
	}

	slog.Debug("Worker finished", "workerId", workerId)
}

// Make a buffered channel for each worker to receive requests on
func makeWorkerChannels(numberOfWorkers int) []chan database.CpuUsageQueryParams {
	requestChannels := make([]chan database.CpuUsageQueryParams, numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		requestChannels[i] = make(chan database.CpuUsageQueryParams, 100)
	}
	return requestChannels
}

// Measure CPU usage queries using a pool of workers to run those queries
// Return all recorded query measurements
func MeasureCpuUsageQueries(numberOfWorkers int, connectionPool *pgxpool.Pool, incomingQueryParameters <-chan database.CpuUsageQueryParams) AggregatedCpuUsageResults {
	// Setup channels for workers to receive requests and send responses
	// Each worker gets its own channel to receive requests on
	requestChannels := makeWorkerChannels(numberOfWorkers)
	responses := make(chan QueryMeasurement, 100)
	for _, requestChannel := range requestChannels {
		defer close(requestChannel)
	}
	defer close(responses)

	// Setup worker pools
	for workerId := 0; workerId < numberOfWorkers; workerId++ {
		go createCpuUsageWorker(workerId, connectionPool, requestChannels[workerId], responses)
	}

	// Assign incoming query parameters to workers and send them for processing
	workerAssigner := WorkerAssigner{
		numberOfWorkers: numberOfWorkers,
		assignedWorkers: make(map[string]int),
	}
	expectedResponseCount := 0
	for queryParameters := range incomingQueryParameters {
		assignedWorker := workerAssigner.AssignHostToWorker(queryParameters.Hostname)
		expectedResponseCount++
		requestChannels[assignedWorker] <- queryParameters
	}
	// Wait for and aggregate all responses
	resultAggregator := NewResultAggregator()
	for responseToAwait := 1; responseToAwait <= expectedResponseCount; responseToAwait++ {
		response := <-responses
		resultAggregator.AggregateCpuMeasure(response)
	}

	return resultAggregator.CalculateAggregates()
}