package pagination

import (
	"fmt"
	"testing"
)

func TestParseLink(t *testing.T) {
	defaultPageSize := 30

	{
		var tests = []struct {
			testName string
			link     string
			basePath string
			page     int
			pageSize int
		}{
			{"happy input", "api.example.com/books?author=jk&page=2&pageSize=5", "api.example.com/books", 2, 5},
			{"input with http scheme", "api.example.com/books?author=jk&page=2&pageSize=5", "api.example.com/books", 2, 5},
			{"input with http scheme", "api.example.com/books?author=jk&page=2&pageSize=5", "api.example.com/books", 2, 5},
			{"input without page", "api.example.com/books?author=jk&pageSize=5", "api.example.com/books", 1, 5},
			{"input without pageSize", "api.example.com/books?author=jk&page=2", "api.example.com/books", 2, 30},
			{"input without page and pageSize", "api.example.com/books?author=jk", "api.example.com/books", 1, 30},
		}

		for _, test := range tests {
			descr := fmt.Sprintf("\nTest %s failed:\n", test.testName)

			basePath, page, pageSize, _ := parseLink(test.link, defaultPageSize)

			if basePath != test.basePath {
				t.Errorf("%s[basePath]: got %s, want %s", descr, basePath, test.basePath)
			}
			if page != test.page {
				t.Errorf("%s[page]: got %d, want %d", descr, page, test.page)
			}
			if pageSize != test.pageSize {
				t.Errorf("%s[pageSize]: got %d, want %d", descr, pageSize, test.pageSize)
			}
		}
	}
}
