package cache

import (
	"context"
	"github.com/wojnosystems/go-cache/capacity"
	"github.com/wojnosystems/go-cache/lru"
)

// ValueSizer should return the size of a value that is or to be stored in the cache
// the units don't matter, but need to match with what you specify in NewLRU as the cap variable
// for example: if cap represents the total number of items, this should return 1 for each item.
// for example: if cap represnets the total number of bytes, this should return the number of bytes for each item.
type ValueSizer func(value interface{}) uint

type lruBase struct {
	*unbounded
	tracker    lru.Tracker
	limit      capacity.TrackMutator
	valueSizer ValueSizer
}

// NewLRU creates a cache that has the ability to limit the size however you wish to track it
// cap: is the maximum "size" of this cache. The size is defined by you when you implement valueSizer
// valueSizer: Added items will use the size returned by valueSizer. Items removed will use the same
// valueMapper: looks up values based on keys
func NewLRU(cap uint, valueSizer ValueSizer, valueMapper ValueMapper) GetInvalidater {
	l := &lruBase{
		tracker:    lru.NewTracker(),
		limit:      capacity.NewMaxLen(cap),
		valueSizer: valueSizer,
	}
	l.unbounded = newUnbounded(func(ctx context.Context, key interface{}) (value interface{}, err error) {
		value, err = valueMapper(ctx, key)
		if err != nil {
			return
		}
		valueSize := l.valueSizer(value)
		if l.limit.IsLargerThanCapacity(valueSize) {
			return nil, ErrInsufficientCapacity
		}
		for !l.limit.Add(valueSize) {
			leastRecentlyUsedItem, _ := l.tracker.LRU()
			l.Invalidate(leastRecentlyUsedItem)
		}
		return
	})
	return l
}

func (l *lruBase) Get(ctx context.Context, key interface{}) (value interface{}, err error) {
	value, err = l.unbounded.Get(ctx, key)
	if err != nil {
		return
	}
	l.tracker.Touch(key)
	return
}

func (l *lruBase) Invalidate(key interface{}) {
	if value, ok := l.unbounded.cache[key]; ok {
		l.unbounded.Invalidate(key)
		l.tracker.Remove(key)
		l.limit.Remove(l.valueSizer(value))
	}
	return
}
