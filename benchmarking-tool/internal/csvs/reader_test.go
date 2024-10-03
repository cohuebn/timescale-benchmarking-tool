package csvs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderWithValidCsv(t *testing.T) {
	streamingResults := StreamCsvFile("testdata/valid-csv.csv")

	allStreamedResults := make([][]string, 0)
	for csvRow := range streamingResults {
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