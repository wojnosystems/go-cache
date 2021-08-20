package cache

import (
	"context"
	"fmt"
	"github.com/wojnosystems/go-cache/capacity"
	"github.com/wojnosystems/go-cache/lru"
)

var ErrInsufficientCapacity = fmt.Errorf("insufficient capacity")

type lruItem struct {
	*unbounded
	maxLen  capacity.Tracker
	tracker lru.Tracker
}

// NewLRUItem is a cache that evicts the least recently used (oldest) item when a new item needs to
// be cached and there's insufficient space
func NewLRUItem(maxItems int, valueFactory ValueMapper) GetInvalidater {
	l := &lruItem{
		tracker: lru.NewTracker(),
	}
	l.maxLen = capacity.NewMaxLen(uint(maxItems), func() uint {
		return uint(l.tracker.Len())
	})
	l.unbounded = newUnbounded(func(ctx context.Context, key interface{}) (value interface{}, err error) {
		value, err = valueFactory(ctx, key)
		if err != nil {
			return
		}
		if l.maxLen.IsLargerThanCapacity(1) {
			return nil, ErrInsufficientCapacity
		}
		for !l.maxLen.HasCapacity(1) {
			leastRecentlyUsed, _ := l.tracker.LRU()
			l.Invalidate(leastRecentlyUsed)
		}
		return
	})
	return l
}

func (l *lruItem) Get(ctx context.Context, key interface{}) (value interface{}, err error) {
	value, err = l.unbounded.Get(ctx, key)
	if err != nil {
		return
	}
	l.tracker.Touch(key)
	return
}

func (l *lruItem) Invalidate(key interface{}) {
	l.unbounded.Invalidate(key)
	l.tracker.Remove(key)
}
