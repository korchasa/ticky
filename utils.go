package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func readAllIssues() ([]Issue, error) {
	iss := make([]Issue, 0, 10)
	for _, status := range Flags.Statuses {
		files, err := filepath.Glob(fmt.Sprintf("%s/%s/*.md", Flags.IssuesDir, status))
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
			iss = append(iss, is)
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