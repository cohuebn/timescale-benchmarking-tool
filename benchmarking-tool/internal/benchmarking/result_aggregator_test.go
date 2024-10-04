package benchmarking

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResultAggregatorWithNoMeasurements(test *testing.T) {
	aggregator := NewResultAggregator()

	result := aggregator.CalculateAggregates()

	zeroDuration := time.Duration(0)
	assert.Equal(test, 0, result.NumberOfQueriesProcessed)
	assert.Equal(test, 0, result.ErrorCount)
	assert.Equal(test, zeroDuration, result.MaximumQueryTime)
	assert.Equal(test, zeroDuration, result.MinimumQueryTime)
	assert.Equal(test, zeroDuration, result.MeanQueryTime)
	assert.Equal(test, zeroDuration, result.MedianQueryTime)
}

func TestResultAggregatorWithSingleSuccessMeasurement(test *testing.T) {
	aggregator := NewResultAggregator()

	queryTime := time.Duration(100)
	aggregator.AggregateCpuMeasure(QueryMeasurement{
		QueryTime: queryTime,
	})

	result := aggregator.CalculateAggregates()

	assert.Equal(test, 1, result.NumberOfQueriesProcessed)
	assert.Equal(test, 0, result.ErrorCount)
	assert.Equal(test, queryTime, result.MaximumQueryTime)
	assert.Equal(test, queryTime, result.MinimumQueryTime)
	assert.Equal(test, queryTime, result.MeanQueryTime)
	assert.Equal(test, queryTime, result.MedianQueryTime)
}

func TestResultAggregatorWithEvenNumberOfMeasurements(test *testing.T) {
	aggregator := NewResultAggregator()

	measurements := []QueryMeasurement{
		{QueryTime: time.Duration(100)},
		{QueryTime: time.Duration(200)},
		{QueryTime: time.Duration(300)},
		{QueryTime: time.Duration(500)},
	}
	for _, measurement := range measurements {
		aggregator.AggregateCpuMeasure(measurement)
	}

	result := aggregator.CalculateAggregates()

	assert.Equal(test, 4, result.NumberOfQueriesProcessed)
	assert.Equal(test, 0, result.ErrorCount)
	assert.Equal(test, time.Duration(500), result.MaximumQueryTime)
	assert.Equal(test, time.Duration(100), result.MinimumQueryTime)
	assert.Equal(test, time.Duration(275), result.MeanQueryTime)
	assert.Equal(test, time.Duration(250), result.MedianQueryTime)
}

func TestResultAggregatorWithOddNumberOfMeasurements(test *testing.T) {
	aggregator := NewResultAggregator()

	measurements := []QueryMeasurement{
		{QueryTime: time.Duration(100)},
		{QueryTime: time.Duration(200)},
		{QueryTime: time.Duration(500)},
	}
	for _, measurement := range measurements {
		aggregator.AggregateCpuMeasure(measurement)
	}

	result := aggregator.CalculateAggregates()

	assert.Equal(test, 3, result.NumberOfQueriesProcessed)
	assert.Equal(test, 0, result.ErrorCount)
	assert.Equal(test, time.Duration(500), result.MaximumQueryTime)
	assert.Equal(test, time.Duration(100), result.MinimumQueryTime)
	assert.Equal(test, time.Duration(266), result.MeanQueryTime)
	assert.Equal(test, time.Duration(200), result.MedianQueryTime)
}

func TestResultAggregatorRecordsSingleError(test *testing.T) {
	aggregator := NewResultAggregator()

	measurements := []QueryMeasurement{
		{QueryTime: time.Duration(100), Error: errors.New("Ouchies!")},
		{QueryTime: time.Duration(200)},
		{QueryTime: time.Duration(500)},
	}
	for _, measurement := range measurements {
		aggregator.AggregateCpuMeasure(measurement)
	}

	result := aggregator.CalculateAggregates()

	assert.Equal(test, 3, result.NumberOfQueriesProcessed)
	assert.Equal(test, 1, result.ErrorCount)
	assert.Equal(test, time.Duration(500), result.MaximumQueryTime)
	assert.Equal(test, time.Duration(100), result.MinimumQueryTime)
	assert.Equal(test, time.Duration(266), result.MeanQueryTime)
	assert.Equal(test, time.Duration(200), result.MedianQueryTime)
}

func TestResultAggregatorRecordsMultipleErrors(test *testing.T) {
	aggregator := NewResultAggregator()

	measurements := []QueryMeasurement{
		{QueryTime: time.Duration(100), Error: errors.New("Ouchies!")},
		{QueryTime: time.Duration(200)},
		{QueryTime: time.Duration(500), Error: errors.New("There's a snake in my boot!")},
	}
	for _, measurement := range measurements {
		aggregator.AggregateCpuMeasure(measurement)
	}

	result := aggregator.CalculateAggregates()

	assert.Equal(test, 3, result.NumberOfQueriesProcessed)
	assert.Equal(test, 2, result.ErrorCount)
	assert.Equal(test, time.Duration(500), result.MaximumQueryTime)
	assert.Equal(test, time.Duration(100), result.MinimumQueryTime)
	assert.Equal(test, time.Duration(266), result.MeanQueryTime)
	assert.Equal(test, time.Duration(200), result.MedianQueryTime)
}