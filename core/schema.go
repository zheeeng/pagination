package core

import "net/url"

type PaginationSchema struct {
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	Total    int        `json:"total"`
	First    string     `json:"first"`
	Last     string     `json:"last"`
	Previous string     `json:"previous"`
	Next     string     `json:"next"`
	Query    url.Values `json:"query"`
}

type Result struct {
	Pagination PaginationSchema `json:"pagination"`
	Result     interface{}      `json:"result"`
}
