// Package shardrouter maps users to databases using a fixed number of
// logical shards, so adding a physical database moves only the logical
// shards it takes over instead of reshuffling almost every key.
package shardrouter

import (
	"errors"
	"fmt"
	"hash/fnv"
	"sync"
)

// LogicalShards is fixed forever: changing it remaps almost every user.
// 480 = 2^5 * 3 * 5 divides evenly by 2, 3, 4, 5, 6, 8, 10, 12, 15, 16...
const LogicalShards = 480

var (
	ErrNoDatabases       = errors.New("at least one physical database is required")
	ErrTooManyDatabases  = errors.New("more physical databases than logical shards")
	ErrEmptyDatabaseName = errors.New("physical database name must not be empty")
	ErrLogicalOutOfRange = errors.New("logical shard out of range")
)

// Router resolves users to physical databases. It is safe for concurrent use:
// requests call PhysicalFor while a migration calls Move.
type Router struct {
	mu        sync.RWMutex
	placement []string // index = logical shard, value = physical database name
}

// New spreads the logical shards across the physical databases in contiguous,
// equally sized blocks (as equal as the division allows).
func New(physical []string) (*Router, error) {
	if len(physical) == 0 {
		return nil, ErrNoDatabases
	}
	if len(physical) > LogicalShards {
		return nil, ErrTooManyDatabases
	}
	for _, name := range physical {
		if name == "" {
			return nil, ErrEmptyDatabaseName
		}
	}

	placement := make([]string, LogicalShards)
	for logical := range placement {
		placement[logical] = physical[logical*len(physical)/LogicalShards]
	}
	return &Router{placement: placement}, nil
}

// LogicalShard returns the logical shard that owns userID. It depends only on
// the user and LogicalShards, never on how many physical databases exist.
func LogicalShard(userID string) int {
	h := fnv.New32a()
	h.Write([]byte(userID)) // hash.Hash.Write never returns an error
	return int(h.Sum32() % LogicalShards)
}

// PhysicalFor returns the physical database that currently holds userID's data.
func (r *Router) PhysicalFor(userID string) string {
	logical := LogicalShard(userID) // pure computation: no lock needed

	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.placement[logical]
}

// Move reassigns one logical shard to a physical database. It is the last
// step of a migration, after the shard's data was copied and caught up.
func (r *Router) Move(logical int, physical string) error {
	if logical < 0 || logical >= LogicalShards {
		return fmt.Errorf("move logical shard %d: %w", logical, ErrLogicalOutOfRange)
	}
	if physical == "" {
		return fmt.Errorf("move logical shard %d: %w", logical, ErrEmptyDatabaseName)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.placement[logical] = physical
	return nil
}
