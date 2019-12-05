package main

import (
	"fmt"
	"sort"
)

type MyIssuesCommand struct {}

func (c *MyIssuesCommand) Execute(_ []string) error {
	user, err := getUserFromGitConfig()
	if err != nil {
		return err
	}
	iss, err := readAllIssues()
	if err != nil {
		return err
	}

	sortIssuesByStatusAndTitle(iss, Flags.Statuses)

	fmt.Printf("Tickets assigned to @%s:\n", user)
	for _, i := range iss {
		for _, a := range i.Assignee {
			if a == user {
				fmt.Printf("  - %s [%s]\n", i.Title, i.Status)
			}
		}
	}

	fmt.Println()

	return nil
}

func sortIssuesByStatusAndTitle(iss []Issue, statuses []string) {
	indexOf := func(data []string, element string) int {
		for k, v := range data {
			if element == v {
				return k
			}
		}
		return -1
	}
	sort.Slice(iss, func(i, j int) bool {
		if iss[i].Status != iss[j].Status {
			return indexOf(statuses, iss[i].Status) < indexOf(statuses, iss[j].Status)
		}
		return iss[i].Title < iss[j].Title
	})
}