# timescale-benchmarking-tool

A tool to benchmark Timescale query performance using multiple workers to run multiple queries concurrently.

## Running locally via Docker

To run the benchmarking tool locally within Docker/Docker Compose, take the following steps:

1. Ensure you have [Docker/Docker Compose installed](https://docs.docker.com/compose/install/).
2. Run the database using Docker Compose. This command will use the local.env file in this repository and ensure the latest changes are included in running containers: `docker compose --env-file local.env up --build`
3. Launch the benchmarking tool into the Docker Compose network created in the previous step: `docker compose --env-file local.env run benchmarking-tool --filename=/query-params/query-params.csv`.

### Sample CSV files

When running in Docker Compose, the benchmarking tool has access to the [query-params](./query-params/) directory. This directory contains sample CSV files containing query params for the benmarking tool to run.
The directory is mounted at the path `/query-params` on the container.

This directory contains a the following files for testing the tool:

- [query_params.csv](./query-params/query-params.csv): This is the file provided as part of the assignment by Timescale
- [reordered-headers.csv](./query-params/reordered-headers.csv): This contains the same data as [query_params.csv](./query-params/query-params.csv), but the `hostname` header has been moved to the end of the file to test handling any order of the three required columns
