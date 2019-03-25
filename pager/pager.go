package pager

// Pager provides basic calculations
// if total is greater than, page is restrict to a range between 0 and maxpage
type Pager struct {
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

// getDefaultNavigation returns navigation info when missing total value
func (p *Pager) getDefaultNavigation() Navigation {
	return Navigation{
		Total:    0,
		Page:     p.page,
		PageSize: p.pageSize,
		First:    0,
		Last:     0,
		Prev:     0,
		Next:     0,
	}
}

// GetNavigation returns navigation info
func (p *Pager) GetNavigation(total int) Navigation {
	if total <= 0 {
		return p.getDefaultNavigation()
	}

	last := divCeil(total, p.pageSize)

	return Navigation{
		Total:    total,
		Page:     p.page,
		PageSize: p.pageSize,
		First:    1,
		Last:     last,
		Prev:     compact(1, last, p.page-1),
		Next:     compact(1, last, p.page+1),
	}
}

func (p *Pager) getRange(total int) (start, end int) {
	offset, length := p.GetOffsetRange(total)
	start = offset
	end = offset + length

	return
}

// GetOffsetRange returns start and end offsets of items
func (p *Pager) GetOffsetRange(total int) (offset, length int) {
	offset = (p.page - 1) * p.pageSize
	length = p.pageSize

	if total > 0 {
		offset = compact(0, total, offset)
		length = compact(0, total-offset, length)
	}

	return
}
