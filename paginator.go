package pagination

import (
	"errors"
	"net/url"
)

// Paginator provides methods to manipulate pagination fields
type Paginator interface {
	Wrap(items interface{}) interface{}
	WrapWithTruncate(func(startIndex, endIndex int) (items interface{})) interface{}
	GetIndicator() (page, pageSize, total int)
	SetIndicator(page, pageSize, total int) error
	SetTotal(total int) error
	SetPageSize(pageSize int) error
}

type paginatorImpl struct {
	Query           url.Values
	FirstQuery      url.Values
	LastQuery       url.Values
	PreviousQuery   url.Values
	NextQuery       url.Values
	defaultPageSize int
	page            int
	pageSize        int
	total           int
	firstPage       int
	lastPage        int
	previousPage    int
	nextPage        int
}

func (p *paginatorImpl) Wrap(items interface{}) interface{} {
	return items
}

func (p *paginatorImpl) WrapWithTruncate(truncate func(startIndex, endIndex int) (items interface{})) interface{} {
	page := p.page
	pageSize := p.pageSize
	total := p.total
	offset := (page - 1) * pageSize

	startIndex := 0
	endIndex := 0

	if total > offset+pageSize {
		startIndex = offset
		endIndex = pageSize + offset
	} else if total-pageSize >= 0 {
		startIndex = total - pageSize
		endIndex = total
	}

	return truncate(startIndex, endIndex)
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
		p.previousPage = page - 1
		if p.previousPage < 1 {
			p.previousPage = 1
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
	p.previousPage = page - 1
	if p.previousPage < p.firstPage {
		p.previousPage = p.firstPage
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
		p.previousPage = p.page - 1
		p.nextPage = p.page + 1
	}
	if p.previousPage < p.firstPage {
		p.previousPage = p.firstPage
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
		p.previousPage = p.page - 1
		p.nextPage = p.page + 1
	}
	if p.previousPage < p.firstPage {
		p.previousPage = p.firstPage
	}
	if p.nextPage > p.lastPage {
		p.nextPage = p.lastPage
	}
	return nil
}