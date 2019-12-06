package main

import (
	"fmt"
)

type ListIssuesCommand struct {
	Assignee   string `long:"assignee" short:"a" description:"Filter by assignee"`
	Status   string `long:"status" short:"s" description:"Filter by status"`
}

func (c *ListIssuesCommand) Execute(_ []string) error {
	iss, err := readAllIssues(Issue{Status:c.Status, Assignee: []string{c.Assignee}})
	if err != nil {
		return err
	}
	fmt.Println("Issues:")
	fmt.Println(buildIssuesListText(iss))

	return nil
}
