package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/benchmarking"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/cli"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/csv"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/database"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/reporting"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

// Setup a connection pool used by the application using CLI arguments
func setupConnectionPool(cliArguments cli.CliArguments) *pgxpool.Pool {
	// Initialize the connection pool to the database
	connectionString := cli.CreateConnectionString(cliArguments)
	connectionPool, err := database.CreateConnectionPool(connectionString)
	if (err != nil) {
		// If we can't create a connection pool to the database, don't continue
		log.Panic(err)
	}
	return connectionPool
}

func exitOnUnhandledError() {
	if r := recover(); r != nil {
		slog.Error("Stopping benchmarking due to an unexpected error. See details above")
		os.Exit(1)
	}
}

func main() {
	defer exitOnUnhandledError()
	
	// Get user-provided CLI arguments
	cliArguments := cli.ParseCliArguments()
	slog.SetLogLoggerLevel(cliArguments.LogLevel)
	slog.Info("Benchmarking tool started", "filename", cliArguments.Filename, "workers", cliArguments.Workers)
	
	// Initialize the connection pool to the database
	connectionPool := setupConnectionPool(cliArguments)
	defer connectionPool.Close()

	// Ensure database connectivity before benchmarking; if we can't connect, don't continue
	connectivityCheck := database.RunConnectivityCheck(connectionPool)
	if (connectivityCheck.Error != nil) {
		log.Panic(connectivityCheck.Error)
	}
	
	streamingErrGroup, streamingErrContext := errgroup.WithContext(context.Background())
	// Open the CSV file and stream its contents
	csvStream, openFileErr := csv.StreamCsvFile(streamingErrContext, cliArguments.Filename, streamingErrGroup)
	if openFileErr != nil {
		log.Panic(openFileErr)
	}
	// Process all rows and aggregate results
	results := benchmarking.ProcessCsv(streamingErrContext, cliArguments.Workers, connectionPool, csvStream, streamingErrGroup)

	// If there's an error in context from other streams, stop processing
	if streamingErr := streamingErrGroup.Wait(); streamingErr != nil {
		log.Panic(streamingErr)
	}

	slog.Info("Benchmarking tool finished. Results below")
	reporting.LogCpuUsageResultsToConsole(results)
}