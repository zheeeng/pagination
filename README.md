<h1 align="center">:star: pagination</h1>

<div align="center">

A pagination wrapper for decorating resource list response.

[![CircleCI](https://img.shields.io/circleci/project/github/zheeeng/pagination/master.svg?label=tests)](https://circleci.com/gh/zheeeng/pagination)
[![Coveralls](https://img.shields.io/coveralls/github/zheeeng/pagination.svg)](https://circleci.com/api/v1.1/project/github/zheeeng/pagination/latest/artifacts/0/tmp/artifacts/coverage.html)
[![Language](https://img.shields.io/github/languages/top/zheeeng/pagination.svg?color=71e1ff)](https://golang.org/)
[![Release](https://img.shields.io/github/tag/zheeeng/pagination.svg)](https://github.com/zheeeng/pagination/releases)
[![License](https://img.shields.io/github/license/zheeeng/pagination.svg)](https://github.com/zheeeng/pagination/blob/master/LICENSE)
</div>

## :paperclip: Go doc

Get document at: https://godoc.org/github.com/zheeeng/pagination

## :fire: Features

1. Parse pagination info from request URI:
    - Extract page and page size
    - Extract quires and feed them to paginated response
2. Decorate response body with pagination info:
    - Feedback page navigation: `page`, `page_size`, `total`
    - Feedback hyper links: `first`, `last`, `prev`, `next`
    - Feedback query pairs
3. Get calculated valuable pagination params:
    - Get whether the URI provided pagination info
    - Calculate the offset and the chunk length
    - Calculate the start and end offsets, for manually truncate the list by yourself
    - Calculate values above from your specified page or item index
4. Manipulate pagination info:
    - Modify quires
    - Reset page and pageSize, maybe sometimes you want to overwrite them
5. Truncate resource list by demands:
    - If the list length is greater than pageSize
6. Config default params:
    - Change the default page size

## :bulb: Note

This pagination wrapper requires the resource list implements `Trunctable` interface.

```go
type Truncatable interface {
    Len() int
    Slice(startIndex, endIndex int) Truncatable
}
```

e.g.
```go
type TrunctableBooks []Book

func (tb TrunctableBooks) Slice(startIndex, endIndex int) pagination.Truncatable {
	return tb[startIndex:endIndex]
}
func (tb TrunctableBooks) Len() int {
	return len(tb)
}
```

## Usage :point_down:

**Init a pagination instance:**
```go
pg := pagination.DefaultPagination()
```

```go
pg := pagination.NewPagination(PaginatorConfiguration{
    PageSize: 50,
})
```

**Parse URI and get a manipulable paginator**
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

**Manipulate queries**

```go
// got url.Values
query := pgt.Query()
if query.Get("publisher") == "" {
	query.Add("publisher", "Tada Publications")
}
if query.Get("author") == "Fisher" {
	query.Set("author", "F.isher")
}
schema.Parse(query.Encode(), &someQueryBookStruct)
```

**Wrap your list**

```go
response := pgt.Wrap(TruncatableItems(partialItems), total)
```

```go
// WrapWithTruncate helps truncating the list
response := pgt.WrapWithTruncate(TruncatableItems(allItems), total)
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
