# timescale-benchmarking-tool

A tool to benchmark Timescale query performance using multiple workers to run multiple queries concurrently.

Requirements for this tool are documented here, along with some rationale around decisions made: [requirements.md](./requirements.md)

Some useful diagrams detailing system architecture can be found here: [architecture.md](./architecture.md)

## Dependencies

To run this tool locally, you'll need the following:

1. [Docker](https://docs.docker.com/compose/install/).

If you want to run entirely in Docker, none of the subsequent dependencies are necessary. If you want to
run the benchmarking tool Go application outside Docker, the following will either be needed or helpful:

1. [Go](https://go.dev/doc/install) (Required) - Go is needed to compile, test, etc.
2. [Make](https://www.gnu.org/software/make/) (Optional) - Make is not required, but the included [Makefile](./benchmarking-tool/Makefile) can simplify installing, building, etc. See the targets in that file for `install`, `build`, etc.

## Running locally

This section details how to run the tool locally.

### Docker

The simplest way to test this tool is to run it within Docker/Docker Compose. To do so, take the following steps:

1. Run the database using Docker Compose. This command will use the local.env file in this repository and ensure the latest changes are included in running containers: `docker compose --env-file local.env up --build`
2. Launch the benchmarking tool into the Docker Compose network created in the previous step: `docker compose --env-file local.env run --build benchmarking-tool --filename=/query-params/query-params.csv`.

### Running the benchmarking tool outside of Docker

You might want to run the tool outside of Docker (faster debugging, connecting to databases outside Docker, etc.). In case you want to run the benchmarking tool outside of Docker, take the following steps:

1. You still need to run the database in Docker Compose. Run steps 1 from the [Running locally via Docker](#running-locally-via-docker) section above
2. Go to the [benchmarking-tool](./benchmarking-tool/) directory: `cd benchmarking-tool/`
3. Ensure you have all dependencies installed. Run: `go mod download`
4. Run the benchmarking tool against the Docker database. A helper script has been created to automatically wire up to that database. Run: `./run-against-local-database.sh`.

Note: Parameters can be passed to that [run-against-local-database.sh](./benchmarking-tool/run-against-local-database.sh)
script to change the CSV file or number of workers. Example: `./run-against-local-database.sh --filename=../query-params/reordered-headers.csv --workers=16`

To see all options for the CLI, run `./run-against-local-database.sh --help` or `go run cmd/benchmarking-tool/main.go --help`.

### Sample CSV files

By default, this tool runs against [the sample CSV file](./query-params/query-params.csv) provided as part of the assignment. However, you can test against other files in that directory using the `--filename` parameter of the CLI.

When running in Docker Compose, the benchmarking tool has access to the [query-params](./query-params/) directory mounted at the path `/query-params` on the container.

This directory contains a the following files for testing the tool:

- [query_params.csv](./query-params/query-params.csv): This is the file provided as part of the assignment by Timescale
- [reordered-headers.csv](./query-params/reordered-headers.csv): This contains the same data as [query_params.csv](./query-params/query-params.csv), but the `hostname` header has been moved to the end of the file to test handling any order of the three required columns
- [not-a-csv.json](./query-params/not-a-csv.json): This is a JSON file that can be used to see the error
  returned when a non-CSV file is provided to the benchmarking tool
- [generate-large-file.sql](./query-params/generate-large-file.sql): This is a SQL query that can be run against the database to get a large dataset for testing

## Running tests

If you'd like to run unit tests, take the following steps:

1. Go to the [benchmarking-tool](./benchmarking-tool/) directory: `cd benchmarking-tool/`
2. Ensure you've downloaded all dependencies locally: `go mod download`
3. Run all tests: `go test ./...`
