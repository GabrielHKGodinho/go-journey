package besttime

import "testing"

func TestMaxProfit(t *testing.T) {
	tests := []struct {
		name   string
		prices []int
		want   int
	}{
		{"leetcode example 1", []int{7, 1, 5, 3, 6, 4}, 5},
		{"leetcode example 2: only falling", []int{7, 6, 4, 3, 1}, 0},
		{"single day", []int{0}, 0},
		{"new minimum after the best sell", []int{2, 5, 1, 8}, 7},
		{"new minimum but the old trade stays best", []int{2, 9, 1, 3}, 7},
		{"flat prices", []int{3, 3, 3}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := maxProfit(tt.prices); got != tt.want {
				t.Errorf("maxProfit(%v) = %d, want %d", tt.prices, got, tt.want)
			}
		})
	}
}