package benchmarking

import "time"

// The results of aggregating multiple CPU usage queries
type AggregatedCpuUsageResults struct {
	NumberOfQueriesProcessed int
	ErrorCount int
	TotalProcessingTime time.Duration
	MinimumQueryTime time.Duration
	MaximumQueryTime time.Duration
	MeanQueryTime time.Duration
	MedianQueryTime time.Duration
}