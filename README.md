# timescale-benchmarking-tool

A tool to benchmark Timescale query performance using multiple workers to run multiple queries concurrently.

## Running locally via Docker

To run the benchmarking tool locally within Docker/Docker Compose, take the following steps:

1. Ensure you have [Docker/Docker Compose installed](https://docs.docker.com/compose/install/).
2. Run the database using Docker Compose. This command will use the local.env file in this repository and ensure the latest changes are included in running containers: `docker compose --env-file local.env up --build`
3. Launch the benchmarking tool into the Docker Compose network created in the previous step. This command will take care of building the image, splice in environment variables, and run the benchmarking tool in the database's network: `benchmarking-tool/run-locally`
