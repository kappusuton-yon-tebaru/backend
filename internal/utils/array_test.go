package utils_test

import (
	"testing"

	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	tcs := []struct {
		name     string
		inputs   []int
		expected []int
	}{
		{
			"no filtered",
			[]int{1, 2, 3, 4, 5},
			[]int{},
		},
		{
			"partial filtered",
			[]int{4, 5, 6, 7, 8},
			[]int{6, 7, 8},
		},
		{
			"full filtered",
			[]int{6, 7, 8, 9, 10},
			[]int{6, 7, 8, 9, 10},
		},
		{
			"no inputs",
			[]int{},
			[]int{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			filtered := utils.Filter(tc.inputs, func(e int) bool {
				return e > 5
			})

			assert.Equal(t, filtered, tc.expected)
		})
	}
}

func TestPaginate(t *testing.T) {
	tcs := []struct {
		name     string
		page     int
		limit    int
		inputs   []int
		expected []int
	}{
		{
			"first page in range limit",
			1,
			3,
			[]int{1, 2, 3, 4, 5, 6, 7},
			[]int{1, 2, 3},
		},
		{
			"second page in range limit",
			2,
			3,
			[]int{1, 2, 3, 4, 5, 6, 7},
			[]int{4, 5, 6},
		},
		{
			"third page partially in range limit",
			3,
			3,
			[]int{1, 2, 3, 4, 5, 6, 7},
			[]int{7},
		},
		{
			"bigger limit",
			1,
			10,
			[]int{1, 2, 3, 4, 5, 6, 7},
			[]int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			"non-positive page",
			0,
			10,
			[]int{1, 2, 3, 4, 5, 6, 7},
			[]int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			"overflow page",
			2,
			10,
			[]int{1, 2, 3, 4, 5, 6, 7},
			[]int{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			filtered := utils.Paginate(tc.inputs, tc.page, tc.limit)
			assert.Equal(t, filtered, tc.expected)
		})
	}
}
