package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

func readAllIssues(tpl Issue, statuses []string, issuesDir string) ([]Issue, error) {
	iss := make([]Issue, 0, 10)
	for _, status := range statuses {
		files, err := filepath.Glob(fmt.Sprintf("%s/%s/*.md", issuesDir, status))
		if err != nil {
			return iss, err
		}
		for _, f := range files {
			b, err := ioutil.ReadFile(f)
			if err != nil {
				return iss, err
			}
			var is Issue
			if err := is.UnmarshalYAML(b); err != nil {
				return iss, err
			}
			is.File = f
			is.Status = status
			if is.Match(tpl) {
				iss = append(iss, is)
			}
		}
	}
	return iss, nil
}

func getUserFromGitConfig() (string, error) {
	var stdout bytes.Buffer
	cmd := exec.Command("git", []string{"config", "--get", "--null", "user.name"}...)
	cmd.Stdout = &stdout
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 {
				return "", fmt.Errorf("git config not found")
			}
		}
		return "", err
	}

	return strings.TrimRight(stdout.String(), "\000"), nil
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

func buildIssuesListText(iss []Issue, statuses []string) string {
	var buff []string
	sortIssuesByStatusAndTitle(iss, statuses)

	maxTitle := 20
	maxStatus := 0
	for _, i := range iss {
		if len(i.Title) > maxTitle { maxTitle = len(i.Title) }
		if len(i.Status) > maxStatus { maxStatus = len(i.Status) }
	}

	for _, i := range iss {
		var ass []string
		for _, a := range i.Assignee {
			ass = append(ass, "@" + a)
		}
		buff = append(
			buff,
			fmt.Sprintf(
				"      %s    %s    %s",
				i.Title + strings.Repeat(" ", maxTitle - len(i.Title)),
				i.Status + strings.Repeat(" ", maxStatus - len(i.Status)),
				strings.Join(ass, ", ")))
	}

	return strings.Join(buff, "\n") + "\n"
}