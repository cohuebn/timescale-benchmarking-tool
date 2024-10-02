#!/bin/sh

# This script is used to build and launch the Benchmarking Tool into the local Docker Compose environment.
# It will run against the local TimescaleDB instance from the Docker Compose stack.

set -e

script_directory="$( cd "$( dirname "$0" )" && pwd )"

# Build the image
docker build -t benchmarking-tool $script_directory
# Find the network name of the Docker Compose stack
docker_network=$(docker network ls --filter name=timescale-benchmarking --format "{{.Name}}")
# Run the container
docker run --env-file ../local.env --network $docker_network benchmarking-tool