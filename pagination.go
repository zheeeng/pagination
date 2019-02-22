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
	"net/url"
	"strconv"
)

var defaultPageSize = 30

type runInContext func(p Paginator) interface{}

// Pagination instance
type Pagination interface {
	Wrap(link string, r runInContext) Paginated
}

// PaginatorConfiguration defines the default pagination parameters
type PaginatorConfiguration struct {
	PageSize int
}

type pagination struct {
	defaultPaginatorConfiguration PaginatorConfiguration
}

// DefaultPagination returns a default pagination instance
func DefaultPagination() Pagination {
	return &pagination{
		defaultPaginatorConfiguration: PaginatorConfiguration{
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
		defaultPaginatorConfiguration: cfg,
	}
}

type v = url.Values

func parseQueries(u *url.URL) (getQueries func() (v, v, v, v, v), cleanPaginationInQueries func()) {
	q := u.Query()
	f := u.Query()
	l := u.Query()
	p := u.Query()
	n := u.Query()

	getQueries = func() (v, v, v, v, v) {
		return q, f, l, p, n
	}

	cleanPaginationInQueries = func() {
		q.Del("page")
		q.Del("page_size")
		f.Del("page")
		f.Del("page_size")
		l.Del("page")
		l.Del("page_size")
		p.Del("page")
		p.Del("page_size")
		n.Del("page")
		n.Del("page_size")
	}

	return
}

func (p *pagination) Wrap(link string, run runInContext) Paginated {
	parsedURL, err := url.Parse(link)
	basePath := ""
	var getQueries func() (v, v, v, v, v)
	var cleanPaginationInQueries func()
	if err == nil {
		if parsedURL.Scheme != "" {
			basePath = parsedURL.Scheme + "://"
		}
		basePath = basePath + parsedURL.Host + parsedURL.Path
		getQueries, cleanPaginationInQueries = parseQueries(parsedURL)
	}

	query, firstQuery, lastQuery, previousQuery, nextQuery := getQueries()

	page := 0
	queryPage := query.Get("page")
	if queryPage != "" {
		page, err = strconv.Atoi(queryPage)
		if err != nil {
			page = 0
		}
	}

	pageSize := 0
	queryPageSize := query.Get("page_size")
	if queryPageSize != "" {
		pageSize, err = strconv.Atoi(queryPageSize)
		if err != nil {
			pageSize = p.defaultPaginatorConfiguration.PageSize
		}
	}

	cleanPaginationInQueries()

	pgt := paginatorImpl{
		Query:           query,
		FirstQuery:      firstQuery,
		LastQuery:       lastQuery,
		PreviousQuery:   previousQuery,
		NextQuery:       nextQuery,
		defaultPageSize: p.defaultPaginatorConfiguration.PageSize,
	}

	pgt.SetIndicator(page, pageSize, 0)

	{
		result := run(&pgt)

		pgt.Query.Set("page", strconv.Itoa(pgt.page))
		pgt.Query.Set("page_size", strconv.Itoa(pgt.pageSize))

		first := ""
		pgt.FirstQuery.Set("page", strconv.Itoa(pgt.firstPage))
		pgt.FirstQuery.Set("page_size", strconv.Itoa(pgt.pageSize))
		first = basePath + "?" + pgt.FirstQuery.Encode()

		last := ""
		if pgt.lastPage != 0 {
			pgt.LastQuery.Set("page", strconv.Itoa(pgt.lastPage))
			pgt.LastQuery.Set("page_size", strconv.Itoa(pgt.pageSize))
			last = basePath + "?" + pgt.LastQuery.Encode()
		}

		previous := ""
		pgt.PreviousQuery.Set("page", strconv.Itoa(pgt.previousPage))
		pgt.PreviousQuery.Set("page_size", strconv.Itoa(pgt.pageSize))
		previous = basePath + "?" + pgt.PreviousQuery.Encode()

		next := ""
		pgt.NextQuery.Set("page", strconv.Itoa(pgt.nextPage))
		pgt.NextQuery.Set("page_size", strconv.Itoa(pgt.pageSize))
		next = basePath + "?" + pgt.NextQuery.Encode()

		return Paginated{
			Pagination: PageFields{
				Page:     pgt.page,
				PageSize: pgt.pageSize,
				Total:    pgt.total,
				First:    first,
				Last:     last,
				Previous: previous,
				Next:     next,
				Query:    pgt.Query,
			},
			Result: result,
		}
	}
}
