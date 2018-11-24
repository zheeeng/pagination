package pagination

import (
	"github.com/zheeeng/pagination/core"
)

type (
	PaginationSchema       = core.PaginationSchema
	Pagination             = core.Pagination
	Paginator              = core.Paginator
	PaginatorConfiguration = core.PaginatorConfiguration
)

var (
	DefaultPagination = core.DefaultPagination
	NewPagination     = core.NewPagination
)
