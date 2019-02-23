/*
Package pagination wraps the response items with pagination infomation.

Example output in JSON format:

    {
        "pagination": {
            "page": 2,
            "pageSize": 5,
            "total": 20,
            "first": "api.example.com/books?author=jk\u0026page=1\u0026pageSize=5",
            "last": "api.example.com/books?author=jk\u0026page=4\u0026pageSize=5",
            "prev": "api.example.com/books?author=jk\u0026page=1\u0026pageSize=5",
            "next": "api.example.com/books?author=jk\u0026page=3\u0026pageSize=5",
            "query": {
                "author": [
                    "jk"
                ],
                "page": [
                    "2"
                ],
                "pageSize": [
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

Get details: https://github.com/zheeeng/pagination
*/
package pagination

import (
	"strconv"
)

type runInContext func(p Paginator) interface{}

// Pagination instance
type Pagination interface {
	Wrap(link string, r runInContext) Paginated
}

var defaultPageSize = 30

// PaginatorConfiguration defines the default pagination parameters. By default:
//
// -- PageSize: 30
type PaginatorConfiguration struct {
	PageSize int
}

type pagination struct {
	paginatorConfiguration PaginatorConfiguration
}

// DefaultPagination returns a default pagination instance
func DefaultPagination() Pagination {
	return &pagination{
		paginatorConfiguration: PaginatorConfiguration{
			PageSize: defaultPageSize,
		},
	}
}

// NewPagination create a fresh pagination instance
func NewPagination(cfg PaginatorConfiguration) Pagination {
	if cfg.PageSize == 0 {
		cfg.PageSize = defaultPageSize
	}

	return &pagination{
		paginatorConfiguration: cfg,
	}
}

func (p *pagination) Wrap(link string, run runInContext) Paginated {
	basePath, page, pageSize, queries := parseLink(link, p.paginatorConfiguration.PageSize)

	pgt := paginatorImpl{
		queries:         queries,
		defaultPageSize: p.paginatorConfiguration.PageSize,
	}

	pgt.SetIndicator(page, pageSize, 0)

	{
		result := run(&pgt)

		fields := PageFields{
			Page:     pgt.page,
			PageSize: pgt.pageSize,
			Total:    pgt.total,
			Query:    pgt.queries.query,
		}

		pgt.queries.query.Set("page", strconv.Itoa(pgt.page))
		pgt.queries.query.Set("pageSize", strconv.Itoa(pgt.pageSize))

		pgt.queries.firstQuery.Set("page", strconv.Itoa(pgt.firstPage))
		pgt.queries.firstQuery.Set("pageSize", strconv.Itoa(pgt.pageSize))
		fields.First = basePath + "?" + pgt.queries.firstQuery.Encode()

		if pgt.lastPage != 0 {
			pgt.queries.lastQuery.Set("page", strconv.Itoa(pgt.lastPage))
			pgt.queries.lastQuery.Set("pageSize", strconv.Itoa(pgt.pageSize))
			fields.Last = basePath + "?" + pgt.queries.lastQuery.Encode()
		}

		pgt.queries.prevQuery.Set("page", strconv.Itoa(pgt.prevPage))
		pgt.queries.prevQuery.Set("pageSize", strconv.Itoa(pgt.pageSize))
		fields.Prev = basePath + "?" + pgt.queries.prevQuery.Encode()

		pgt.queries.nextQuery.Set("page", strconv.Itoa(pgt.nextPage))
		pgt.queries.nextQuery.Set("pageSize", strconv.Itoa(pgt.pageSize))
		fields.Next = basePath + "?" + pgt.queries.nextQuery.Encode()

		return Paginated{
			Pagination: fields,
			Result:     result,
		}
	}
}
