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
