package cli

import (
	"log/slog"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/database"
)

// Create a connection string using the provided CLI arguments
func CreateConnectionString(cliArguments CliArguments) string {
	// Intentionally don't log out the password
	slog.Debug(
		"Creating connection string from CLI arguments",
		"host", cliArguments.DatabaseHost,
		"port", cliArguments.DatabasePort,
		"username", cliArguments.DatabaseUsername,
		"database", cliArguments.DatabaseName,
	)
	return database.CreateConnectionString(database.ConnectionStringInputs{
		Host: cliArguments.DatabaseHost,
		Port: cliArguments.DatabasePort,
		Username: cliArguments.DatabaseUsername,
		Password: cliArguments.DatabasePassword,
		Database: cliArguments.DatabaseName,
	})
}