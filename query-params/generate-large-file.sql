-- A SQL script to create a large dataset using random, but valid data from the cpu_usage table

-- First, figure out the hosts and date ranges that are valid
with data_options as (
	select
		cu.host,
		min(cu.ts) as earliest_time,
		max(cu.ts) as latest_time,
		row_number() OVER () as data_option
	from cpu_usage cu
	group by cu.host
),
-- How many options are available total?
option_count as (
	select count(1) as value
	from data_options
),
-- Random selection of one of the available options in a big dataset
option_selection as (
	select floor(random() * oc.value + 1)::int as data_option
	from option_count oc
	cross join pg_catalog.generate_series(1, 10000)
),
-- Get a random, but valid option using each chosen host row
random_data as (
	select
		dop.host,
		dop.earliest_time + random() * (dop.latest_time - dop.earliest_time) time_a,
		dop.earliest_time + random() * (dop.latest_time - dop.earliest_time) time_b
	from data_options dop
	join option_selection os on dop.data_option = os.data_option
)
-- Ensure the start time is at or before the end time
select
	rd.host as hostname,
	least(rd.time_a, rd.time_b) start_time,
	greatest(rd.time_a, rd.time_b) end_time
from random_data rd;