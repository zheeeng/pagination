package core

import (
	"net/url"
	"strconv"
)

type configure func(p Paginator)

type Pagination interface {
	NewWrapper(link string, c configure) func(interface{}) Result
}

type PaginatorConfiguration struct {
	pageSize int
}

type pagination struct {
	defaultPaginatorConfiguration PaginatorConfiguration
}

func DefaultPagination() Pagination {
	return &pagination{
		defaultPaginatorConfiguration: PaginatorConfiguration{
			pageSize: 30,
		},
	}
}

func NewPagination(cfg PaginatorConfiguration) Pagination {
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

func (p *pagination) NewWrapper(link string, h configure) func(interface{}) Result {
	parsedUrl, err := url.Parse(link)
	basePath := ""
	var getQueries func() (v, v, v, v, v)
	var cleanPaginationInQueries func()
	if err == nil {
		if parsedUrl.Scheme != "" {
			basePath = parsedUrl.Scheme + "://"
		}
		basePath = basePath + parsedUrl.Host + parsedUrl.Path
		getQueries, cleanPaginationInQueries = parseQueries(parsedUrl)
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
			pageSize = p.defaultPaginatorConfiguration.pageSize
		}
	}

	cleanPaginationInQueries()

	pgt := paginatorImpl{
		Query:           query,
		FirstQuery:      firstQuery,
		LastQuery:       lastQuery,
		PreviousQuery:   previousQuery,
		NextQuery:       nextQuery,
		defaultPageSize: p.defaultPaginatorConfiguration.pageSize,
	}

	pgt.SetIndicator(page, pageSize, 0)

	return func(result interface{}) Result {
		h(&pgt)

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

		return Result{
			Pagination: PaginationSchema{
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
