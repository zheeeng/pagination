package queries

import (
	"fmt"
	"testing"
)

func TestParseLinkFail(t *testing.T) {
	defaultPageSize := 30
	var testLinks = []string{":::", ":", "::/..", "::/.", "::\\"}

	for i, testLink := range testLinks {
		basePath, page, pageSize, queries, hasPage, hasPageSize := ParseLink(testLink, defaultPageSize)

		if basePath != "" ||
			page != 0 || pageSize != 0 ||
			len(queries.FirstQuery) != 0 ||
			len(queries.LastQuery) != 0 ||
			len(queries.PrevQuery) != 0 ||
			len(queries.NextQuery) != 0 ||
			hasPage != false || hasPageSize != false {
			t.Errorf(
				"%d. invalid test link `%s` should be parsed to zero values result. got\npage: %d, pageSize: %d\nfirstQuery: %s, lastQuery: %s, prevQuery: %s, nextQuery: %s\nhasPage: %v, hasPageSize: %v",
				i, testLink,
				page, pageSize,
				queries.FirstQuery.Encode(),
				queries.LastQuery.Encode(),
				queries.PrevQuery.Encode(),
				queries.NextQuery.Encode(),
				hasPage, hasPageSize,
			)
		}
	}

}

func TestParseLink(t *testing.T) {
	defaultPageSize := 30

	{
		var tests = []struct {
			testName     string
			link         string
			basePath     string
			page         int
			pageSize     int
			queryEncoded string
			hasPage      bool
			hasPageSize  bool
		}{
			{"happy input", "api.example.com/books?author=jk&page=2&page_size=5",
				"api.example.com/books", 2, 5, "author=jk", true, true,
			},
			{"input with http scheme", "http://api.example.com/books?author=jk&page=2&page_size=5",
				"http://api.example.com/books", 2, 5, "author=jk", true, true,
			},
			{"input with https scheme", "https://api.example.com/books?author=jk&page=2&page_size=5",
				"https://api.example.com/books", 2, 5, "author=jk", true, true,
			},
			{"input without page", "api.example.com/books?author=jk&page_size=5",
				"api.example.com/books", 1, 5, "author=jk", false, true,
			},
			{"input with invalid page", "api.example.com/books?author=jk&page=foo&page_size=5",
				"api.example.com/books", 1, 5, "author=jk", false, true,
			},
			{"input without page_size", "api.example.com/books?author=jk&page=2",
				"api.example.com/books", 2, 30, "author=jk", true, false,
			},
			{"input with invalid page_size", "api.example.com/books?author=jk&page=2&page_size=bar",
				"api.example.com/books", 2, 30, "author=jk", true, false,
			},
			{"input without page and page_size", "api.example.com/books?author=jk",
				"api.example.com/books", 1, 30, "author=jk", false, false,
			},
			{"input with multiple query terms", "api.example.com/books?author=jk&name=heaven",
				"api.example.com/books", 1, 30, "author=jk&name=heaven", false, false,
			},
		}

		for i, test := range tests {
			descr := fmt.Sprintf("\n%d. Test %s failed:\n", i, test.testName)

			basePath, page, pageSize, queries, hasPage, hasPageSize := ParseLink(test.link, defaultPageSize)

			if basePath != test.basePath {
				t.Errorf("%s[basePath]: got %s, want %s", descr, basePath, test.basePath)
			}
			if page != test.page {
				t.Errorf("%s[page]: got %d, want %d", descr, page, test.page)
			}
			if pageSize != test.pageSize {
				t.Errorf("%s[pageSize]: got %d, want %d", descr, pageSize, test.pageSize)
			}
			if hasPage != test.hasPage {
				t.Errorf("%s[hasPage]: got %v, want %v", descr, hasPage, test.hasPage)
			}
			if hasPageSize != test.hasPageSize {
				t.Errorf("%s[hasPageSize]: got %v, want %v", descr, hasPageSize, test.hasPageSize)
			}
			if queries.Query.Encode() != test.queryEncoded {
				t.Errorf("%s[query encoded]: got %s, want %s", descr, queries.Query.Encode(), test.queryEncoded)
			}

			queryBase := queries.Query.Encode()
			queryTests := []string{
				queries.FirstQuery.Encode(),
				queries.LastQuery.Encode(),
				queries.NextQuery.Encode(),
				queries.FirstQuery.Encode(),
			}
			for _, q := range queryTests {
				if q != queryBase {
					t.Errorf("%s[all queries are same]: got (%s != %s), want (%[2]s == %[3]s)", descr, queryBase, q)
					break
				}
			}
		}
	}
}
