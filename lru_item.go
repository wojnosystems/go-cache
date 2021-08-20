package cache

import (
	"fmt"
)

var ErrInsufficientCapacity = fmt.Errorf("insufficient capacity")

// NewLRUItem is a cache that evicts the least recently used (oldest) item when a new item needs to
// be cached and there's insufficient space
func NewLRUItem(maxItems int, valueMapper ValueMapper) GetInvalidater {
	return NewLRU(uint(maxItems), func(value interface{}) uint {
		return uint(1)
	}, valueMapper)
}
