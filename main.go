package main

import (
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Opts struct {
	Verbose     bool `long:"verbose" short:"v" description:"Verbose output"`
	Statuses    []string
	StatusesStr string `long:"statuses" default:"todo,in-progress"`
	IssuesDir   string `long:"issues-dir" description:"Issues directory" default:"issues"`
}

func main() {
	log.SetLevel(log.WarnLevel)
	var opts struct {
		*Opts
		InitCmd    InitCommand       `command:"init" description:"Init ticky"`
		NewCmd    NewIssueCommand       `command:"new" description:"Create new issue"`
		MyCmd     MyIssuesCommand       `command:"my" description:"Show my issues"`
		ListCmd   ListIssuesCommand     `command:"list" description:"List issues"`
		ReadmeCmd GenerateReadmeCommand `command:"readme" description:"Generate README.md"`
	}
	opts.Opts = &Opts{}
	opts.InitCmd.Opts = opts.Opts
	opts.NewCmd.Opts = opts.Opts
	opts.MyCmd.Opts = opts.Opts
	opts.ListCmd.Opts = opts.Opts
	opts.ReadmeCmd.Opts = opts.Opts

	parser := flags.NewParser(&opts, flags.Default)
	parser.ShortDescription = "Tiny issue tracking in source code"
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		if opts.Verbose {
			log.SetLevel(log.InfoLevel)
		}
		for _, st := range strings.Split(opts.StatusesStr, ",") {
			opts.Statuses = append(opts.Statuses, strings.Trim(st, " "))
		}
		return command.Execute(args)
	}
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
