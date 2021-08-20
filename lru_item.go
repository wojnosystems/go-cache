package cache

import (
	"container/list"
	"context"
	"fmt"
)

var ErrInsufficientCapacity = fmt.Errorf("insufficient capacity")

// recencyIndexType maps keys to recency elements
type recencyIndexType map[interface{}]*list.Element

type lruItem struct {
	*unbounded
	maxItems     int
	recency      *list.List // of interface{} keys, front = oldest, back = LRU
	recencyIndex recencyIndexType
}

// NewLRUItem is a cache that evicts the least recently used (oldest) item when a new item needs to
// be cached and there's insufficient space
func NewLRUItem(maxItems int, valueFactory ValueMapper) GetInvalidater {
	l := &lruItem{
		maxItems:     maxItems,
		recency:      list.New(),
		recencyIndex: make(recencyIndexType),
	}
	l.unbounded = newUnbounded(func(ctx context.Context, key interface{}) (value interface{}, err error) {
		value, err = valueFactory(ctx, key)
		if err != nil {
			return
		}
		if l.isAtOrAboveCapacity() {
			l.removeLRU()
		}
		if l.isAtOrAboveCapacity() {
			return nil, ErrInsufficientCapacity
		}
		return
	})
	return l
}

func (l *lruItem) isAtOrAboveCapacity() bool {
	return len(l.recencyIndex) >= l.maxItems
}

func (l *lruItem) removeLRU() {
	if oldest := l.recency.Front(); oldest != nil {
		l.Invalidate(oldest.Value)
	}
}

func (l *lruItem) Get(ctx context.Context, key interface{}) (value interface{}, err error) {
	value, err = l.unbounded.Get(ctx, key)
	if err != nil {
		return
	}
	l.touch(key)
	return
}

func (l *lruItem) touch(key interface{}) {
	vm, ok := l.recencyIndex[key]
	if ok {
		l.recency.Remove(vm)
	}
	l.recencyIndex[key] = l.recency.PushBack(key)
}

func (l *lruItem) Invalidate(key interface{}) {
	l.unbounded.Invalidate(key)
	if vm, ok := l.recencyIndex[key]; ok {
		l.recency.Remove(vm)
		delete(l.recencyIndex, key)
	}
}
