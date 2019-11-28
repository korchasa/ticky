package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	IssuesDir = ".issues"
	DefaultStatuses = []string{"todo", "in-progress", "done"}
	ReadmeFile = "README.md"
	ReadmeAltFile = "_README.md"
)

type Config struct {
	IssuesPath string
	Statuses   []string
}

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

func main() {
	conf := Config{
		IssuesPath: fmt.Sprintf("%s/*.md", IssuesDir),
		Statuses:   DefaultStatuses,
	}
	files, err := filepath.Glob(conf.IssuesPath)
	if err != nil {
		panic(err)
	}
	iss := make([]Issue, 0, len(files))
	for _, f := range files {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}
		is, err := UnmarshalIssue(string(b))
		if err != nil {
			panic(err)
		}
		is.File = f
		iss = append(iss, is)
	}
	b, err := BuildTable(iss, conf.Statuses)
	if err != nil {
		panic(err)
	}

	tpl, err := ioutil.ReadFile(ReadmeAltFile)
	if err != nil {
		panic(err)
	}

	readme := strings.Replace(string(tpl), "[:ticky:]", b, 1)

	err = ioutil.WriteFile(ReadmeFile, []byte(readme), 0644)
	if err != nil {
		panic(err)
	}
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

	re = regexp.MustCompile(`\[([^]]*)\]`)
	statuses := re.FindAllString(ls[0], -1)
	switch len(statuses) {
	case 0:
		return is, fmt.Errorf("can'r find any status in line: %s", ls[0])
	case 1:
		is.Status = statuses[0][1:len(statuses[0])-1]
		ls[0] = strings.Replace(ls[0], statuses[0], "", 1)
	default:
		return is, fmt.Errorf("too many statuses in line: %s", ls[0])
	}

	is.Title = strings.Trim(ls[0], " #")
	is.Text = strings.Trim(strings.Join(ls[1:], "\n"), "\n")
	return
}

func BuildTable(iss []Issue, statuses []string) (s string, err error) {
	isps := make(map[string][]Issue)

	for _, status := range statuses {
		isps[status] = make([]Issue, 0, 3)
		s = fmt.Sprintf("%s| %s ", s, status)
	}
	s = fmt.Sprintf("%s|\n", s)

	for i := 0; i < len(statuses); i++ {
		s = fmt.Sprintf("%s|---", s)
	}
	s = fmt.Sprintf("%s|\n", s)

	sort.Slice(iss, func(i, j int) bool {
		return iss[i].Title < iss[j].Title
	})

	for _, is := range iss {
		isps[is.Status] = append(isps[is.Status], is)
	}

	for _, status := range statuses {
		s = fmt.Sprintf("%s|", s)
		for _, is := range isps[status] {
			s = fmt.Sprintf("%s [%s](%s)<br/>", s, is.Title, is.File)
		}
	}
	s = fmt.Sprintf("%s|", s)
	return
}
