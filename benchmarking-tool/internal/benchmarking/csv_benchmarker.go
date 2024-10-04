package benchmarking

import (
	"log"
	"slices"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/csv"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/database"
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
func getQueryParamsStream(csvStream <-chan csv.CsvStreamingResult) <-chan database.CpuUsageQueryParams {
	queryParamsStream := make(chan database.CpuUsageQueryParams, 100)
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
				queryParamsStream <- database.ParseCpuUsageCsvRow(headerRow, csvRow.Row)
			}
		}
	}()

	return queryParamsStream
}

// Stream CSV rows through worker pools, run queries using those workers, and return aggregate results
func ProcessCsv(numberOfWorkers int, connectionPool *pgxpool.Pool, csvStream <-chan csv.CsvStreamingResult) AggregatedCpuUsageResults {
	queryParamsStream := getQueryParamsStream(csvStream)
	return RunCpuUsageQueries(numberOfWorkers, connectionPool, queryParamsStream)
}