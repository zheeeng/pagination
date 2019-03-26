package queries

import (
	"net/url"
	"strconv"
)

// PaginationQueries defines query fields
type PaginationQueries struct {
	Query      url.Values
	FirstQuery url.Values
	LastQuery  url.Values
	PrevQuery  url.Values
	NextQuery  url.Values
}

func (q *PaginationQueries) initPaginationQueries(u *url.URL) *PaginationQueries {
	q.Query = u.Query()
	q.FirstQuery = u.Query()
	q.LastQuery = u.Query()
	q.PrevQuery = u.Query()
	q.NextQuery = u.Query()

	return q
}

func (q *PaginationQueries) cleanPaginations() *PaginationQueries {
	q.Query.Del("page")
	q.Query.Del("page_size")
	q.FirstQuery.Del("page")
	q.FirstQuery.Del("page_size")
	q.LastQuery.Del("page")
	q.LastQuery.Del("page_size")
	q.PrevQuery.Del("page")
	q.PrevQuery.Del("page_size")
	q.NextQuery.Del("page")
	q.NextQuery.Del("page_size")

	return q
}

// ParseLink parse link to infomation components
func ParseLink(link string, defaultPageSize int) (
	basePath string,
	page, pageSize int,
	queries PaginationQueries,
	hasPage, hasPageSize bool,
) {
	parsedURL, err := url.Parse(link)

	if err != nil {
		return
	}

	if parsedURL.Scheme != "" {
		basePath = parsedURL.Scheme + "://"
	}

	basePath = basePath + parsedURL.Host + parsedURL.Path
	queries.initPaginationQueries(parsedURL)

	if queryPage := queries.Query.Get("page"); queryPage != "" {
		hasPage = true
		if page, err = strconv.Atoi(queryPage); err != nil {
			page = 1
		}
	} else {
		page = 1
	}

	if queryPageSize := queries.Query.Get("page_size"); queryPageSize != "" {
		hasPageSize = true
		if pageSize, err = strconv.Atoi(queryPageSize); err != nil {
			pageSize = defaultPageSize
		}
	} else {
		pageSize = defaultPageSize
	}

	queries.cleanPaginations()

	return
}
