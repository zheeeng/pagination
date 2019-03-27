<h1 align="center">:star: pagination</h1>

<div align="center">

A pagination wrapper for decorating resource list response.

![CircleCI](https://img.shields.io/circleci/project/github/zheeeng/pagination/master.svg?label=tests)
![Coverage Status](https://coveralls.io/repos/github/zheeeng/pagination/badge.svg)
![Language](https://img.shields.io/github/languages/top/zheeeng/pagination.svg?color=71e1ff)
![GitHub tag](https://img.shields.io/github/tag/zheeeng/pagination.svg)
![GitHub](https://img.shields.io/github/license/zheeeng/pagination.svg)
</div>

## :paperclip: Go doc

Get document at: https://godoc.org/github.com/zheeeng/pagination

## :fire: Features

1. Parse pagination info from request uri
    - Extract page and page size
    - Extract quires and feed them to paginated response
2. Decorate response body with pagination info
    - Feedback page navigation: `page`, `page_size`, `total`
    - Feedback hyper links: `first`, `last`, `prev`, `next`
    - Feedback query pairs
3. Get calculated valuable pagination params:
    - Get whether the URI provided pagination info
    - Calculate the offset and the chunk length
    - Calculate the start and end offsets, for manually truncate the list by yourself
    - Calculate values above from your specified page or item index
4. Manipulate pagination info
    - Modify quires
    - Reset pageSize, if you
5. Truncate resource list by demands
    - If the list length is greater than pageSize
6. Config default params
    - Change the default pageSize

## :bulb: Note

This pagination wrapper requires the resource list to be implemented with `Trunctable` interface. e.g.
```go
type TrunctableBooks []Book

func (tb TrunctableBooks) Slice(startIndex, endIndex int) pagination.Truncatable {
	return tb[startIndex:endIndex]
}
func (tb TrunctableBooks) Len() int {
	return len(tb)
}
```

## Example :point_down:

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

## Usage :point_down:

**Init default pagination configuration:**
```go
pg := pagination.DefaultPagination()
```

```go
pg := pagination.NewPagination(PaginatorConfiguration{
    PageSize: 50,
})
```

**Parse URI and get a paginator instance:**
```go
pgt := pg.Parse(someURI)

```

**Get/set page information**
```go
offset, length := pgt.GetOffsetRange()

total, items := db.Offset(offset).Limit(length).Query()
```

```go
start, end := pgt.GetRange()

total, items := db.QueryAll()
items = items[start:end]
```

```go
// Put all items into one page
if !pgt.HasRawPagination() {
    total, items := db.QueryAll()
    pgt.SetPageInfo(1, total)
}
```

**Wrap your list**

```go
response := pgt.Wrap(TruncatableItems(partialItems), total)
```

```go
// WrapWithTruncate helps truncating the list
response := pgt.WrapWithTruncate(TruncatableItems(allItems), total)
```

## Sample output :point_down:

```json
{
    "pagination": {
        "page": 2,
        "page_size": 5,
        "total": 20,
        "first": "api.example.com/books?author=jk&page=1&page_size=5",
        "last": "api.example.com/books?author=jk&page=4&page_size=5",
        "prev": "api.example.com/books?author=jk&page=1&page_size=5",
        "next": "api.example.com/books?author=jk&page=3&page_size=5",
        "query": {
            "author": [
                "jk"
            ],
            "page": [
                "2"
            ],
            "page_size": [
                "5"
            ]
        }
    },
    "result": [
        {
            "id": 5,
            "author": "jk",
            "name": "book"
        },
        {
            "id": 6,
            "author": "jk",
            "name": "book"
        },
        {
            "id": 7,
            "author": "jk",
            "name": "book"
        },
        {
            "id": 8,
            "author": "jk",
            "name": "book"
        },
        {
            "id": 9,
            "author": "jk",
            "name": "book"
        }
    ]
}
```
