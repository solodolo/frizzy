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
		{1, -3, 0},
	}

	for _, test := range tests {
		contentPaths := make([]string, test.numPaths)
		numPerPage := test.numPerPage

		paginationContext := buildPaginationContext(contentPaths, 1, numPerPage)
		if numPages, ok := paginationContext.At("numPages"); ok {
			result := numPages.result.(IntResult)
			if int(result) != test.expected {
				t.Errorf("expected %d pages, got %d", test.expected, int(result))
			}
		} else if numPerPage > 0 {
			t.Errorf("expected context to contain \"numPages\" key")
		}
	}
}
