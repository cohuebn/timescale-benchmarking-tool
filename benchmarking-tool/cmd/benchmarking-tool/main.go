package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/benchmarking"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/cli"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/csv"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/database"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/reporting"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Setup a connection pool used by the application using CLI arguments
func setupConnectionPool(cliArguments cli.CliArguments) *pgxpool.Pool {
	// Initialize the connection pool to the database
	connectionString := cli.CreateConnectionString(cliArguments)
	connectionPool, err := database.CreateConnectionPool(connectionString)
	if (err != nil) {
		// If we can't connect to the database, there's no point in continuing
		log.Panic(err)
	}
	return connectionPool
}

func exitOnUnhandledError() {
	if r := recover(); r != nil {
		slog.Error("Stopping benchmarking due to an unexpected error")
		os.Exit(1)
	}
}

func main() {
	defer exitOnUnhandledError()

	// Get user-provided CLI arguments
	cliArguments := cli.ParseCliArguments()
	slog.Info("Benchmarking tool started", "filename", cliArguments.Filename, "workers", cliArguments.Workers)
	
	// Initialize the connection pool to the database
	connectionPool := setupConnectionPool(cliArguments)
	defer connectionPool.Close()

	// Ensure database connectivity before benchmarking; if we can't connect, there's no point in continuing
	connectivityCheck := database.RunConnectivityCheck(connectionPool)
	if (connectivityCheck.Error != nil) {
		log.Panic(connectivityCheck.Error)
	}
	
	// Read the CSV file and stream its contents
	csvStream := csv.StreamCsvFile(cliArguments.Filename)
	// Process all rows and aggregate results
	results := benchmarking.ProcessCsv(cliArguments.Workers, connectionPool, csvStream)
	slog.Info("Benchmarking tool finished. Results below")
	reporting.LogCpuUsageResultsToConsole(results)
}