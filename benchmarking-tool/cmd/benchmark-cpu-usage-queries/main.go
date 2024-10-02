package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/cli"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/database"
)

func main() {
	cliArguments := cli.ParseCliArguments()
	
	// Initialize the connection pool to the database
	connectionString := os.Getenv("DATABASE_CONNECTION_STRING")
	connectionPool, err := database.CreateConnectionPool(connectionString)
	if (err != nil) {
		log.Panic(err)
	}
	defer connectionPool.Close()

	connectivityCheck := database.RunConnectivityCheck(connectionPool)
	
	fmt.Println("This is a stub for the benchmark-cpu-usage-queries command; it is not yet implemented.")
	fmt.Printf("The filename is: %s and the number of workers is %d\n", cliArguments.Filename, cliArguments.Workers)
	fmt.Printf("Database connectivity check result: %t\n", connectivityCheck.SuccessfulConnection)
}