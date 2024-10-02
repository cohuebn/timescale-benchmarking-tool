package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateConnectionString(t *testing.T) {
	inputs := ConnectionStringInputs{
		Host: "localhost",
		Port: 5432,
		Username: "walter_sobachek",
		Password: "therearerules",
		Database: "homework",
	}

	result := CreateConnectionString(inputs)

	expected := "postgresql://walter_sobachek:therearerules@localhost:5432/homework"
	assert.Equal(t, expected, result)
}