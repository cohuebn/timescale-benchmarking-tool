# timescale-benchmarking-tool

A tool to benchmark Timescale query performance using multiple workers to run multiple queries concurrently.

## Running locally

To run the benchmarking tool locally within Docker Compose, take the following steps:

1. Ensure you have [Docker/Docker Compose installed](https://docs.docker.com/compose/install/).
2. Run the stack using Docker Compose. This command will use the local.env file in this repository and ensure the latest changes are included in running containers: `docker compose --env-file local.env up --build`
