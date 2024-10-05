# Requirements

This document is meant as a checklist for ensuring requirements are met.

## Functional requirements

- [x] Accepts a CSV file: File can be passed using the `--filename` parameter of the CLI
- [x] Accepts a parameter to determine how many workers are used: Number of workers can be passed using the `--workers` parameter of the CLI
- [x] CLI generates the expected outputs
  - [x] Number of queries processed
  - [x] Total processing time across all queries
  - [x] Minimum query time
  - [x] Maximum query time
  - [x] Median query time
  - [x] Average (mean) query time
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

## Other criteria

- Robustness
  - The program handles a CSV file as long as it contains the expected headers (not just in the order)
    provided in the sample CSV. Example: `/run-against-local-database.sh --filename=../query-params/extra-headers.csv`
- Efficiency
  - The CSV file is processed line-by-line (not pulled entirely into memory)
  - Work is distributed [as evenly](benchmarking-tool/internal/benchmarking/worker_assigner_test.go#L23-L45) as possible [using a hash-ring type algorithm](./benchmarking-tool/internal/benchmarking/worker_assigner.go#L9-L18)

## Optional functionality

- I added one additional output for the number of failed queries in case any queries fail during the run. I thought that was useful to ensure all queries ran successfully as part of benchmarking
- Unit tests can be run using instructions [here](./README.md#running-tests)
