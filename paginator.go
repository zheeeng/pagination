package pagination

import (
	"errors"
	"reflect"
)

// WrappedItems is a nominal type for Pagination::Wrap func consuming.
// It ensure user always call Paginator::Wrap or Paginator::WrapWithTruncate to return items
type WrappedItems interface{}

// Paginator provides methods to manipulate pagination fields
type Paginator interface {
	Wrap(items interface{}) WrappedItems
	WrapWithTruncate(items interface{}) WrappedItems
	GetPaginationRangeByPage(page int) (startIndex, endIndex int)
	GetPaginationRangeByIndex(index int) (startIndex, endIndex int)
	GetPaginationRange() (startIndex, endIndex int)
	GetIndicator() (page, pageSize, total int)
	SetIndicator(page, pageSize, total int) error
	SetTotal(total int) error
	SetPageSize(pageSize int) error
}

type paginatorImpl struct {
	queries         paginationQueries
	defaultPageSize int
	page            int
	pageSize        int
	total           int
	firstPage       int
	lastPage        int
	prevPage        int
	nextPage        int
}

// Wrap is used for putting the input items to Result field of the Paginated struct.
func (p *paginatorImpl) Wrap(items interface{}) WrappedItems {
	return items
}

// WrapWithTruncate does the same thing with Wrap,
// and it truncates the input items by the pagination range.
// It may cause a panic if items is not Slice kind
func (p *paginatorImpl) WrapWithTruncate(items interface{}) WrappedItems {
	if kind := reflect.TypeOf(items).Kind(); kind != reflect.Slice {
		return items
	}

	startIndex, endIndex := p.GetPaginationRange()

	return reflect.ValueOf(items).Slice(startIndex, endIndex).Interface()
}

func (p *paginatorImpl) GetPaginationRangeByPage(page int) (startIndex, endIndex int) {
	pageSize := p.pageSize
	total := p.total
	offset := (page - 1) * pageSize

	if total > offset+pageSize {
		startIndex = offset
		endIndex = pageSize + offset
	} else if total-pageSize >= 0 {
		startIndex = total - pageSize
		endIndex = total
	}

	return
}

func (p *paginatorImpl) GetPaginationRangeByIndex(index int) (startIndex, endIndex int) {
	return p.GetPaginationRangeByPage((index / p.pageSize) + 1)
}

func (p *paginatorImpl) GetPaginationRange() (startIndex, endIndex int) {
	return p.GetPaginationRangeByPage(p.page)
}

func (p *paginatorImpl) GetIndicator() (page, pageSize, total int) {
	page = p.page
	pageSize = p.pageSize
	total = p.total
	return
}

func (p *paginatorImpl) SetIndicator(page, pageSize, total int) error {
	if page < 0 {
		return errors.New("page can't be a negative number")
	}
	if page == 0 {
		page = 1
	}
	if pageSize < 0 {
		return errors.New("pageSize can't be a negative number")
	}
	if pageSize == 0 {
		pageSize = p.defaultPageSize
	}
	if total < 0 {
		return errors.New("total can't be a negative number")
	}

	if total == 0 {
		p.page = page
		p.pageSize = pageSize
		p.firstPage = 1
		p.lastPage = 0
		p.total = 0
		p.prevPage = page - 1
		if p.prevPage < 1 {
			p.prevPage = 1
		}
		p.nextPage = page + 1
		return nil
	}

	p.page = page
	p.pageSize = pageSize
	p.total = total
	p.firstPage = 1
	p.lastPage = ((total-1)/pageSize + 1)
	if p.page > p.lastPage {
		p.page = p.lastPage
	}
	p.prevPage = page - 1
	if p.prevPage < p.firstPage {
		p.prevPage = p.firstPage
	}
	p.nextPage = page + 1
	if p.nextPage > p.lastPage {
		p.nextPage = p.lastPage
	}
	return nil
}

func (p *paginatorImpl) SetTotal(total int) error {
	if total < 0 {
		return errors.New("total can't be a negative number")
	}

	if total == 0 {
		p.lastPage = 0
		p.total = 0
		return nil
	}

	p.total = total
	p.lastPage = ((total-1)/p.pageSize + 1)
	if p.page > p.lastPage {
		p.page = p.lastPage
		p.prevPage = p.page - 1
		p.nextPage = p.page + 1
	}
	if p.prevPage < p.firstPage {
		p.prevPage = p.firstPage
	}
	if p.nextPage > p.lastPage {
		p.nextPage = p.lastPage
	}
	return nil
}

func (p *paginatorImpl) SetPageSize(pageSize int) error {
	if pageSize < 0 {
		return errors.New("pageSize can't be a negative number")
	}
	if pageSize == 0 {
		pageSize = p.defaultPageSize
	}

	p.pageSize = pageSize
	if p.total == 0 {
		return nil
	}
	p.lastPage = ((p.total-1)/pageSize + 1)
	if p.page > p.lastPage {
		p.page = p.lastPage
		p.prevPage = p.page - 1
		p.nextPage = p.page + 1
	}
	if p.prevPage < p.firstPage {
		p.prevPage = p.firstPage
	}
	if p.nextPage > p.lastPage {
		p.nextPage = p.lastPage
	}
	return nil
}
