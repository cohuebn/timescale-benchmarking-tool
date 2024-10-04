package benchmarking

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/csv"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

func validateHeaderRow(row []string) error {
	requiredColumns := []string{"hostname", "start_time", "end_time"}
	missingColumns := []string{}
	for _, column := range requiredColumns {
		if !slices.Contains(row, column) {
			missingColumns = append(missingColumns, column)
		}
	}
	if len(missingColumns) > 0 {
		slog.Debug("Detected missing columns in header row", "missingColumns", missingColumns)
		return fmt.Errorf("missing required columns in header row: %v", strings.Join(missingColumns, ", "))
	}
	return nil
}

// If processing fails, read the CSV stream and log a count of rows that were read from the file, but not processed
// fully
func cleanoutCsvStream(csvStream <-chan csv.CsvStreamingResult) {
	skippedRowCount := 0
	for range csvStream {
		skippedRowCount++
	}
	if (skippedRowCount > 0) {
		slog.Debug("CSV processing failed; some rows were skipped", "skippedRows", skippedRowCount)
	}
}

// Convert a stream of CSV rows into a stream of query parameters
func getQueryParamsStream(ctx context.Context, csvStream <-chan csv.CsvStreamingResult, errGroup *errgroup.Group) <-chan database.CpuUsageQueryParams {
	queryParamsStream := make(chan database.CpuUsageQueryParams, 100)
	errGroup.Go(func() error {
		defer func() {
			cleanoutCsvStream(csvStream)
			close(queryParamsStream)
			slog.Debug("Finished cleaning up query param streaming")
		}()

		var headerRow []string
		for {
			select {
			case <-ctx.Done():
				slog.Debug("Stopping query param streaming due to an external error")
				return ctx.Err()
			case csvRow, ok := <-csvStream:
				if !ok {
					return nil
				}
				if headerRow == nil {
					// Validate and hang on to the header row so it can be used to label subsequent rows
					headerRow = csvRow.Row
					headerRowErr := validateHeaderRow(csvRow.Row)
					if (headerRowErr != nil) {
						return headerRowErr
					}
				} else {
					queryParamsStream <- database.ParseCpuUsageCsvRow(headerRow, csvRow.Row)
				}
			}
		}
	})

	return queryParamsStream
}

// Stream CSV rows through worker pools, run queries using those workers, and return aggregate results
func ProcessCsv(ctx context.Context, numberOfWorkers int, connectionPool *pgxpool.Pool, csvStream <-chan csv.CsvStreamingResult, errGroup *errgroup.Group) AggregatedCpuUsageResults {
	queryParamsStream := getQueryParamsStream(ctx, csvStream, errGroup)
	return RunCpuUsageQueries(ctx, numberOfWorkers, connectionPool, queryParamsStream, errGroup)
}