package cache

import (
	"strconv"
	"sync"
	"testing"
)

func TestCacheConcurrentAccess(t *testing.T) {
	c := NewCache()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := strconv.Itoa(n)
			c.Set(key, "value-"+key)
			c.Get(key)
		}(i)
	}

	wg.Wait()
}
