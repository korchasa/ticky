package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalIssue(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    Issue
		wantErr error
	}{
		{
			name: "positive",
			content: `# Описать минимальный формат тикета @korchasa

Надо:

- [x] определить список полей
- [x] описать их в md
- [ ] проверить на наглядность diff'а

Поля: статус, ответственный

`,
			want: Issue{
				Title:    "Описать минимальный формат тикета",
				Assignee: []string{"korchasa"},
				Text:     `Надо:

- [x] определить список полей
- [x] описать их в md
- [ ] проверить на наглядность diff'а

Поля: статус, ответственный`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalIssue(tt.content)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBuildTable(t *testing.T) {
	statuses := []string{"todo", "in-progress", "done"}
	tests := []struct {
		name    string
		iss      []Issue
		wantS   string
	}{
		{
			name: "simple",
			iss: []Issue{
				{
					Title:    "td2",
					File: "./f1",
					Status:   "todo",
					Assignee: []string{"user1"},
				},
				{
					Title:    "td1",
					File: "./f2",
					Status:   "todo",
					Assignee: []string{"user2"},
				},
				{
					Title:    "ip",
					File: "./f3",
					Status:   "in-progress",
					Assignee: []string{"user3"},
				},
				{
					Title:    "arch",
					File: "./f4",
					Status:   "archive",
					Assignee: []string{"user4"},
				},
			},
			wantS: `| todo | in-progress | done |
| --- | --- | --- |
| [td1](./f2) [[@user2](https://github.com/user2)]<br/> [td2](./f1) [[@user1](https://github.com/user1)]<br/> | [ip](./f3) [[@user3](https://github.com/user3)]<br/> | |`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := BuildTable(tt.iss, statuses)
			assert.Equal(t, tt.wantS, gotS)
			if err != nil {
				assert.Failf(t, "error on table build: %s", err.Error())
			}
		})
	}
}