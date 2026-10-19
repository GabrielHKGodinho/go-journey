package combinations

import (
	"fmt"
	"slices"
	"testing"
)

var implementations = []struct {
	name string
	fn   func(n, k int) [][]int
}{
	{name: "Combine", fn: Combine},
	{name: "CombineTakeSkip", fn: CombineTakeSkip},
}

func TestExamples(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		n, k int
		want [][]int
	}{
		{name: "leetcode example 1", n: 4, k: 2, want: [][]int{{1, 2}, {1, 3}, {1, 4}, {2, 3}, {2, 4}, {3, 4}}},
		{name: "leetcode example 2", n: 1, k: 1, want: [][]int{{1}}},
		{name: "k equals n", n: 3, k: 3, want: [][]int{{1, 2, 3}}},
		{name: "k equals 1", n: 3, k: 1, want: [][]int{{1}, {2}, {3}}},
	}

	for _, impl := range implementations {
		for _, tt := range tests {
			t.Run(impl.name+"/"+tt.name, func(t *testing.T) {
				t.Parallel()

				got := impl.fn(tt.n, tt.k)
				sortCombinations(got)
				if !equalCombinations(got, tt.want) {
					t.Errorf("%s(%d, %d) = %v, want %v", impl.name, tt.n, tt.k, got, tt.want)
				}
			})
		}
	}
}

// TestProperties checks every (n, k) up to n = 10 without a hand-written answer:
// the right count, each combination valid, and no combination repeated.
func TestProperties(t *testing.T) {
	t.Parallel()

	for _, impl := range implementations {
		for n := 1; n <= 10; n++ {
			for k := 1; k <= n; k++ {
				t.Run(fmt.Sprintf("%s/n=%d,k=%d", impl.name, n, k), func(t *testing.T) {
					t.Parallel()

					got := impl.fn(n, k)
					if want := binomial(n, k); len(got) != want {
						t.Fatalf("got %d combinations, want C(%d,%d) = %d", len(got), n, k, want)
					}

					seen := make(map[string]bool, len(got))
					for _, c := range got {
						if len(c) != k {
							t.Fatalf("combination %v has %d numbers, want %d", c, len(c), k)
						}
						for i, x := range c {
							if x < 1 || x > n {
								t.Fatalf("combination %v has %d outside 1..%d", c, x, n)
							}
							if i > 0 && c[i-1] >= x {
								t.Fatalf("combination %v is not strictly increasing", c)
							}
						}
						key := fmt.Sprint(c)
						if seen[key] {
							t.Fatalf("combination %v appears twice", c)
						}
						seen[key] = true
					}
				})
			}
		}
	}
}

// TestNoSharedBackingArray catches a missing slices.Clone: if two results
// shared one array, writing to one would change the other.
func TestNoSharedBackingArray(t *testing.T) {
	t.Parallel()

	for _, impl := range implementations {
		got := impl.fn(5, 3)
		for i := range got {
			for j := i + 1; j < len(got); j++ {
				if &got[i][0] == &got[j][0] {
					t.Fatalf("%s: results %d and %d share a backing array", impl.name, i, j)
				}
			}
		}
	}
}

// The pruning doesn't change the output, only the work, so benchmarks show it:
// go test -bench . -benchmem
func BenchmarkCombine(b *testing.B) {
	for _, impl := range implementations {
		for _, c := range []struct{ n, k int }{{20, 10}, {20, 18}} {
			b.Run(fmt.Sprintf("%s/n=%d,k=%d", impl.name, c.n, c.k), func(b *testing.B) {
				for b.Loop() {
					impl.fn(c.n, c.k)
				}
			})
		}
	}
}

func sortCombinations(cs [][]int) {
	slices.SortFunc(cs, slices.Compare[[]int])
}

func equalCombinations(a, b [][]int) bool {
	return slices.EqualFunc(a, b, slices.Equal[[]int])
}
