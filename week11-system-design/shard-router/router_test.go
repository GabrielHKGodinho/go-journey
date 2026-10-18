package shardrouter

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

// userIDs returns n deterministic, distinct user IDs.
func userIDs(n int) []string {
	ids := make([]string, n)
	for i := range ids {
		ids[i] = fmt.Sprintf("user-%d", i)
	}
	return ids
}

// userInShard returns some user ID whose logical shard is target.
func userInShard(t *testing.T, target int) string {
	t.Helper()
	for i := 0; i < 1_000_000; i++ {
		id := fmt.Sprintf("probe-%d", i)
		if LogicalShard(id) == target {
			return id
		}
	}
	t.Fatalf("no user found for logical shard %d", target)
	return ""
}

func countPerDatabase(r *Router) map[string]int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	counts := make(map[string]int)
	for _, name := range r.placement {
		counts[name]++
	}
	return counts
}

func TestLogicalShard_DeterministicAndInRange(t *testing.T) {
	t.Parallel()
	for _, id := range userIDs(10_000) {
		got := LogicalShard(id)
		if got < 0 || got >= LogicalShards {
			t.Fatalf("LogicalShard(%q) = %d, want 0..%d", id, got, LogicalShards-1)
		}
		if again := LogicalShard(id); again != got {
			t.Fatalf("LogicalShard(%q) not deterministic: %d then %d", id, got, again)
		}
	}
}

// Golden values pin the hash function: routing must give the same answer in
// every process and every release, so a seeded hash (hash/maphash, map hashing)
// or a changed LogicalShards breaks this test.
func TestLogicalShard_StableAcrossProcesses(t *testing.T) {
	t.Parallel()
	golden := map[string]int{
		"user-0":                               39,
		"user-42":                              459,
		"3f6c1a2e-9b7d-4c1e-8a55-0d2f7b9e1c44": 267,
	}
	for id, want := range golden {
		if got := LogicalShard(id); got != want {
			t.Errorf("LogicalShard(%q) = %d, want %d", id, got, want)
		}
	}
}

func TestLogicalShard_SpreadsUsers(t *testing.T) {
	t.Parallel()
	const users = 480_000 // 1,000 per logical shard on average
	counts := make([]int, LogicalShards)
	for _, id := range userIDs(users) {
		counts[LogicalShard(id)]++
	}
	for logical, c := range counts {
		// Generous bounds: catch "everything in a few shards", not normal noise.
		if c < 800 || c > 1200 {
			t.Errorf("logical shard %d got %d users, want 800..1200", logical, c)
		}
	}
}

func TestNew_SpreadsLogicalShardsEvenly(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		physical []string
		want     int // logical shards per database
	}{
		{"one database", []string{"pg-1"}, 480},
		{"four databases", []string{"pg-1", "pg-2", "pg-3", "pg-4"}, 120},
		{"six databases", []string{"pg-1", "pg-2", "pg-3", "pg-4", "pg-5", "pg-6"}, 80},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, err := New(tt.physical)
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			counts := countPerDatabase(r)
			if len(counts) != len(tt.physical) {
				t.Fatalf("got %d databases in use, want %d", len(counts), len(tt.physical))
			}
			for _, name := range tt.physical {
				if counts[name] != tt.want {
					t.Errorf("%s holds %d logical shards, want %d", name, counts[name], tt.want)
				}
			}
		})
	}
}

func TestNew_RejectsInvalidInput(t *testing.T) {
	t.Parallel()
	tooMany := make([]string, LogicalShards+1)
	for i := range tooMany {
		tooMany[i] = fmt.Sprintf("pg-%d", i)
	}
	tests := []struct {
		name     string
		physical []string
		wantErr  error
	}{
		{"nil", nil, ErrNoDatabases},
		{"empty", []string{}, ErrNoDatabases},
		{"more databases than logical shards", tooMany, ErrTooManyDatabases},
		{"empty name", []string{"pg-1", ""}, ErrEmptyDatabaseName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, err := New(tt.physical)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("New error = %v, want %v", err, tt.wantErr)
			}
			if r != nil {
				t.Fatalf("New returned a router together with an error")
			}
		})
	}
}

func TestMove_ChangesOnlyThatLogicalShard(t *testing.T) {
	t.Parallel()
	r, err := New([]string{"pg-1", "pg-2", "pg-3", "pg-4"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	moved := userInShard(t, 5)
	untouched := userInShard(t, 6)
	before := r.PhysicalFor(untouched)

	if err := r.Move(5, "pg-5"); err != nil {
		t.Fatalf("Move: %v", err)
	}
	if got := r.PhysicalFor(moved); got != "pg-5" {
		t.Errorf("user in moved shard is on %q, want pg-5", got)
	}
	if got := r.PhysicalFor(untouched); got != before {
		t.Errorf("user in another shard moved from %q to %q", before, got)
	}
}

func TestMove_RejectsInvalidInput(t *testing.T) {
	t.Parallel()
	r, err := New([]string{"pg-1"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	tests := []struct {
		name     string
		logical  int
		physical string
		wantErr  error
	}{
		{"negative shard", -1, "pg-2", ErrLogicalOutOfRange},
		{"shard equal to LogicalShards", LogicalShards, "pg-2", ErrLogicalOutOfRange},
		{"empty database name", 0, "", ErrEmptyDatabaseName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Move(tt.logical, tt.physical); !errors.Is(err, tt.wantErr) {
				t.Fatalf("Move error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// Adding a fifth database by moving 24 logical shards from each of the four
// existing ones must move about 1/5 of the users, not the ~80% of hash % N.
func TestResharding_FourToFiveMovesOneFifth(t *testing.T) {
	t.Parallel()
	r, err := New([]string{"pg-1", "pg-2", "pg-3", "pg-4"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	users := userIDs(100_000)
	before := make([]string, len(users))
	for i, id := range users {
		before[i] = r.PhysicalFor(id)
	}

	// Each database owns a block of 120; hand the last 24 of each block to pg-5.
	for block := 0; block < 4; block++ {
		for logical := block*120 + 96; logical < (block+1)*120; logical++ {
			if err := r.Move(logical, "pg-5"); err != nil {
				t.Fatalf("Move(%d): %v", logical, err)
			}
		}
	}

	moved := 0
	for i, id := range users {
		if r.PhysicalFor(id) != before[i] {
			moved++
		}
	}
	fraction := float64(moved) / float64(len(users))
	if fraction < 0.18 || fraction > 0.22 {
		t.Fatalf("moved %.1f%% of users, want about 20%%", fraction*100)
	}
	if got := countPerDatabase(r)["pg-5"]; got != 96 {
		t.Fatalf("pg-5 holds %d logical shards, want 96", got)
	}
}

// Run with -race: lookups must be safe while a migration is moving shards.
func TestRouter_ConcurrentLookupsAndMoves(t *testing.T) {
	t.Parallel()
	r, err := New([]string{"pg-1", "pg-2"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, id := range userIDs(2_000) {
				_ = r.PhysicalFor(id)
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for logical := 0; logical < LogicalShards; logical++ {
			if err := r.Move(logical, "pg-3"); err != nil {
				t.Errorf("Move(%d): %v", logical, err)
				return
			}
		}
	}()
	wg.Wait()
}
