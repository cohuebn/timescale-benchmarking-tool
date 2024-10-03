#! /bin/sh

# A script to make it easier to run the benchmarking tool against the local Docker database.

# Ensure this works regardless of current working directory
script_path="${BASH_SOURCE:-$0}"
absolute_script_path="$(realpath "${script_path}")"
script_dir="$(dirname "${absolute_script_path}")"


source "${script_dir}/../local.env";
go run cmd/benchmarking-tool/main.go \
		--database-host=${DB_HOST} \
		--database-port=${DB_PORT} \
		--database-name=${DB_NAME} \
		--database-username=${BENCHMARKING_USER} \
		--database-password=${BENCHMARKING_PASSWORD} \
		"$@";