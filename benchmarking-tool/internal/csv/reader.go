package csv

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"golang.org/x/sync/errgroup"
)

// An error that occurs when parsing a CSV file
type CsvParsingError struct {
	Filename string
	Err   error
}

func (err *CsvParsingError) Error() string {
	return fmt.Sprintf("CSV parsing error for file %s. Caused by: %s", err.Filename, err.Err)
}

// A single result while streaming a CSV file
type CsvStreamingResult struct {
	Row []string
}

// Read a CSV file incrementally and stream its contents to the provided channel
func StreamCsvFile(ctx context.Context, filename string, errGroup *errgroup.Group) (<-chan CsvStreamingResult, error) {
	outputChannel := make(chan CsvStreamingResult, 100)

	file, err := os.Open(filename)
	// If the file can't be opened, stop processing
	if err != nil {
		return nil, &CsvParsingError{
			Filename: filename,
			Err: 		err,
		}
	}

	errGroup.Go(func() error {
		csvReader := csv.NewReader(file)
		defer func() {
			file.Close()
			close(outputChannel)
			slog.Debug("Finished cleaning up CSV reader", "filename", filename)
		}()

		for {
			select {
			case <-ctx.Done():
				slog.Debug("Stopping CSV file streaming due to an external error")
				return ctx.Err()
			default:
				record, err := csvReader.Read()
				// If we reach the end of the file, stop processing
				if errors.Is(err, io.EOF) {
					return nil
				}
				// If there's an unexpected error while reading a line, send it to the error group
				if err != nil {
					return &CsvParsingError{
						Filename: filename,
						Err:      err,
					}
				}
				// If everything is fine, send the row to the output channel
				outputChannel <- CsvStreamingResult{
					Row: record,
				}
			}
		}
	})

	return outputChannel, nil
}