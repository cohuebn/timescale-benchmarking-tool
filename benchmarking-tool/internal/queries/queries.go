package queries

import "fmt"

func ConvertCpuUsageParamsToQuery(params CpuUsageQueryParams) string {
	return fmt.Sprintf(`
		select
			time_bucket('1 minute', ts) as bucket,
			max(usage) as max_usage,
			min(usage) as min_usage
		from cpu_usage
		where host = '%s'
		and ts between '%s' and '%s'
		group by bucket;
	`, params.Hostname, params.StartTime, params.EndTime)
}