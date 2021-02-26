package processor

import (
	"testing"
)

func TestBuildPaginationContextsContainsCorrectNumberOfContexts(t *testing.T) {
	tests := []struct {
		numPaths, numPerPage, expected int
	}{
		{50, 8, 7},
		{50, 5, 10},
		{1, 1, 1},
		{0, 0, 0},
		{0, 200, 0},
		{2, 3, 1},
		{2, 4000, 1},
		{100, 11, 10},
	}

	for _, test := range tests {
		contentPaths := make([]string, test.numPaths)
		numPerPage := test.numPerPage

		paginationContexts := buildPaginationContexts(contentPaths, numPerPage)

		if len(paginationContexts) != test.expected {
			t.Errorf("expected %d pages, got %d", test.expected, len(paginationContexts))
		}
	}
}
