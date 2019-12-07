package main

type MyIssuesCommand struct {
	*Opts
}

func (c *MyIssuesCommand) Execute(a []string) error {
	user, err := getUserFromGitConfig()
	if err != nil {
		return err
	}

	subcmd := ListIssuesCommand{Assignee: user}
	return subcmd.Execute(a)
}