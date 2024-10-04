package csv

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestReaderWithValidCsv(t *testing.T) {
	errGroup, ctx := errgroup.WithContext(context.Background())

	csvStream, _ := StreamCsvFile(ctx, "testdata/valid-csv.csv", errGroup)

	allStreamedResults := make([][]string, 0)
	for csvRow := range csvStream {
		allStreamedResults = append(allStreamedResults, csvRow.Row)
	}

	err := errGroup.Wait()
	
	assert.Nil(t, err)
	assert.Equal(t, len(allStreamedResults), 4)
	assert.Equal(t, allStreamedResults[0], []string{"hostname", "start_time", "end_time"})
	assert.Equal(t, allStreamedResults[1], []string{"host_000008", "2017-01-01 08:59:22", "2017-01-01 09:59:22"})
	assert.Equal(t, allStreamedResults[2], []string{"host_000001", "2017-01-02 13:02:02", "2017-01-02 14:02:02"})
	assert.Equal(t, allStreamedResults[3], []string{"host_000008", "2017-01-02 18:50:28", "2017-01-02 19:50:28"})
}

func TestReaderErrorWithMissingFileImmediatelyReturnsError(t *testing.T) {
	errGroup, ctx := errgroup.WithContext(context.Background())

	_, err := StreamCsvFile(ctx, "testdata/wah-wah-wee-wah-im-not-here.csv", errGroup)

	assert.Contains(t, err.Error(), "testdata/wah-wah-wee-wah-im-not-here.csv")
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestReaderErrorWithNonCsvReturnsErrorAsItsFound(t *testing.T) {
	errGroup, ctx := errgroup.WithContext(context.Background())

	StreamCsvFile(ctx, "testdata/not-a-csv.json", errGroup)

	err := errGroup.Wait()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "CSV parsing error for file testdata/not-a-csv.json")
	assert.Contains(t, err.Error(), "parse error on line 1, column 3")
}