package pagination

import (
	"net/url"
	"strconv"
)

type paginationQueries struct {
	query      url.Values
	firstQuery url.Values
	lastQuery  url.Values
	prevQuery  url.Values
	nextQuery  url.Values
}

func (q *paginationQueries) InitPaginationQueries(u *url.URL) *paginationQueries {
	q.query = u.Query()
	q.firstQuery = u.Query()
	q.lastQuery = u.Query()
	q.prevQuery = u.Query()
	q.nextQuery = u.Query()

	return q
}

func (q *paginationQueries) CleanAllPaginations() *paginationQueries {
	q.query.Del("page")
	q.query.Del("page_size")
	q.firstQuery.Del("page")
	q.firstQuery.Del("page_size")
	q.lastQuery.Del("page")
	q.lastQuery.Del("page_size")
	q.prevQuery.Del("page")
	q.prevQuery.Del("page_size")
	q.nextQuery.Del("page")
	q.nextQuery.Del("page_size")

	return q
}

func parseLink(link string, defaultPageSize int) (basePath string, page, pageSize int, queries paginationQueries, hasPage, hasPageSize bool) {
	parsedURL, err := url.Parse(link)

	if err != nil {
		return
	}

	if parsedURL.Scheme != "" {
		basePath = parsedURL.Scheme + "://"
	}

	basePath = basePath + parsedURL.Host + parsedURL.Path
	queries.InitPaginationQueries(parsedURL)

	if queryPage := queries.query.Get("page"); queryPage != "" {
		hasPage = true
		if page, err = strconv.Atoi(queryPage); err != nil {
			page = 1
		}
	} else {
		page = 1
	}

	if queryPageSize := queries.query.Get("page_size"); queryPageSize != "" {
		hasPageSize = true
		if pageSize, err = strconv.Atoi(queryPageSize); err != nil {
			pageSize = defaultPageSize
		}
	} else {
		pageSize = defaultPageSize
	}

	queries.CleanAllPaginations()

	return
}
