package pagination

import (
	"strconv"

	"github.com/zheeeng/pagination/pager"
)

// Truncatable is used for feeding Paginator::Wrap and Paginator::WrapWithTruncate, to wrap items into paginated result
type Truncatable interface {
	Len() int
	Slice(startIndex, endIndex int) Truncatable
}

// Indicator contains pagination information, e.g. page, pageSize, total, first, last...
type Indicator = pager.Navigation

// Paginator provides methods to manipulate pagination fields
type Paginator struct {
	pager           *pager.Pager
	basePath        string
	queries         paginationQueries
	DefaultPageSize int
	hasPage         bool
	hasPageSize     bool
}

func (p *Paginator) buildFields() *PageFields {
	nav := p.pager.GetNavigation()
	fields := &PageFields{
		Page:     nav.Page,
		PageSize: nav.PageSize,
		Total:    nav.Total,
		Query:    p.queries.query,
	}

	p.queries.query.Set("page", strconv.Itoa(nav.Page))
	p.queries.query.Set("page_size", strconv.Itoa(nav.PageSize))

	p.queries.firstQuery.Set("page", strconv.Itoa(nav.First))
	p.queries.firstQuery.Set("page_size", strconv.Itoa(nav.PageSize))
	fields.First = p.basePath + "?" + p.queries.firstQuery.Encode()

	if nav.Last > 0 {
		p.queries.lastQuery.Set("page", strconv.Itoa(nav.Last))
		p.queries.lastQuery.Set("page_size", strconv.Itoa(nav.PageSize))
		fields.Last = p.basePath + "?" + p.queries.lastQuery.Encode()
	}

	p.queries.prevQuery.Set("page", strconv.Itoa(nav.Prev))
	p.queries.prevQuery.Set("page_size", strconv.Itoa(nav.PageSize))
	fields.Prev = p.basePath + "?" + p.queries.prevQuery.Encode()

	p.queries.nextQuery.Set("page", strconv.Itoa(nav.Next))
	p.queries.nextQuery.Set("page_size", strconv.Itoa(nav.PageSize))
	fields.Next = p.basePath + "?" + p.queries.nextQuery.Encode()

	return fields
}

// Wrap is used for putting the input items to Result field of the Paginated struct.
func (p *Paginator) Wrap(items Truncatable, total int) Paginated {
	p.pager.SetTotal(total)
	fields := p.buildFields()

	return Paginated{
		Pagination: fields,
		Result:     items,
	}
}

// WrapWithTruncate does the same thing with Wrap,
// and it truncates the input items by the pagination range.
// It may cause a panic if items is not Slice kind
func (p *Paginator) WrapWithTruncate(items Truncatable, total int) Paginated {
	p.pager.SetTotal(total)
	fields := p.buildFields()

	length := items.Len()

	startIndex, endIndex := p.GetRange()

	if endIndex > length {
		endIndex = length
	}

	return Paginated{
		Pagination: fields,
		Result:     items.Slice(startIndex, endIndex),
	}
}

// GetRangeByIndex returns the corresponding start and end offsets by a specific item index number
func (p *Paginator) GetRangeByIndex(index int) (start, end int) {
	return p.pager.ClonePagerWithCursor(index, p.pager.GetNavigation().PageSize).GetRange()
}

// GetRange returns the corresponding start and end offsets by Paginator context
func (p *Paginator) GetRange() (start, end int) {
	return p.pager.GetRange()
}

// GetOffsetRangeByIndex returns the corresponding offset and range length by a specific item index number
func (p *Paginator) GetOffsetRangeByIndex(index int) (offset, length int) {
	return p.pager.ClonePagerWithCursor(index, p.pager.GetNavigation().PageSize).GetOffsetRange()
}

// GetOffsetRange returns the corresponding offset and range length by Paginator context
func (p *Paginator) GetOffsetRange() (offset, length int) {
	return p.pager.GetOffsetRange()
}

// GetIndicator returns current page, pageSize, total and tother info in its context
func (p *Paginator) GetIndicator() Indicator {
	return p.pager.GetNavigation()
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
