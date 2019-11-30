package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

var (
	IssuesDir     = ".issues"
	ShowStatuses  = []string{"todo", "in-progress"}
	ReadmeFile    = "README.md"
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
	iss := make([]Issue, 0, 10)
	for _, status := range ShowStatuses {
		files, err := filepath.Glob(fmt.Sprintf("%s/%s/*.md", IssuesDir, status))
		if err != nil {
			panic(err)
		}
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
			is.Status = status
			iss = append(iss, is)
		}
	}
	b, err := BuildTable(iss, ShowStatuses)
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

	is.Title = strings.Trim(ls[0], " #")
	is.Text = strings.Trim(strings.Join(ls[1:], "\n"), "\n")
	return
}

func BuildTable(iss []Issue, statuses []string) (s string, err error) {
	issuesMap := make(map[string][]Issue)

	sort.Slice(iss, func(i, j int) bool {
		return iss[i].Title < iss[j].Title
	})

	for _, status := range statuses {
		issuesMap[status] = make([]Issue, 0, 3)
	}
	for _, is := range iss {
		issuesMap[is.Status] = append(issuesMap[is.Status], is)
	}

	t := template.Must(template.New("table").Parse(`
|{{range .Statuses }} {{.}} |{{end}}
|{{range .Statuses }} --- |{{end}}
|{{range .Statuses }} {{range index $.Issues . }}[{{.Title}}]({{.File}}) {{range .Assignee}}[@{{.}}](https://github.com/{{.}}){{end}}<br/> {{end}}|{{end}}
`))
	type data struct {
		Issues map[string][]Issue
		Statuses []string
	}

	buff := new(bytes.Buffer)
	if err = t.Execute(buff, data{issuesMap, statuses}); err != nil {
		return "", err
	}

	return strings.Trim(buff.String(), " \n"), nil
}
