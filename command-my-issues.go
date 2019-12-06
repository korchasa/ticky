package main

type MyIssuesCommand struct {
}

func (c *MyIssuesCommand) Execute(a []string) error {
	user, err := getUserFromGitConfig()
	if err != nil {
		return err
	}

	cmd := ListIssuesCommand{Assignee: user}
	return cmd.Execute(a)
}