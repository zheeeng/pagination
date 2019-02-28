package pagination

import (
	"errors"
)

// Truncatable is a nominal type for Pagination::Wrap func consuming.
// It ensure user always call Paginator::Wrap or Paginator::WrapWithTruncate to return items
type Truncatable interface {
	Slice(startIndex, endIndex int) Truncatable
}

var (
	// ErrNegativePage -- page can't be a negative number
	ErrNegativePage = errors.New("page can't be a negative number")
	// ErrNegativePageSize -- pageSize can't be a negative number
	ErrNegativePageSize = errors.New("pageSize can't be a negative number")
	// ErrNegativeTotal -- total can't be a negative number
	ErrNegativeTotal = errors.New("total can't be a negative number")
)

// Paginator provides methods to manipulate pagination fields
type Paginator struct {
	queries         paginationQueries
	DefaultPageSize int
	page            int
	pageSize        int
	total           int
	firstPage       int
	lastPage        int
	prevPage        int
	nextPage        int
	hasPage         bool
	hasPageSize     bool
}

// Wrap is used for putting the input items to Result field of the Paginated struct.
func (p *Paginator) Wrap(items Truncatable) Truncatable {
	return items
}

// WrapWithTruncate does the same thing with Wrap,
// and it truncates the input items by the pagination range.
// It may cause a panic if items is not Slice kind
func (p *Paginator) WrapWithTruncate(items Truncatable) Truncatable {
	startIndex, endIndex := p.GetPaginationRange()

	return items.Slice(startIndex, endIndex)
}

// GetPaginationRangeByPage returns the corresponding start and end indices by a specific page number
func (p *Paginator) GetPaginationRangeByPage(page int) (startIndex, endIndex int) {
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

// GetPaginationRangeByIndex returns the corresponding start and end indices by a specific item index number
func (p *Paginator) GetPaginationRangeByIndex(index int) (startIndex, endIndex int) {
	return p.GetPaginationRangeByPage((index / p.pageSize) + 1)
}

// GetPaginationRange returns the corresponding start and end indices by Paginator context
func (p *Paginator) GetPaginationRange() (startIndex, endIndex int) {
	return p.GetPaginationRangeByPage(p.page)
}

// GetIndicator returns current page, pageSize, total in its context
func (p *Paginator) GetIndicator() (page, pageSize, total int) {
	page = p.page
	pageSize = p.pageSize
	total = p.total
	return
}

// HasRawPagination returns whether the test link contains pagination fields
func (p *Paginator) HasRawPagination() bool {
	// if link contains page, coz we assigned default page_size
	// or if link contains page_size, we assigned default page 1,
	// we think there has a specificed pagination infomation.
	return p.hasPage || p.hasPageSize
}

// HasRawPage returns whether the test link contains 'page' field
func (p *Paginator) HasRawPage() bool {
	return p.hasPage
}

// HasRawPageSize returns whether the test link contains 'page_size' field
func (p *Paginator) HasRawPageSize() bool {
	return p.hasPageSize
}

// SetIndicator sets current page, pageSize, total, it checks the inputs validation
func (p *Paginator) SetIndicator(page, pageSize, total int) error {
	if page < 0 {
		return ErrNegativePage
	}
	if pageSize < 0 {
		return ErrNegativePageSize
	}
	if total < 0 {
		return ErrNegativeTotal
	}

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = p.DefaultPageSize
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

// SetTotal tells Paginator the total number of items
func (p *Paginator) SetTotal(total int) error {
	if total < 0 {
		return ErrNegativeTotal
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

// SetPage sets the page of paginator,
// by default this value have be parsed from link's query fields
func (p *Paginator) SetPage(page int) error {
	if page < 0 {
		return ErrNegativePage
	}
	if page == 0 {
		page = 1
	}

	p.page = page

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

// SetPageSize sets the pageSize of paginator,
// by default this value have be parsed from link's query fields
func (p *Paginator) SetPageSize(pageSize int) error {
	if pageSize < 0 {
		return errors.New("pageSize can't be a negative number")
	}
	if pageSize == 0 {
		pageSize = p.DefaultPageSize
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
