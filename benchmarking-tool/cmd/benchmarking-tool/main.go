package main

import (
	"fmt"
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
	// For now, just get a count of query params in the CSV; this is
	// just a sanity check to test integration but will be removed later
	var headerRow []string
	queryCount := 0
	for csvRow := range csvStream {
		if csvRow.Error != nil {
			log.Panic(csvRow.Error)
		}
		if headerRow == nil {
			headerRow = csvRow.Row
		} else {
			queryParams := queries.ParseCpuUsageCsvRow(headerRow, csvRow.Row)
			slog.Debug("Parsed query params", "queryParams", queryParams)
			queryCount++
		}
	}

	connectivityCheck := database.RunConnectivityCheck(connectionPool)

	fmt.Println("This is a stub for the benchmarking-tool command; it is not yet implemented.")
	fmt.Printf("Running with %d worker(s)\n", cliArguments.Workers)
	fmt.Printf("Read the CSV file: %s. Would've run %d queries\n", cliArguments.Filename, queryCount)
	fmt.Printf("Database connectivity check result: %t\n", connectivityCheck.SuccessfulConnection)
}