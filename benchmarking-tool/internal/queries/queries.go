package queries

import "fmt"

func ConvertCpuUsageParamsToQuery(params CpuUsageQueryParams) string {
	return fmt.Sprintf(`
		select
			hostname,
			time_bucket('1 minute', time) as time,
			max(usage) as max_usage,
			min(usage) as min_usage,
		from cpu_usage
		where hostname = '%s'
		and time between '%s' and '%s';
	`, params.Hostname, params.StartTime, params.EndTime)
}