package main

import "fmt"

type NewIssueCommand struct {
}

func (c *NewIssueCommand) Execute(args []string) error {
	fmt.Println(args)
	return nil
}