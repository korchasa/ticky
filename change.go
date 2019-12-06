package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Change struct {
	Description string
	Breaking    bool
	Feature     bool
	Fix         bool
}

var (
	conventionalCommitRegex = regexp.MustCompile(`(?s)^(?P<type>\S+?)?(?P<scope>\(\S+\))?(?P<breaking>!?)?: (?P<desc>[^\n\r]+)?([\n\r]{2}(?P<body>.*))?`)
)

func (c *Change) UnmarshalString(data string) error {

	data = strings.Trim(data, " \n")

	match := conventionalCommitRegex.FindStringSubmatch(data)

	if len(match) == 0 {
		return fmt.Errorf("not a conventional commit message")
	}

	result := make(map[string]string)
	for i, name := range conventionalCommitRegex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	c.Breaking = result["breaking"] == "!" || strings.Contains(data, "BREAKING CHANGE")
	c.Description = result["desc"]

	t := result["type"]
	c.Feature = t == "feat" || t == "feature"
	c.Fix = t == "fix"

	return nil
}