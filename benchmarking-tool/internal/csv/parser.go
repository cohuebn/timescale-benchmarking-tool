package csv

import (
	"github.com/cohuebn/timescale-benchmarking-tool/internal/collections"
)

// A row of CPU usage data
type CpuUsageRow struct {
	Hostname string
	// Note that no time logic is needed in this application, so we can just use strings for the time fields.
	StartTime string
	EndTime string
}

func ParseCpuUsageRow(headers []string, row []string) CpuUsageRow {
	return CpuUsageRow{
		Hostname: row[collections.IndexOf(headers, "hostname")],
		StartTime: row[collections.IndexOf(headers, "start_time")],
		EndTime: row[collections.IndexOf(headers, "end_time")],
	}
}