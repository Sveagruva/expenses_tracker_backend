package utils

import (
	"sync"
)

type IncreasingCounter struct {
	mu     sync.Mutex
	number int
}

func (c *IncreasingCounter) Next() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.number++
	return c.number
}
