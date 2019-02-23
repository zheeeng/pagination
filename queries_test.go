package pagination

import (
	"fmt"
	"testing"
)

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
		}{
			{"happy input", "api.example.com/books?author=jk&page=2&pageSize=5",
				"api.example.com/books", 2, 5, "author=jk",
			},
			{"input with http scheme", "http://api.example.com/books?author=jk&page=2&pageSize=5",
				"http://api.example.com/books", 2, 5, "author=jk",
			},
			{"input with https scheme", "https://api.example.com/books?author=jk&page=2&pageSize=5",
				"https://api.example.com/books", 2, 5, "author=jk",
			},
			{"input without page", "api.example.com/books?author=jk&pageSize=5",
				"api.example.com/books", 1, 5, "author=jk",
			},
			{"input without pageSize", "api.example.com/books?author=jk&page=2",
				"api.example.com/books", 2, 30, "author=jk",
			},
			{"input without page and pageSize", "api.example.com/books?author=jk",
				"api.example.com/books", 1, 30, "author=jk",
			},
			{"input with multiple query terms", "api.example.com/books?author=jk&name=heaven",
				"api.example.com/books", 1, 30, "author=jk&name=heaven",
			},
		}

		for _, test := range tests {
			descr := fmt.Sprintf("\nTest %s failed:\n", test.testName)

			basePath, page, pageSize, queries := parseLink(test.link, defaultPageSize)

			if basePath != test.basePath {
				t.Errorf("%s[basePath]: got %s, want %s", descr, basePath, test.basePath)
			}
			if page != test.page {
				t.Errorf("%s[page]: got %d, want %d", descr, page, test.page)
			}
			if pageSize != test.pageSize {
				t.Errorf("%s[pageSize]: got %d, want %d", descr, pageSize, test.pageSize)
			}
			if queries.query.Encode() != test.queryEncoded {
				t.Errorf("%s[query encoded]: got %s, want %s", descr, queries.query.Encode(), test.queryEncoded)
			}

			queryBase := queries.query.Encode()
			queryTests := []string{queries.firstQuery.Encode(), queries.lastQuery.Encode(), queries.nextQuery.Encode(), queries.firstQuery.Encode()}
			for _, q := range queryTests {
				if q != queryBase {
					t.Errorf("%s[all queries are same]: got (%s != %s), want (%[2]s == %[3]s)", descr, queryBase, q)
					break
				}
			}
		}
	}
}
