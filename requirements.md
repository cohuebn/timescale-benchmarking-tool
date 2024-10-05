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
- [x] The tool should handle any invalid input appropriately

## Other criteria

- Robustness
  - The program handles a CSV file as long as it contains the expected headers (not just in the order)
    provided in the sample CSV. Example: `/run-against-local-database.sh --filename=../query-params/extra-headers.csv`
- Efficiency
  - The CSV file is processed line-by-line (not pulled entirely into memory)
  - Work is distributed as evenly as possible [using a hash-ring type algorithm](./benchmarking-tool/internal/benchmarking/worker_assigner.go#L9-L18)

## Optional functionality

- I added one additional output for the number of failed queries in case any queries fail during the run. I thought that was useful to ensure all queries ran successfully as part of benchmarking
- Unit tests can be run using instructions [here](./README.md#running-tests)
