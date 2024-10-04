package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertCpuUsageParamsToQuery(t *testing.T) {
	params := CpuUsageQueryParams{
		Hostname:  "host_12345",
		StartTime: "2017-01-02 03:04:05",
		EndTime:   "2017-02-03 06:59:59",
	}

	result := ConvertCpuUsageParamsToQuery(params)

	// Not testing everything here, just the parameter replacement
	// In a production application, I'd write integration tests to validate the query returns expected results
	// when run against an actual TimescaleDB instance
	assert.Contains(t, result, "where host = 'host_12345'")
	assert.Contains(t, result, "and ts between '2017-01-02 03:04:05' and '2017-02-03 06:59:59'")
}