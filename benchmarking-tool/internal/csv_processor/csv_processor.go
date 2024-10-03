package csv_processor

import (
	"log"
	"slices"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/csv"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/queries"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/workers"
	"github.com/jackc/pgx/v5/pgxpool"
)

func validateHeaderRow(row []string) {
	requiredColumns := []string{"hostname", "start_time", "end_time"}
	missingColumns := []string{}
	for _, column := range requiredColumns {
		if !slices.Contains(row, column) {
			missingColumns = append(missingColumns, column)
		}
	}
	if len(missingColumns) > 0 {
		log.Panicf("Missing required columns in header row: %v", missingColumns)
	}
}

// Convert a stream of CSV rows into a stream of query parameters
func getQueryParamsStream(csvStream <-chan csv.CsvStreamingResult) <-chan queries.CpuUsageQueryParams {
	queryParamsStream := make(chan queries.CpuUsageQueryParams, 100)
	go func() {
		defer close(queryParamsStream)
		var headerRow []string
		for csvRow := range csvStream {
			if csvRow.Error != nil {
				log.Panic(csvRow.Error)
			}
			
			if headerRow == nil {
				validateHeaderRow(csvRow.Row)
				headerRow = csvRow.Row
			} else {
				queryParamsStream <- queries.ParseCpuUsageCsvRow(headerRow, csvRow.Row)
			}
		}
	}()

	return queryParamsStream
}

// Stream CSV rows through worker pools, run queries using those workers, and return aggregate results
func ProcessCsv(numberOfWorkers int, connectionPool *pgxpool.Pool, csvStream <-chan csv.CsvStreamingResult) workers.AggregatedCpuUsageResults {
	queryParamsStream := getQueryParamsStream(csvStream)
	return workers.MeasureCpuUsageQueries(numberOfWorkers, connectionPool, queryParamsStream)
}