# Requirements

This document is meant as a checklist for ensuring requirements are met.

## Functional requirements

- [x] Accepts a CSV file
  - File can be passed using the `--filename` parameter of the CLI
- [x] Accepts a parameter to determine how many workers are used
  - Number of workers can be passed using the `--workers` parameter of the CLI
- [x] CLI generates the expected outputs
  - [x] Number of queries processed
  - [x] Total processing time across all queries
  - [x] Minimum query time
  - [x] Maximum query time
  - [x] Median query time
  - [x] Average (mean) query time
  - I added one additional output for the number of failed queries in case any queries fail during the run.
