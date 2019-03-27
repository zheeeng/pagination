package pager

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestCompact(t *testing.T) {
	tests := []struct {
		min    int
		max    int
		value  int
		expect int
	}{
		{1, 10, 5, 5},
		{1, 10, 0, 1},
		{1, 10, 1, 1},
		{1, 10, 10, 10},
		{1, 10, 11, 10},
		// swap min max
		{10, 1, 5, 5},
		{10, 1, 0, 1},
		{10, 1, 1, 1},
		{10, 1, 10, 10},
		{10, 1, 11, 10},
	}

	for _, test := range tests {
		if compact(test.min, test.max, test.value) != test.expect {
			t.Errorf("Error happens, boundary1: %d, boundary2: %d, value: %d, expect: %d",
				test.min, test.max, test.value, test.expect,
			)
		}
	}
}

func TestDivCeil(t *testing.T) {
	tests := []struct {
		a      int
		b      int
		expect int
	}{
		{51, 10, 6},
		{50, 10, 5},
		{49, 10, 5},
	}

	for _, test := range tests {
		value := divCeil(test.a, test.b)
		if value != test.expect {
			t.Errorf("Error happens, input: %d, %d, expect get: %d, actually got: %d",
				test.a, test.b, test.expect, value,
			)
		}
	}
}

func TestPager(t *testing.T) {
	tests := []struct {
		caseName              string
		page, pageSize, total int
		start, end            int
		navigation            Navigation
	}{
		{
			"total is zero value",
			5, 10, 0,
			40, 50,
			Navigation{0, 5, 10, 1, 0, 4, 6},
		},
		{
			"total value is below 0",
			5, 10, -1,
			40, 50,
			Navigation{0, 5, 10, 1, 0, 4, 6},
		},
		{
			"total is zero value, page in zero value",
			0, 10, 0,
			0, 10,
			Navigation{0, 1, 10, 1, 0, 1, 2},
		},
		{
			"total is zero value, page is the lowest bound value",
			1, 10, 0,
			0, 10,
			Navigation{0, 1, 10, 1, 0, 1, 2},
		},
		{
			"basic",
			5, 10, 100,
			40, 50,
			Navigation{100, 5, 10, 1, 10, 4, 6},
		},
		{
			"total is on upper bound",
			5, 10, 50,
			40, 50,
			Navigation{50, 5, 10, 1, 5, 4, 5},
		},
		{
			"total is below to upper bound",
			5, 10, 49,
			40, 49,
			Navigation{49, 5, 10, 1, 5, 4, 5},
		},
	}

	for i, test := range tests {
		desc := fmt.Sprintf("[%d]: test [%s] functionality\n", i, test.caseName)

		pager := NewPager(test.page, test.pageSize)
		pager.SetTotal(test.total)

		navigation := pager.GetNavigation()
		start, end := pager.GetRange()

		if navigation != test.navigation {
			t.Errorf("%s[navigation] output doesn't match expected, `pager` is %v, `total` is %d, expects `navigation` is %v, got %v",
				desc, pager, test.total, test.navigation, navigation,
			)
		}

		if start != test.start || end != test.end {
			t.Errorf("%s[start, end offsets] output doesn't match expected, `pager` is %v, `total` is %d, expects (`start`, `end`) is (%d, %d), got (%d, %d)",
				desc, pager, test.total, test.start, test.end, start, end,
			)
		}
	}

}

func TestSetPageInfoAndClonePager(t *testing.T) {
	tests := []struct {
		fromPage, fromPageSize, fromTotal int
		toPage, toPageSize, toTotal       int
	}{
		{5, 10, 0, 5, 10, -1},
		{5, 10, -1, 0, 10, 0},
		{0, 10, 0, 1, 10, 0},
		{1, 10, 0, 5, 10, 100},
		{5, 10, 100, 5, 10, 50},
		{5, 10, 50, 5, 10, 49},
	}

	for i, test := range tests {
		pager := NewPager(test.fromPage, test.fromPageSize)
		pager.SetTotal(test.fromTotal)

		toComparePager := NewPager(test.toPage, test.toPageSize)

		clonedPager := pager.ClonePager(test.toPage, test.toPageSize)
		clonedPager2 := pager.ClonePagerWithCursor(
			test.toPage*test.toPageSize-rand.Intn(test.toPageSize),
			test.toPageSize,
		)

		pager.SetPageInfo(test.toPage, test.toPageSize)

		if clonedPager.total != pager.total ||
			clonedPager.page != toComparePager.page ||
			clonedPager.pageSize != toComparePager.pageSize {
			t.Errorf("%d. [Clone failed], got clonedPager: %v,", i, clonedPager)
		}
		if clonedPager2.total != pager.total ||
			clonedPager2.page != toComparePager.page ||
			clonedPager2.pageSize != toComparePager.pageSize {
			t.Errorf("%d .[Clone from cursor failed], got colonedPager2: %v", i, clonedPager2)
		}
		if pager.page != toComparePager.page ||
			pager.pageSize != toComparePager.pageSize {
			t.Errorf("%d .[SetPageInfo failed], got resetted page: %v,", i, pager)

		}
	}

}
