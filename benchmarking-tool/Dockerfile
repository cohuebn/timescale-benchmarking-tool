FROM golang:1.23.2-alpine3.20 AS builder

WORKDIR /benchmarking-tool

# Download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code and build the cpu usage queries benchmarking tool
COPY . .
RUN go build -o bin/benchmarking-tool cmd/benchmarking-tool/main.go

FROM golang:1.23.2-alpine3.20

WORKDIR /benchmarking-tool
COPY --from=builder /benchmarking-tool/bin/benchmarking-tool .

ENTRYPOINT [ "./benchmarking-tool" ]