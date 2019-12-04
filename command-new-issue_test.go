package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_buildIssuesPath(t *testing.T) {
	tests := []struct {
		name string
		i Issue
		want string
	}{
		{
			"simple",
			Issue{
				Title:    "foo bar",
				Status:   "status1",
			},
			"issues-test/status1/foo-bar.md",
		},
		{
			"unicode",
			Issue{
				Title:    "привет мир",
				Status:   "status2",
			},
			"issues-test/status2/привет-мир.md",
		},
		{
			"cleanup",
			Issue{
				Title:    "black/white > 1:2",
				Status:   "status2",
			},
			"issues-test/status2/black-white-1-2.md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, buildIssuesPath(tt.i))
		})
	}
}