package pagination

import (
	"github.com/zheeeng/pagination/core"
)

type (
	PaginationSchema       = core.PaginationSchema
	Paginator              = core.Paginator
	PaginatorConfiguration = core.PaginatorConfiguration
)

var (
	DefaultPagination = core.DefaultPagination
	NewPagination     = core.NewPagination
)
