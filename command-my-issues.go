package main

type MyIssuesCommand struct {
	Executed bool
	EArgs    []string
}

func (c *MyIssuesCommand) Execute(args []string) error {
	c.Executed = true
	c.EArgs = args
	return nil
}