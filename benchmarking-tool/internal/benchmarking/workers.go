package benchmarking

import (
	"context"
	"log/slog"
	"sync"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

// TODO: In a production application, I'd add integration tests to test the interaction
// between the main thread, the workers, and the actual database. In the interest
// of simplicity, I'm leaving that out of scope for this assignment.

// Create a "CPU usage worker" that will run CPU usage queries and measure important metrics on each query
func runCpuUsageWorker(ctx context.Context, workerId int, connectionPool *pgxpool.Pool, requests <-chan database.CpuUsageQueryParams, responses chan<- QueryMeasurement) error {
	slog.Debug("Worker started", "workerId", workerId)

	for {
		select {
		case <-ctx.Done():
			slog.Debug("Stopping worker due to an external error", "workerId", workerId)
			return ctx.Err()
		case request, ok := <-requests:
			if (!ok) {
				slog.Debug("Worker finished processing all requests", "workerId", workerId)
				return nil
			}
			measurement := MeasureCpuUsageQuery(connectionPool, request)
			responses <- measurement
		}
	}
}

// Make a buffered channel for each worker to receive requests on
func makeWorkerChannels(numberOfWorkers int) []chan database.CpuUsageQueryParams {
	requestChannels := make([]chan database.CpuUsageQueryParams, numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		requestChannels[i] = make(chan database.CpuUsageQueryParams, 100)
	}
	return requestChannels
}

// Assign incoming query parameters to workers and send the request to that worker for processing
func assignIncomingRequests(incomingQueryParameters <-chan database.CpuUsageQueryParams, requestChannels []chan database.CpuUsageQueryParams) {
	workerAssigner := NewWorkerAssigner(len(requestChannels))
	for queryParameters := range incomingQueryParameters {
		assignedWorker := workerAssigner.AssignHostToWorker(queryParameters.Hostname)
		requestChannels[assignedWorker] <- queryParameters
	}
}

// Run a pool of workers to process incoming CPU usage queries
func runWorkerPool(
	ctx context.Context,
	numberOfWorkers int,
	connectionPool *pgxpool.Pool,
	incomingQueryParameters <-chan database.CpuUsageQueryParams,
	errGroup *errgroup.Group) <-chan QueryMeasurement {
	// Setup channels for workers to receive requests and send responses
	// Each worker gets its own channel to receive requests on
	requestChannels := makeWorkerChannels(numberOfWorkers)
	responses := make(chan QueryMeasurement, 100)

	// Setup workers
	var workerWaitGroup sync.WaitGroup
	for workerId := 0; workerId < numberOfWorkers; workerId++ {
		workerWaitGroup.Add(1)
		errGroup.Go(func () error {
			defer workerWaitGroup.Done()
			return runCpuUsageWorker(ctx, workerId, connectionPool, requestChannels[workerId], responses)
		})
	}

	// Send requests to workers for processing
	go func() {
		assignIncomingRequests(incomingQueryParameters, requestChannels)
		for _, requestChannel := range requestChannels {
			close(requestChannel)
		}
	}()

	// When all workers are done, close the responses channel
	go func() {
		workerWaitGroup.Wait()
		close(responses)
	}()

	return responses
}

// Run CPU usage queries using a pool of workers. Return all recorded query measurements.
func RunCpuUsageQueries(ctx context.Context, numberOfWorkers int, connectionPool *pgxpool.Pool, incomingQueryParameters <-chan database.CpuUsageQueryParams, errGroup *errgroup.Group) AggregatedCpuUsageResults {
	responses := runWorkerPool(ctx, numberOfWorkers, connectionPool, incomingQueryParameters, errGroup)

	// Aggregate all responses
	resultProgress := GetProgressBar()
	resultAggregator := NewResultAggregator()
	for response := range responses {
		resultAggregator.AggregateCpuMeasure(response)
		resultProgress.Add(1)
	}

	return resultAggregator.CalculateAggregates()
}