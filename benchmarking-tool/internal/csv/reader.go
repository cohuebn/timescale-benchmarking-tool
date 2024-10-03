package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

// An error that occurs when parsing a CSV file
type CsvParsingError struct {
	Filename string
	Err   error
}

func (err *CsvParsingError) Error() string {
	return fmt.Sprintf("CSV parsing error for file %s: %s", err.Filename, err.Err)
}

// A single result while streaming a CSV file
type CsvStreamingResult struct {
	Row []string
	Error error
}

// Read a CSV file incrementally and stream its contents to the provided channel
func StreamCsvFile(filename string) <-chan CsvStreamingResult {
	outputChannel := make(chan CsvStreamingResult, 10)

	file, err := os.Open(filename)
	// If the file can't be opened, stream the error and stop processing
	if err != nil {
		outputChannel <- CsvStreamingResult{
			Row: nil,
			Error: &CsvParsingError{
				Filename: filename,
				Err:      err,
			},
		}
		return outputChannel
	}

	go func() {
		csvReader := csv.NewReader(file)
		defer file.Close()
		defer close(outputChannel)
		for {
			record, err := csvReader.Read()
			// If we reach the end of the file, stop processing
			if errors.Is(err, io.EOF) {
				break
			}
			// If there's an error while reading a line, stream the error and stop processing
			if err != nil {
				outputChannel <- CsvStreamingResult{
					Row: nil,
					Error: &CsvParsingError{
						Filename: filename,
						Err:      err,
					},
				}
				break
			}
			// Stream the current row to the output channel
			outputChannel <- CsvStreamingResult{
				Row: record,
				Error: nil,
			}
		}
	}()

	return outputChannel
}