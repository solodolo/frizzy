package processor

import (
	"testing"
)

func TestBuildPaginationContextReturnsCorrectNumberOfPages(t *testing.T) {
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

func TestBuildPaginationContextReturnsCorrectContentKeys(t *testing.T) {
	tests := []struct {
		numPaths, numPerPage, curPage, expected int
	}{
		{50, 8, 1, 8},    // 7 pages
		{50, 8, 8, 0},    // 7 pages
		{50, 5, 2, 5},    // 10 pages
		{1, 1, 1, 1},     // 1 page
		{0, 0, 0, 0},     // 0 pages
		{0, 200, 0, 0},   // 0 pages
		{2, 3, 8, 0},     // 1 page
		{2, 4000, 1, 2},  // 1 page
		{2, 4000, 0, 0},  // 1 page
		{100, 11, 10, 1}, // 10 pages
		{1, -3, 0, 0},    // 0 pages
	}

	for _, test := range tests {
		contentPaths := make([]string, test.numPaths)
		numPerPage := test.numPerPage
		curPage := test.curPage

		paginationContext := buildPaginationContext(contentPaths, curPage, numPerPage)
		content, ok := paginationContext.At("content")

		if !ok {
			t.Errorf("%v: expected context to have content key\n", test)
			return
		}

		contextsOnPage := content.child

		if len(*contextsOnPage) != test.expected {
			t.Errorf("%v: expected %d contexts on page, got %d\n", test, len(*contextsOnPage), test.expected)
		}
	}
}
