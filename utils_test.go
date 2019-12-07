package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_sortIssuesByStatusAndTitle(t *testing.T) {
	{
		byStatus := []Issue{{Status: "sa", Title: "a"}, {Status: "sb", Title: "z"}}
		wantByStatus := []Issue{{Status: "sb", Title: "z"}, {Status: "sa", Title: "a"}}
		sortIssuesByStatusAndTitle(byStatus, []string{"sb", "sa"})
		assert.Equal(t, wantByStatus, byStatus)
	}
	{
		byTitle := []Issue{{Title: "z"}, {Title: "a"}}
		wantByTitle := []Issue{{Title: "a"}, {Title: "z"}}
		sortIssuesByStatusAndTitle(byTitle, []string{"sb", "sa"})
		assert.Equal(t, wantByTitle, byTitle)
	}
}