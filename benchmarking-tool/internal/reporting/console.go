package reporting

import (
	"os"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/benchmarking"
	"github.com/jedib0t/go-pretty/v6/table"
)

func LogCpuUsageResultsToConsole(aggregatedResults benchmarking.AggregatedCpuUsageResults) {
	tableWriter := table.NewWriter()
	tableWriter.SetOutputMirror(os.Stdout)
	tableWriter.AppendHeader(table.Row{"Metric", "Value"})
	tableWriter.AppendRow(table.Row{"Total number of queries processed", aggregatedResults.NumberOfQueriesProcessed})
	tableWriter.AppendRow(table.Row{"Total number of failed queries", aggregatedResults.ErrorCount})
	tableWriter.AppendRow(table.Row{"Total processing time", aggregatedResults.TotalProcessingTime})
	tableWriter.AppendRow(table.Row{"Maximum query time", aggregatedResults.MaximumQueryTime})
	tableWriter.AppendRow(table.Row{"Minimum query time", aggregatedResults.MinimumQueryTime})
	tableWriter.AppendRow(table.Row{"Mean query time", aggregatedResults.MeanQueryTime})
	tableWriter.AppendRow(table.Row{"Median query time", aggregatedResults.MedianQueryTime})
	tableWriter.SetStyle(table.StyleRounded)
	tableWriter.Render()
}