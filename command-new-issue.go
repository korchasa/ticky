package main

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"
)

type NewIssueCommand struct {
}

func (c *NewIssueCommand) Execute(args []string) error {
	username, err := getUserFromGitConfig()
	if err != nil {
		return fmt.Errorf("can't get username from git config: %s", err)
	}

	i := Issue{
		Title:    strings.Join(args, " "),
		Assignee: []string{username},
		Status:   Flags.Statuses[0],
	}
	logrus.WithFields(logrus.Fields{"title":i.Title,"assignee":i.Assignee[0],"status":i.Status}).Info("Issue created")

	cont, err := i.MarshalYAML()
	if err != nil {
		return fmt.Errorf("can't marshal issue to YAML: %s", err)
	}
	logrus.WithField("symbols", len(cont)).Info("Content prepared")

	path := buildIssuesPath(i)
	if err = ioutil.WriteFile(path, cont, 0644); err != nil {
		return err
	}
	logrus.WithField("path", path).Info("File saved")

	return nil
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

func buildIssuesPath(i Issue) string {
	file := i.Title
	replace := "<>:\"/\\|?* "
	for _, char := range replace {
		file = strings.ReplaceAll(file, string(char), "-")
	}
	for strings.Contains(file, "--") {
		file = strings.ReplaceAll(file, "--", "-")
	}
	return fmt.Sprintf("%s/%s/%s.md", Flags.IssuesDir, i.Status, file)
}