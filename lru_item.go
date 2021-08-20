package cache

import (
	"context"
	"fmt"
	"github.com/wojnosystems/go-cache/lru"
)

var ErrInsufficientCapacity = fmt.Errorf("insufficient capacity")

type lruItem struct {
	*unbounded
	maxItems int
	tracker  lru.Tracker
}

// NewLRUItem is a cache that evicts the least recently used (oldest) item when a new item needs to
// be cached and there's insufficient space
func NewLRUItem(maxItems int, valueFactory ValueMapper) GetInvalidater {
	l := &lruItem{
		maxItems: maxItems,
		tracker:  lru.NewTracker(),
	}
	l.unbounded = newUnbounded(func(ctx context.Context, key interface{}) (value interface{}, err error) {
		value, err = valueFactory(ctx, key)
		if err != nil {
			return
		}
		if !l.canFitInTotal() {
			return nil, ErrInsufficientCapacity
		}
		for !l.canFitInRemaining() {
			leastRecentlyUsed, _ := l.tracker.LRU()
			l.Invalidate(leastRecentlyUsed)
		}
		return
	})
	return l
}

func (l *lruItem) canFitInRemaining() bool {
	return l.tracker.Len() < l.maxItems
}

func (l *lruItem) canFitInTotal() bool {
	return l.maxItems > 0
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
