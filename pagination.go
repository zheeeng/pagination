/*
Package pagination wraps the response items with pagination infomation.

Example output in JSON format:

    {
        "pagination": {
            "page": 2,
            "page_size": 5,
            "total": 20,
            "first": "api.example.com/books?author=jk\u0026page=1\u0026page_size=5",
            "last": "api.example.com/books?author=jk\u0026page=4\u0026page_size=5",
            "prev": "api.example.com/books?author=jk\u0026page=1\u0026page_size=5",
            "next": "api.example.com/books?author=jk\u0026page=3\u0026page_size=5",
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

Get details: https://github.com/zheeeng/pagination
*/
package pagination

import (
	"github.com/zheeeng/pagination/pager"
	"github.com/zheeeng/pagination/queries"
)

type runInContext func(p *Paginator) Truncatable

// Pagination instance
type Pagination interface {
	Parse(link string) *Paginator
}

const defaultPageSize = 30

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

func (p *pagination) Parse(link string) *Paginator {
	basePath, page, pageSize, queries, hasPage, hasPageSize := queries.ParseLink(link, p.paginatorConfiguration.PageSize)

	pgt := &Paginator{
		pager:           pager.NewPager(page, pageSize),
		basePath:        basePath,
		queries:         queries,
		defaultPageSize: p.paginatorConfiguration.PageSize,
		hasPage:         hasPage,
		hasPageSize:     hasPageSize,
	}

	return pgt
}
