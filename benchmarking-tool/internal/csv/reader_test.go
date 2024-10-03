package csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderWithValidCsv(t *testing.T) {
	csvStream := StreamCsvFile("testdata/valid-csv.csv")

	allStreamedResults := make([][]string, 0)
	for csvRow := range csvStream {
		if csvRow.Error != nil {
			panic(csvRow.Error)
		}
		allStreamedResults = append(allStreamedResults, csvRow.Row)
	}

	assert.Equal(t, len(allStreamedResults), 4)
	assert.Equal(t, allStreamedResults[0], []string{"hostname", "start_time", "end_time"})
	assert.Equal(t, allStreamedResults[1], []string{"host_000008", "2017-01-01 08:59:22", "2017-01-01 09:59:22"})
	assert.Equal(t, allStreamedResults[2], []string{"host_000001", "2017-01-02 13:02:02", "2017-01-02 14:02:02"})
	assert.Equal(t, allStreamedResults[3], []string{"host_000008", "2017-01-02 18:50:28", "2017-01-02 19:50:28"})
}

func TestReaderErrorWithMissingFile(t *testing.T) {
	csvStream := StreamCsvFile("testdata/wah-wah-wee-wah-im-not-here.csv")

	result := <-csvStream

	assert.Nil(t, result.Row)
	assert.NotNil(t, result.Error)
	err := result.Error.Error()
	assert.Contains(t, err, "CSV parsing error for file testdata/wah-wah-wee-wah-im-not-here.csv")
	assert.Contains(t, err, "no such file or directory")
}

func TestReaderErrorWithNonCsv(t *testing.T) {
	csvStream := StreamCsvFile("testdata/not-a-csv.json")

	result := <-csvStream

	assert.Nil(t, result.Row)
	assert.NotNil(t, result.Error)
	err := result.Error.Error()
	assert.Contains(t, err, "CSV parsing error for file testdata/not-a-csv.json")
	assert.Contains(t, err, "parse error on line 1, column 3")
}