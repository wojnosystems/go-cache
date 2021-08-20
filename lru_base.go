package cache

import (
	"context"
	"github.com/wojnosystems/go-cache/capacity"
	"github.com/wojnosystems/go-cache/lru"
)

type valueMapSizer func(value interface{}) uint

type lruBase struct {
	*unbounded
	tracker    lru.Tracker
	limit      capacity.TrackMutator
	valueSizer valueMapSizer
}

func newLRUBase(cap uint, valueSizer valueMapSizer, valueMapper ValueMapper) *lruBase {
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
		valueSize := valueSizer(value)
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

func (l *lruBase) Invalidate(key interface{}) (ok bool) {
	var value interface{}
	if value, ok = l.unbounded.cache[key]; ok {
		l.unbounded.Invalidate(key)
		l.tracker.Remove(key)
		l.limit.Remove(l.valueSizer(value))
	}
	return
}
