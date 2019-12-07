package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type NewIssueCommand struct {
	*Opts
}

func (c *NewIssueCommand) Execute(args []string) error {
	username, err := getUserFromGitConfig()
	if err != nil {
		return fmt.Errorf("can't get username from git config: %s", err)
	}

	i := Issue{
		Title:    strings.Join(args, " "),
		Assignee: []string{username},
		Status:   c.Statuses[0],
	}
	logrus.WithFields(logrus.Fields{"title": i.Title, "assignee": i.Assignee[0], "status": i.Status}).Info("Issue created")

	cont, err := i.MarshalYAML()
	if err != nil {
		return fmt.Errorf("can't marshal issue to YAML: %s", err)
	}
	logrus.WithField("symbols", len(cont)).Info("Content prepared")

	path := c.buildIssuesPath(i)
	if err = ioutil.WriteFile(path, cont, 0644); err != nil {
		return err
	}
	logrus.WithField("path", path).Info("File saved")

	editor := os.Getenv("EDITOR")
	if editor != "" {
		err = exec.Command(editor, path).Run()
		if err != nil {
			return fmt.Errorf("editor finished with error: %s", err)
		}
	}

	return nil
}

func (c *NewIssueCommand) buildIssuesPath(i Issue) string {
	file := i.Title
	replace := "<>:\"/\\|?*() "
	for _, char := range replace {
		file = strings.ReplaceAll(file, string(char), "-")
	}
	for strings.Contains(file, "--") {
		file = strings.ReplaceAll(file, "--", "-")
	}
	return fmt.Sprintf("%s/%s/%s.md", c.IssuesDir, i.Status, strings.Trim(file, "-"))
}
