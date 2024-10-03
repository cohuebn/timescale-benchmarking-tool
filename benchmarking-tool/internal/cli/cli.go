package cli

import (
	"flag"
	"os"

	"github.com/peterbourgon/ff/v3"
)

// Command line arguments for the benchmarking tool.
type CliArguments struct {
	Filename string
	Workers int
	DatabaseHost string
	DatabasePort int
	DatabaseName string
	DatabaseUsername string
	DatabasePassword string
}

// Parse the command line arguments from user input
func ParseCliArguments() CliArguments {
	flagSet := flag.NewFlagSet("benchmarking-tool", flag.ExitOnError)
	filename := flagSet.String("filename", "../query-params/query-params.csv", "The name of the file providing query inputs.")
	workers := flagSet.Int("workers", 1, "The number of workers to use for running concurrent queries.")
	dbHost := flagSet.String("database-host", "localhost", "The database host.")
	dbPort := flagSet.Int("database-port", 5432, "The database port.")
	dbName := flagSet.String("database-name", "postgres", "The name of the database.")
	dbUsername := flagSet.String("database-username", "postgres", "The username for the database connection.")
	dbPassword := flagSet.String("database-password", "", "The password for the database connection.")
	ff.Parse(flagSet, os.Args[1:], ff.WithEnvVarPrefix("BENCHMARKING_TOOL"))
	return CliArguments{
		Filename: *filename,
		Workers: *workers,
		DatabaseHost: *dbHost,
		DatabasePort: *dbPort,
		DatabaseName: *dbName,
		DatabaseUsername: *dbUsername,
		DatabasePassword: *dbPassword,
	}
}