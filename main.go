package main

import (
	"github.com/JulienQNN/comai/cmd"
)

var (
	Version   = "dev"
	BuildDate = "unknown"
	GitBranch = "unknown"
	GitCommit = "unknown"
)

func main() {
	cmd.Execute()
}
