package main

import (
	"log"
	"log/slog"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/cli"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/csv"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/database"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/queries"
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

func main() {
	cliArguments := cli.ParseCliArguments()
	
	// Initialize the connection pool to the database
	connectionPool := setupConnectionPool(cliArguments)
	defer connectionPool.Close()

	// Read the CSV file and stream its contents
	csvStream := csv.StreamCsvFile(cliArguments.Filename)
	// For now, just get a count of query params in the CSV and run a single query; this is
	// just a sanity check to test integration but will be removed later
	var headerRow []string
	var firstQueryResult *queries.QueryMeasurement
	queryParamsCount := 0
	for csvRow := range csvStream {
		if csvRow.Error != nil {
			log.Panic(csvRow.Error)
		}
		if headerRow == nil {
			headerRow = csvRow.Row
			// TODO - get rid of continue control flow by breaking out functions
			// for different paths
			continue
		} 
		
		queryParams := queries.ParseCpuUsageCsvRow(headerRow, csvRow.Row)
		if (firstQueryResult == nil) {
			measurement := queries.MeasureCpuUsageQuery(connectionPool, queries.ParseCpuUsageCsvRow(headerRow, csvRow.Row))
			firstQueryResult = &measurement
		}
		slog.Debug("Parsed query params", "queryParams", queryParams)
		queryParamsCount++
	}

	connectivityCheck := database.RunConnectivityCheck(connectionPool)

	slog.Info("Benchmarking tool started", "filename", cliArguments.Filename, "workers", cliArguments.Workers)
	slog.Info("CSV file read and parsed", "parsedRowCount", queryParamsCount)
	slog.Info("Database connectivity check", "successfulConnection", connectivityCheck.SuccessfulConnection)
	slog.Info("First query result", "queryResult", *firstQueryResult)
}