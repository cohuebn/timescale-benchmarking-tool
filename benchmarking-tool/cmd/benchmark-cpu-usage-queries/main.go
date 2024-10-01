package main

import (
	"fmt"

	"github.com/cohuebn/timescale-benchmarking-tool/internal/cli"
)

func main() {
	cliArguments := cli.ParseCliArguments()
	fmt.Println("This is a stub for the benchmark-cpu-usage-queries command; it is not yet implemented.")
	fmt.Printf("The filename is: %s and the number of workers is %d\n", cliArguments.Filename, cliArguments.Workers)
}