package benchmarking

import "github.com/schollz/progressbar/v3"

// Get a progress bar for tracking the progress of query runs
func GetProgressBar() *progressbar.ProgressBar {
	return progressbar.NewOptions(
		-1,
		progressbar.OptionSetDescription("Queries processed"),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
	)
}