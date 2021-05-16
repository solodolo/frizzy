package processor

import (
	"fmt"
	"testing"
)

func TestBuildPaginationContextReturnsCorrectNumberOfPages(t *testing.T) {
	tests := []struct {
		numPaths, numPerPage, expected int
	}{
		{50, 8, 7},
		{50, 5, 10},
		{1, 1, 1},
		{2, 3, 1},
		{2, 4000, 1},
		{100, 11, 10},
	}

	for _, test := range tests {
		contentPaths := make([]string, test.numPaths)
		numPerPage := test.numPerPage

		paginationContext, err := buildPaginationContext(contentPaths, 1, numPerPage)

		if err != nil {
			t.Errorf("%v: expected no errors, got %q\n", test, err)
			return
		}

		if numPages, ok := paginationContext.At("numPages"); ok {
			result := numPages.result.(IntResult)
			if int(result) != test.expected {
				t.Errorf("%v: expected %d pages, got %d", test, test.expected, int(result))
			}
		} else if numPerPage > 0 {
			t.Errorf("%v: expected context to contain \"numPages\" key", test)
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
		{2, 3, 8, 0},     // 1 page
		{2, 4000, 1, 2},  // 1 page
		{100, 11, 10, 1}, // 10 pages
	}

	for _, test := range tests {
		contentPaths := make([]string, test.numPaths)
		numPerPage := test.numPerPage
		curPage := test.curPage

		paginationContext, err := buildPaginationContext(contentPaths, curPage, numPerPage)

		if err != nil {
			t.Errorf("%v: expected no errors, got %q\n", test, err)
			return
		}

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

func TestInvalidCurPageReturnsError(t *testing.T) {
	var tests = []int{0, -1, -50}

	for _, test := range tests {
		expectedErr := fmt.Errorf("expected current page to be > 0, got %d", test)
		paginationContext, err := buildPaginationContext([]string{}, test, 1)

		if paginationContext != nil {
			t.Errorf("%d: expected pagination context to be nil, got %v", test, paginationContext)
		} else if err == nil || err.Error() != expectedErr.Error() {
			t.Errorf("%d: expected error %q, got %q", test, expectedErr, err)
		}
	}
}

func TestInvalidNumPerPageReturnsError(t *testing.T) {
	var tests = []int{0, -7, -200}

	for _, test := range tests {
		expectedErr := fmt.Errorf("expected number of items per page to be > 0, got %d", test)
		paginationContext, err := buildPaginationContext([]string{}, 1, test)

		if paginationContext != nil {
			t.Errorf("%d: expected pagination context to be nil, got %v", test, paginationContext)
		} else if err == nil || err.Error() != expectedErr.Error() {
			t.Errorf("%d: expected error %q, got %q", test, expectedErr, err)
		}
	}
}

func TestPagesBeforeReturnsCorrectContexts(t *testing.T) {
	var tests = []struct {
		curPage, numBefore, expectedNumPages int
	}{
		{1, 5, 0},
		{5, 1, 1},
		{0, 0, 0},
		{5, 5, 4},
	}

	inputPath := "/foo/bar/baz.html"
	for _, test := range tests {
		result, _ := PagesBefore(test.curPage, test.numBefore, inputPath)
		vals := result.context.Values()

		if len(vals) != test.expectedNumPages {
			t.Errorf("%v: expected %d contexts, got %d", test, test.expectedNumPages, len(vals))
		} else {
			for i, ctx := range vals {
				pageNumResult := (*ctx)["_pageNum"]
				hrefPathResult := (*ctx)["_pageHref"]

				gotPageNum := int(pageNumResult.result.(IntResult))
				gotPageHref := hrefPathResult.result.String()

				expectedPageNum := test.curPage - test.expectedNumPages + i
				expectedHref := GetMarkdownOutputPath(inputPath, expectedPageNum)
				if gotPageNum != expectedPageNum {
					t.Errorf("%v: expected page num to be %d, got %d", test, expectedPageNum, gotPageNum)
				} else if gotPageHref != expectedHref {
					t.Errorf("%v: expected page href to be %q, got %q", test, expectedHref, gotPageHref)
				}
			}
		}
	}
}

func TestPagesAfterReturnsCorrectContexts(t *testing.T) {
	var tests = []struct {
		curPage, numPages, numAfter, expectedNumPages int
	}{
		{1, 1, 5, 0},
		{5, 6, 1, 1},
		{0, 0, 0, 0},
		{5, 9, 5, 4},
	}

	inputPath := "/foo/bar/baz.html"
	for _, test := range tests {
		result, _ := PagesAfter(test.curPage, test.numPages, test.numAfter, inputPath)
		vals := result.context.Values()

		if len(vals) != test.expectedNumPages {
			t.Errorf("%v: expected %d contexts, got %d", test, test.expectedNumPages, len(vals))
		} else {
			for i, ctx := range vals {
				pageNumResult := (*ctx)["_pageNum"]
				hrefPathResult := (*ctx)["_pageHref"]

				gotPageNum := int(pageNumResult.result.(IntResult))
				gotPageHref := hrefPathResult.result.String()

				expectedPageNum := test.curPage + i + 1
				expectedHref := GetMarkdownOutputPath(inputPath, expectedPageNum)
				if gotPageNum != expectedPageNum {
					t.Errorf("%v: expected page num to be %d, got %d", test, expectedPageNum, gotPageNum)
				} else if gotPageHref != expectedHref {
					t.Errorf("%v: expected page href to be %q, got %q", test, expectedHref, gotPageHref)
				}
			}
		}
	}
}
