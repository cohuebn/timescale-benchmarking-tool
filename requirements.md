# Requirements

This document is meant as a checklist for ensuring requirements are met.

## Functional requirements

- [x] Accepts a CSV file: File can be passed using the `--filename` parameter of the CLI
- [x] Accepts a parameter to determine how many workers are used: Number of workers can be passed using the `--workers` parameter of the CLI
- [x] The tool should handle any invalid input appropriately. Some of the failure-cases tested
  - [x] Path to non-existent file provided
  - [x] The file provided is not a CSV
  - [x] The file provided is missing one or more required columns
  - [x] Data in one or more of the data rows does not match expected format
  - [x] Wrong type for an input (e.g. a non-numeric input for the `--workers` parameter)
- [x] The tool should take the CSV row values `hostname`, `start_time`, and `end_time` and use them to
      generate a SQL query for each row. The query should return the max cpu usage and min cpu usage of the
      given hostname for every minute in the time range specified by the start time and end time. This is done [here](./benchmarking-tool/internal/database/queries.go#L5-L16)
- [x] Queries for the same hostname should be executed by the same worker each time. This is tested [here](benchmarking-tool/internal/benchmarking/worker_assigner_test.go#L10-L21)
- [x] The tool generates the expected outputs
  - [x] Number of queries processed
  - [x] Total processing time across all queries
  - [x] Minimum query time
  - [x] Maximum query time
  - [x] Median query time
  - [x] Average (mean) query time

## Other criteria

- Robustness
  - The program handles a CSV file as long as it contains the expected headers (not just in the order
    provided in the sample CSV). Example: `/run-against-local-database.sh --filename=../query-params/extra-headers.csv`
- Efficiency
  - The CSV file is processed line-by-line (not pulled entirely into memory)
  - Work is distributed [as evenly](benchmarking-tool/internal/benchmarking/worker_assigner_test.go#L23-L45) as possible [using a hash-ring type algorithm](./benchmarking-tool/internal/benchmarking/worker_assigner.go#L9-L18)
- Cleanliness
  - Code is organized into `cmd` (the CLI) and `internal` (all library functions used by the CLI). I tried finding a good line of separating the `internal` packages by function, but not going too crazy with overly-small packages.
  - While there may be ways to make the code more conscise by consolidating, I wanted to utilize the single-responsibility principle as best possible to prevent overly complex functions/packages.

## Optional functionality

- I added one additional output for the number of failed queries in case any queries fail during the run. I thought that was useful to ensure all queries ran successfully as part of benchmarking
- Unit tests can be run using instructions [here](./README.md#running-tests)
- A progress bar was added to give the user feedback at progress/speed of processing during
  the benchmark run.

## Potential improvements

Some TODO comments remain for things that would be nice-to-have given more time to work on the tool:

1. Find an efficient way to get the total lines in the CSV so that the progress bar has a finish line
   rather than being indeterminate
2. Come up with a more efficient approximation of median; To get exact median, I needed to store all query times until the end of processing and then get the median. However, there are ways to approximate the median without storing all values. E.g. storing a histogram and using it to approximate the median. However, for the sake of correctness and simplicity, I opted to calculate absolute median.
3. In a production application, I'd add integration tests to have automated validation that the Timescale interaction is working correctly.
4. If I had more time, I'd work on tuning the tool more for performance (e.g. channel buffer sizes). I was seeing about 95-105 queries per second in the tool's current state.
