package csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderWithValidCsv(t *testing.T) {
	errorChannel := make(chan error)

	csvStream, _ := StreamCsvFile("testdata/valid-csv.csv", errorChannel)

	allStreamedResults := make([][]string, 0)
	StreamingLoop:
		for {
			select {
			case csvRow, ok := <-csvStream:
				if !ok {
					break StreamingLoop
				}
				allStreamedResults = append(allStreamedResults, csvRow.Row)
			case err := <-errorChannel:
				assert.Fail(t, "Error while streaming CSV file", err)
			}
		}

	assert.Equal(t, len(allStreamedResults), 4)
	assert.Equal(t, allStreamedResults[0], []string{"hostname", "start_time", "end_time"})
	assert.Equal(t, allStreamedResults[1], []string{"host_000008", "2017-01-01 08:59:22", "2017-01-01 09:59:22"})
	assert.Equal(t, allStreamedResults[2], []string{"host_000001", "2017-01-02 13:02:02", "2017-01-02 14:02:02"})
	assert.Equal(t, allStreamedResults[3], []string{"host_000008", "2017-01-02 18:50:28", "2017-01-02 19:50:28"})
}

func TestReaderErrorWithMissingFile(t *testing.T) {
	errorChannel := make(chan error)

	_, err := StreamCsvFile("testdata/wah-wah-wee-wah-im-not-here.csv", errorChannel)

	assert.Contains(t, err.Error(), "testdata/wah-wah-wee-wah-im-not-here.csv")
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestReaderErrorWithNonCsv(t *testing.T) {
	errorChannel := make(chan error)

	csvStream, _ := StreamCsvFile("testdata/not-a-csv.json", errorChannel)

	select {
	case <-csvStream:
		assert.Fail(t, "Expected error while streaming CSV file")
	case err := <-errorChannel:
		assert.Contains(t, err.Error(), "CSV parsing error for file testdata/not-a-csv.json")
		assert.Contains(t, err.Error(), "parse error on line 1, column 3")
	}
}