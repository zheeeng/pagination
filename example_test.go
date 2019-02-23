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

var total = 20
var requestURI = "api.example.com/books?author=jk&page=2&pageSize=5"
var books = []Book{}

func init() {
	for i := 0; i < 20; i++ {
		book := Book{i, "jk", "book"}
		books = append(books, book)
	}
}

func Example() {
	pg := pagination.DefaultPagination()

	paginatedData := pg.Wrap(
		requestURI,
		func(pgt *pagination.Paginator) pagination.Truncatable {
			pgt.SetTotal(total)
			return pgt.WrapWithTruncate(TrunctableBooks(books))
		})

	responseBody, _ := json.MarshalIndent(paginatedData, "", "    ")

	fmt.Println(string(responseBody))
}
