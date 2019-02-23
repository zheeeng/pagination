package pagination

import "net/url"

// PageFields defines the struct of pagination field
type PageFields struct {
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
	Total    int        `json:"total"`
	First    string     `json:"first"`
	Last     string     `json:"last"`
	Prev     string     `json:"prev"`
	Next     string     `json:"next"`
	Query    url.Values `json:"query"`
}

// Paginated defines the paginated response struct
type Paginated struct {
	Pagination PageFields  `json:"pagination"`
	Result     Truncatable `json:"result"`
}
