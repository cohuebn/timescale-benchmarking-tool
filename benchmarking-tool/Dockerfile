FROM golang:1.23 AS builder

WORKDIR /benchmarking-tool

# Download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code and build the cpu usage queries benchmarking tool
COPY . .
RUN go build -o bin/benchmark-cpu-usage-queries cmd/benchmark-cpu-usage-queries/main.go

FROM scratch

WORKDIR /benchmarking-tool
COPY --from=builder /benchmarking-tool/bin/benchmark-cpu-usage-queries .

ENTRYPOINT [ "./benchmark-cpu-usage-queries" ]