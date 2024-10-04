package cli

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

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
	LogLevel slog.Level
}

// Given a user-input log-level as a string, return the corresponding slog.Level
func parseLogLevel(stringLevel string) slog.Level {
	switch strings.ToLower(stringLevel) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "info":
		return slog.LevelInfo
	default:
		// Can't use log here because we haven't set the log level yet
		fmt.Printf("Unknown log level %s, defaulting to info\n", stringLevel)
		return slog.LevelInfo
	}
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
	logLevel := flagSet.String("log-level", "info", "The log level.")
	ff.Parse(flagSet, os.Args[1:], ff.WithEnvVarPrefix("BENCHMARKING_TOOL"))
	return CliArguments{
		Filename: *filename,
		Workers: *workers,
		DatabaseHost: *dbHost,
		DatabasePort: *dbPort,
		DatabaseName: *dbName,
		DatabaseUsername: *dbUsername,
		DatabasePassword: *dbPassword,
		LogLevel: parseLogLevel(*logLevel),
	}
}