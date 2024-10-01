package cli

import (
	"flag"
)

// Command line arguments for the benchmarking tool.
type CliArguments struct {
	Filename string
	Workers int
}

// Parse the command line arguments from user input
func ParseCliArguments() CliArguments {
	filename := flag.String("filename", "cpu-usage-queries.txt", "The name of the file providing query inputs.")
	workers := flag.Int("workers", 1, "The number of workers to use for running concurrent queries.")
	flag.Parse()
	return CliArguments{
		Filename: *filename,
		Workers: *workers,
	}
}