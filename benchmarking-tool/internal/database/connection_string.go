package database

import "fmt"

type ConnectionStringInputs struct {
	Host string
	Port int
	Username string
	Password string
	Database string
}

// Create a connection string for the database using the provided inputs.
func CreateConnectionString(inputs ConnectionStringInputs) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", inputs.Username, inputs.Password, inputs.Host, inputs.Port, inputs.Database)
}