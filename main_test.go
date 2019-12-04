package main

func init() {
	Flags = FlagsSpec{
		Statuses:  []string{"status1", "status2"},
		IssuesDir: "issues-test",
	}
}