package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"sort"
	"strings"
	"text/template"
)

type GenerateReadmeCommand struct {
	*Opts
	ReadmeFile string `long:"output" short:"o" description:"Output file" default:"README.md"`
	ReadmeTemplateFile string `long:"src" short:"s" description:"Template file" default:"_README.md"`
}

func (c *GenerateReadmeCommand) Execute(_ []string) error {
	iss, err := readAllIssues(Issue{}, c.Statuses, c.IssuesDir)
	if err != nil {
		return err
	}
	logrus.Infof("Founded %d issues in directory `%s`", len(iss), c.IssuesDir)

	b, err := BuildTable(iss, c.Statuses)
	if err != nil {
		return err
	}

	tpl, err := ioutil.ReadFile(c.ReadmeTemplateFile)
	if err != nil {
		return err
	}

	readme := strings.Replace(string(tpl), "[:ticky:]", b, 1)

	logrus.Infof("Write %d bytes to `%s`", len(readme), c.ReadmeFile)
	return ioutil.WriteFile(c.ReadmeFile, []byte(readme), 0644)
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

	funcMap := template.FuncMap{
		"esc": func(t string) string { return strings.ReplaceAll(t, "|", "\\|") },
	}

	t := template.Must(template.New("table").Funcs(funcMap).Parse(`
|{{range .Statuses }} {{. | esc }} |{{end}}
|{{range .Statuses }} --- |{{end}}
|{{range .Statuses }} {{range index $.Issues . }}- [{{.Title|esc}}]({{.File|esc}}){{range .Assignee}} [[@{{.|esc}}](https://github.com/{{.|esc}})]{{end}}<br/> {{end}}|{{end}}
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
