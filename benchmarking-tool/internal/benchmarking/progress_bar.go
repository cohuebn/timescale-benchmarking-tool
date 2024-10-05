package benchmarking

import "github.com/schollz/progressbar/v3"

// Get a progress bar for tracking the progress of query runs
func GetProgressBar() *progressbar.ProgressBar {
	// TODO: find an efficient way to get the total lines in the CSV so that the
	// progress bar has a finish line
	return progressbar.NewOptions(
		-1,
		progressbar.OptionSetDescription("Queries processed"),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
	)
}