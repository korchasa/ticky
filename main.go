package main

import (
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

var Flags struct {
	Verbose     bool `long:"verbose" short:"v" description:"Verbose output"`
	Statuses    []string
	StatusesStr string                `long:"statuses" short:"s" default:"todo,in-progress"`
	IssuesDir   string                `long:"issues-dir" short:"i" description:"Issues directory" default:"issues"`
	ReadmeCmd   GenerateReadmeCommand `command:"readme" description:"Generate README.md"`
	NewCmd      NewIssueCommand       `command:"new" description:"Create new issue"`
	MyCmd       MyIssuesCommand       `command:"my" description:"Show my issues"`
	ListCmd     ListIssuesCommand     `command:"list" description:"List issues"`
	TagCmd     NewTagCommand     `command:"tag" description:"List issues"`
}

func main() {
	log.SetLevel(log.WarnLevel)
	parser := flags.NewParser(&Flags, flags.Default)
	parser.ShortDescription = "Tiny issue tracking in source code"
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		if Flags.Verbose {
			log.SetLevel(log.InfoLevel)
		}
		for _, st := range strings.Split(Flags.StatusesStr, ",") {
			Flags.Statuses = append(Flags.Statuses, strings.Trim(st, " "))
		}
		return command.Execute(args)
	}
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
