package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommit_UnmarshalString(t *testing.T) {
	tests := []struct {
		data string
		c    Change
	}{
		{
			"feat: allow provided config object to extend other configs\nBREAKING CHANGE: 'extends' key in config",
			Change{"allow provided config object to extend other configs", true, true, false},
		},
		{
			"refactor!: drop support for Node 6",
			Change{"drop support for Node 6", true, false, false},
		},
		{
			"docs: correct spelling of CHANGELOG",
			Change{"correct spelling of CHANGELOG", false, false, false},
		},
		{
			"feat(lang): add polish language",
			Change{"add polish language", false, true, false},
		},
		{
			"fix: correct minor typos in code\n\nsee the issue for details\n\non typos fixed.\n\nReviewed-by: Z\nRefs #133",
			Change{"correct minor typos in code", false, false, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			var c Change
			err := c.UnmarshalString(tt.data)
			assert.NoError(t, err)
			assert.Equal(t, tt.c.Description, c.Description)
			assert.Equal(t, tt.c.Breaking, c.Breaking)
			assert.Equal(t, tt.c.Feature, c.Feature)
			assert.Equal(t, tt.c.Fix, c.Fix)
		})
	}
}
