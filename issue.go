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

func (i *Issue) UnmarshalYAML(data []byte) error {
	ls := strings.Split(string(data), "\n")
	if len(ls) == 0 {
		return fmt.Errorf("can't find lines in file")
	}

	re := regexp.MustCompile(`@([^\s]*)`)
	for _, match := range re.FindAllString(ls[0], -1) {
		i.Assignee = append(i.Assignee, match[1:])
		ls[0] = strings.Replace(ls[0], match, "", 1)
	}

	i.Title = strings.Trim(ls[0], " #")
	i.Text = strings.Trim(strings.Join(ls[1:], "\n"), "\n")
	return nil
}

func (i *Issue) MarshalYAML() ([]byte, error) {
	assignees := ""
	for _, as := range i.Assignee {
		assignees += "@" + as
	}
	return []byte(fmt.Sprintf("# %s %s\n\n%s\n\n", i.Title, assignees, i.Text)), nil
}

func (i *Issue) Match(tpl Issue) bool {
	if tpl.File != "" && tpl.File != i.File { return false }
	if tpl.Title != "" && tpl.Title != i.Title { return false }
	if tpl.Status != "" && tpl.Status != i.Status { return false }
	if len(tpl.Assignee) > 0 {
		founded := false
		for _, ta := range tpl.Assignee {
			for _, ia := range i.Assignee {
				if ta == ia {
					founded = true
				}
			}
		}
		if !founded {
			return false
		}
	}
	if tpl.Text != "" && tpl.Text != i.Text { return false }
	return true
}
