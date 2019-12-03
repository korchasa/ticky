package main

import (
	"github.com/jessevdk/go-flags"
	"os"
)

type FlagsSpec struct {
	Verbose bool `long:"verbose" short:"v" description:"Verbose output"`
	Statuses []string `long:"status" short:"s" default:"todo" default:"in-progress"`
	IssuesDir string `long:"issues-dir" short:"i" description:"Issues directory" default:"issues"`
	ReadmeCmd GenerateReadmeCommand `command:"readme" description:"Generate README.md"`
	NewCmd NewIssueCommand `command:"new" description:"Create new issue"`
	MyCmd NewIssueCommand `command:"my" description:"Show my issues"`
}

var Flags FlagsSpec

func main() {
	parser := flags.NewParser(&Flags, flags.Default)
	parser.LongDescription ="Tiny issue tracking in source code"
	parser.ShortDescription = "Tiny issue tracking in source code"
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}