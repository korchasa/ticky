package main

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

type GenerateReadmeCommand struct {
	ReadmeFile string `long:"output" short:"o" description:"Output file" default:"README.md"`
	ReadmeTemplateFile string `long:"src" short:"s" description:"Template file" default:"_README.md"`
}

func (cmd *GenerateReadmeCommand) Execute(args []string) error {
	iss := make([]Issue, 0, 10)
	for _, status := range Flags.Statuses {
		files, err := filepath.Glob(fmt.Sprintf("%s/%s/*.md", Flags.IssuesDir, status))
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
	logrus.Infof("Founded %d issues in directory `%s`", len(iss), Flags.IssuesDir)
	b, err := BuildTable(iss, Flags.Statuses)
	if err != nil {
		return err
	}

	tpl, err := ioutil.ReadFile(cmd.ReadmeTemplateFile)
	if err != nil {
		return err
	}

	readme := strings.Replace(string(tpl), "[:ticky:]", b, 1)

	logrus.Infof("Write %d bytes to `%s`", len(readme), cmd.ReadmeFile)
	return ioutil.WriteFile(cmd.ReadmeFile, []byte(readme), 0644)
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
|{{range .Statuses }} {{range index $.Issues . }}[{{.Title}}]({{.File}}){{range .Assignee}} [[@{{.}}](https://github.com/{{.}})]{{end}}<br/> {{end}}|{{end}}
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
