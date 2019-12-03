package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type Issue struct {
	File string
	Title string
	Assignee []string
	Status string
	Text string
}

func (i Issue) String() string {
	js, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(js)
}

func UnmarshalIssue(c string) (is Issue, err error) {
	ls := strings.Split(c, "\n")
	if len(ls) == 0 {
		return is, fmt.Errorf("can't find lines in file")
	}

	re := regexp.MustCompile(`@([^\s]*)`)
	for _, match := range re.FindAllString(ls[0], -1) {
		is.Assignee = append(is.Assignee, match[1:])
		ls[0] = strings.Replace(ls[0], match, "", 1)
	}

	is.Title = strings.Trim(ls[0], " #")
	is.Text = strings.Trim(strings.Join(ls[1:], "\n"), "\n")
	return
}