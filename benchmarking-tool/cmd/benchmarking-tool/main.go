package main

import (
	"fmt"
	"log"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/cli"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/database"
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

	connectivityCheck := database.RunConnectivityCheck(connectionPool)
	
	fmt.Println("This is a stub for the benchmarking-tool command; it is not yet implemented.")
	fmt.Printf("The filename is: %s and the number of workers is %d\n", cliArguments.Filename, cliArguments.Workers)
	fmt.Printf("Database connectivity check result: %t\n", connectivityCheck.SuccessfulConnection)
}