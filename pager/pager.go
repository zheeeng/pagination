package pager

import (
	"math"
)

// Pager provides basic calculations
// if total is greater than, page is restrict to a range between 0 and maxpage
type Pager struct {
	total    int
	page     int
	pageSize int
}

// Navigation defines pager infomation
type Navigation struct {
	Total    int
	Page     int
	PageSize int
	First    int
	Last     int
	Prev     int
	Next     int
}

func compact(min, max, value int) int {
	if min > max {
		min, max = max, min
	}
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func divCeil(a, b int) int {
	d := a / b
	if a%b > 0 {
		return d + 1
	}

	return d
}

//NewPager returns Pager instance
func NewPager(page, pageSize int) *Pager {
	return &Pager{0, compact(1, math.MaxInt32, page), compact(1, math.MaxInt32, pageSize)}
}

// getDefaultNavigation returns navigation info when missing total value
func (p *Pager) getDefaultNavigation() Navigation {
	return Navigation{
		Total:    0,
		Page:     p.page,
		PageSize: p.pageSize,
		First:    1,
		Last:     0,
		Prev:     compact(1, math.MaxInt32, p.page-1),
		Next:     compact(1, math.MaxInt32, p.page+1),
	}
}

// SetTotal sets total value to pager
func (p *Pager) SetTotal(total int) *Pager {
	p.total = total
	return p
}

// SetPageInfo resets page and pageSize to pager
func (p *Pager) SetPageInfo(page, pageSize int) *Pager {
	p.page = compact(1, math.MaxInt32, page)
	p.pageSize = compact(1, math.MaxInt32, pageSize)
	return p
}

// ClonePager returns a fresh pager with specified page and pageSize
func (p *Pager) ClonePager(page, pageSize int) *Pager {
	return &Pager{p.total, compact(1, math.MaxInt32, page), compact(1, math.MaxInt32, pageSize)}
}

// ClonePagerWithCursor returns a fresh pager with specified cursor value and pageSize
func (p *Pager) ClonePagerWithCursor(cursor, pageSize int) *Pager {
	return p.ClonePager(divCeil(cursor, pageSize), pageSize)
}

// GetNavigation returns navigation info
func (p *Pager) GetNavigation() Navigation {
	if p.total <= 0 {
		return p.getDefaultNavigation()
	}

	last := divCeil(p.total, p.pageSize)

	return Navigation{
		Total:    p.total,
		Page:     p.page,
		PageSize: p.pageSize,
		First:    1,
		Last:     last,
		Prev:     compact(1, last, p.page-1),
		Next:     compact(1, last, p.page+1),
	}
}

// GetRange returns the start and end offset values
func (p *Pager) GetRange() (start, end int) {
	offset, length := p.GetOffsetRange()
	start = offset
	end = offset + length

	return
}

// GetOffsetRange returns start and end offsets of items
func (p *Pager) GetOffsetRange() (offset, length int) {
	offset = (p.page - 1) * p.pageSize
	length = p.pageSize

	if p.total > 0 {
		offset = compact(0, p.total, offset)
		length = compact(0, p.total-offset, length)
	}

	return
}
