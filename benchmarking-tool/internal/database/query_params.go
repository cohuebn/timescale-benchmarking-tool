package database

import (
	"github.com/cohuebn/timescale-benchmarking-tool/internal/collections"
)

// A row of CPU usage data
type CpuUsageQueryParams struct {
	Hostname string
	// Note that no time logic is needed in this application, so we can just use strings for the time fields.
	StartTime string
	EndTime string
}

// Parse a row containing CPU usage query params from a CSV file
func ParseCpuUsageCsvRow(headers []string, row []string) CpuUsageQueryParams {
	return CpuUsageQueryParams{
		Hostname: row[collections.IndexOf(headers, "hostname")],
		StartTime: row[collections.IndexOf(headers, "start_time")],
		EndTime: row[collections.IndexOf(headers, "end_time")],
	}
}