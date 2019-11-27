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
			content: `# Описать минимальный формат тикета @korchasa [in-progress]

Надо:

- [x] определить список полей
- [x] описать их в md
- [ ] проверить на наглядность diff'а

Поля: статус, ответственный

`,
			want: Issue{
				Title:    "Описать минимальный формат тикета",
				Assignee: []string{"korchasa"},
				Status:   "in-progress",
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
		wantErr error
	}{
		{
			name: "simple",
			iss: []Issue{
				{
					Title:    "td2",
					File: "./f1",
					Status:   "todo",
				},
				{
					Title:    "td1",
					File: "./f2",
					Status:   "todo",
				},
				{
					Title:    "ip",
					File: "./f3",
					Status:   "in-progress",
				},
			},
			wantS: `| todo | in-progress | done |
|---|---|---|
| [td1](./f2)<br/> [td2](./f1)<br/>| [ip](./f3)<br/>||`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := BuildTable(tt.iss, statuses)
			assert.Equal(t, tt.wantS, gotS)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}