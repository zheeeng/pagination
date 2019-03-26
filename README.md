# pagination
A pagination wrapper for decorating resource list response.

# Go doc

Get document at: https://godoc.org/github.com/zheeeng/pagination

# Features

1. Parse pagination info from request uri
    - Extract page and page size
    - Extract quires and feed them to paginated response
2. Decorate response body with pagination info
    - Feedback page navigation: `page`, `page_size`, `total`
    - Feedback hyper links: `first`, `last`, `prev`, `next`
    - Feedback query pairs
3. Manipulate pagination info
    - Modify quires
    - Reset pageSize, if you
4. Truncate resource list by demands
    - If the list length is greater than pageSize
5. Config default params
    - Change the default pageSize

# Example

```go
package pagination_test

import (
	"encoding/json"
	"fmt"

	"github.com/zheeeng/pagination"
)

type Book struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Name   string `json:"name"`
}

type TrunctableBooks []Book

func (tb TrunctableBooks) Slice(startIndex, endIndex int) pagination.Truncatable {
	return tb[startIndex:endIndex]
}
func (tb TrunctableBooks) Len() int {
	return len(tb)
}

var total = 20
var requestURI = "api.example.com/books?author=jk&page=2&page_size=5"
var books = []Book{}

func init() {
	for i := 0; i < 20; i++ {
		book := Book{i, "jk", "book"}
		books = append(books, book)
	}
}

func Example() {
	pg := pagination.DefaultPagination()

	pgt := pg.Parse(requestURI)

	paginatedData := pgt.WrapWithTruncate(TrunctableBooks(books), total)

	responseBody, _ := json.MarshalIndent(paginatedData, "", "    ")

	fmt.Println(string(responseBody))
}

```
