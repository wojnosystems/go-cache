package lru

import (
	"container/list"
)

// recencyIndexType maps keys to recency elements
type recencyIndexType map[interface{}]*list.Element

type tracker struct {
	// recency of items (interface{}),
	// front = oldest,
	// back = most recently used
	recency *list.List

	// recencyIndex O(1) lookup for items in the recency list
	recencyIndex recencyIndexType
}

func NewTracker() Tracker {
	return &tracker{
		recency:      list.New(),
		recencyIndex: make(recencyIndexType),
	}
}

func (l *tracker) Touch(key interface{}) {
	vm, ok := l.recencyIndex[key]
	if ok {
		l.recency.Remove(vm)
	}
	l.recencyIndex[key] = l.recency.PushBack(key)
}

func (l *tracker) Remove(key interface{}) {
	if vm, ok := l.recencyIndex[key]; ok {
		l.recency.Remove(vm)
		delete(l.recencyIndex, key)
	}
}

func (l *tracker) LRU() (key interface{}, ok bool) {
	if front := l.recency.Front(); front != nil {
		return front.Value, true
	}
	return nil, false
}

func (l *tracker) Len() int {
	return l.recency.Len()
}
