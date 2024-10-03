package workers

import (
	"sort"
	"time"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/queries"
	"github.com/cohuebn/timescale-benchmarking-tool/internal/results"
)

// A structure that can be used to aggregate the results of multiple CPU usage queries
type ResultAggregator struct {
	numberOfQueriesProcessed int
	errorCount int
	totalProcessingTime time.Duration
	minimumQueryTime time.Duration
	maximumQueryTime time.Duration
	// Use of "mean" instead of "average" just to reduce ambiguity with median
	meanQueryTime time.Duration
	
	// TODO: As best I can tell, median requires storing all query times in memory
	// Ideally, we could use an approximation here to avoid doing that. However,
	// in the interest of simplicity, I'm just going to store all query times in memory
	// All other measures (total, mean, min, max, etc.) can be calculated in a streaming fashion without storing individual query times
	queryTimes []time.Duration
}

// Initialize a new ResultAggregator with an empty slice of query times
// All other aggregator fields can just use their default values
func NewResultAggregator() ResultAggregator {
	return ResultAggregator{
		queryTimes: make([]time.Duration, 0),
	}
}

// Get the amount to increment the error count by based on the measurement
func getErrorIncrement(measurement queries.QueryMeasurement) int {
	if (measurement.Error != nil) {
		return 1
	}
	return 0
}

func (aggregator *ResultAggregator) AggregateCpuMeasure(measurement queries.QueryMeasurement) {
	// If this is the first measurement, set all values using just the measurement
	// Otherwise, update the values based on the new measurement
	if (aggregator.numberOfQueriesProcessed == 0) {
		aggregator.numberOfQueriesProcessed = 1
		aggregator.totalProcessingTime = measurement.QueryTime
		aggregator.errorCount = getErrorIncrement(measurement)
		aggregator.maximumQueryTime = measurement.QueryTime
		aggregator.minimumQueryTime = measurement.QueryTime
		aggregator.meanQueryTime = measurement.QueryTime
		aggregator.queryTimes = []time.Duration{measurement.QueryTime}
	} else {
		aggregator.numberOfQueriesProcessed++
		aggregator.totalProcessingTime += measurement.QueryTime
		aggregator.errorCount += getErrorIncrement(measurement)
		if (measurement.QueryTime < aggregator.minimumQueryTime) {
			aggregator.minimumQueryTime = measurement.QueryTime
		}
		if (measurement.QueryTime > aggregator.maximumQueryTime) {
			aggregator.maximumQueryTime = measurement.QueryTime
		}
		aggregator.meanQueryTime = aggregator.totalProcessingTime / time.Duration(aggregator.numberOfQueriesProcessed)
		aggregator.queryTimes = append(aggregator.queryTimes, measurement.QueryTime)
	}
}

// Get the median query time from a list of all query times
func getMedian(queryTimes []time.Duration) time.Duration {
	// If there are no durations, return 0
	if (len(queryTimes) == 0) {
		return 0
	}
	if (len(queryTimes) == 1) {
		return queryTimes[0]
	}
	
	// Get a sorted copy of the query times
	sortedQueryTimes := make([]time.Duration, len(queryTimes))
	copy(sortedQueryTimes, queryTimes)
	sort.Slice(sortedQueryTimes, func(i, j int) bool {
		return queryTimes[i] < queryTimes[j]
	})
	// If there are an odd number of durations, return the middle duration
	middleDurationIndex := len(sortedQueryTimes) / 2
	if (middleDurationIndex % 2 == 1) {
		return sortedQueryTimes[middleDurationIndex]
	}
	// If there are an even number of durations, return the average of the two middle durations
	return (sortedQueryTimes[middleDurationIndex - 1] + sortedQueryTimes[middleDurationIndex]) / 2.0
}

func (aggregator *ResultAggregator) CalculateAggregates() results.AggregatedCpuUsageResults {
	// Calculate the median query time; all other values can be calculated in a streaming fashion
	medianQueryTime := getMedian(aggregator.queryTimes)
	return results.AggregatedCpuUsageResults{
		NumberOfQueriesProcessed: aggregator.numberOfQueriesProcessed,
		ErrorCount: aggregator.errorCount,
		TotalProcessingTime: aggregator.totalProcessingTime,
		MinimumQueryTime: aggregator.minimumQueryTime,
		MaximumQueryTime: aggregator.maximumQueryTime,
		MeanQueryTime: aggregator.meanQueryTime,
		MedianQueryTime: medianQueryTime,
	}
}