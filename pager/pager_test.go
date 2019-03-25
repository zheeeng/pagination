package pager

import (
	"fmt"
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
		caseName   string
		pager      Pager
		total      int
		start      int
		end        int
		navigation Navigation
	}{
		{
			"default(total is zero value)",
			Pager{5, 10},
			0, 40, 50,
			Navigation{0, 5, 10, 0, 0, 0, 0},
		},
		{
			"default(total value is below 0)",
			Pager{5, 10},
			-1, 40, 50,
			Navigation{0, 5, 10, 0, 0, 0, 0},
		},
		{
			"basic",
			Pager{5, 10},
			100, 40, 50,
			Navigation{100, 5, 10, 1, 10, 4, 6},
		},
		{
			"total is on upper bound",
			Pager{5, 10},
			50, 40, 50,
			Navigation{50, 5, 10, 1, 5, 4, 5},
		},
		{
			"total is below to upper bound",
			Pager{5, 10},
			49, 40, 49,
			Navigation{49, 5, 10, 1, 5, 4, 5},
		},
	}

	for i, test := range tests {
		desc := fmt.Sprintf("[%d]: test [%s] functionality\n", i, test.caseName)

		navigation := test.pager.GetNavigation(test.total)
		start, end := test.pager.getRange(test.total)

		if navigation != test.navigation {
			t.Errorf("%s[navigation] output doesn't match expected, `pager` is %v, `total` is %d, expects `navigation` is %v, got %v",
				desc, test.pager, test.total, test.navigation, navigation,
			)
		}

		if start != test.start || end != test.end {
			t.Errorf("%s[start, end offsets] output doesn't match expected, `pager` is %v, `total` is %d, expects (`start`, `end`) is (%d, %d), got (%d, %d)",
				desc, test.pager, test.total, test.start, test.end, start, end,
			)
		}
	}

}
