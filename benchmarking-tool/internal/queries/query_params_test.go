package queries

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCpuUsageRow(t *testing.T) {
	headers := []string{"hostname", "start_time", "end_time"}
	row := []string{"host_000008", "2017-01-01 08:59:22", "2017-01-01 09:59:22"}

	result := ParseCpuUsageCsvRow(headers, row)

	assert.Equal(t, result.Hostname, "host_000008")
	assert.Equal(t, result.StartTime, "2017-01-01 08:59:22")
	assert.Equal(t, result.EndTime, "2017-01-01 09:59:22")
}

func TestParseCpuUsageRowVariedHeaderOrder(t *testing.T) {
	headers := []string{"end_time", "start_time", "hostname"}
	row := []string{"2017-01-01 09:59:22", "2017-01-01 08:59:22", "host_000008" }

	result := ParseCpuUsageCsvRow(headers, row)

	assert.Equal(t, result.Hostname, "host_000008")
	assert.Equal(t, result.StartTime, "2017-01-01 08:59:22")
	assert.Equal(t, result.EndTime, "2017-01-01 09:59:22")
}