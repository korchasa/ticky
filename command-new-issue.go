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
		logrus.Infof("can't get username from git config: %s", err)
	}
	fmt.Println(args, username)

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