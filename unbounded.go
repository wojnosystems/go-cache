package cache

import "context"

type stringKeyCache map[interface{}]interface{}

type unbounded struct {
	cache        stringKeyCache
	valueFactory ValueMapper
}

// NewUnbounded creates a cache without any internal limits on how many items
// can be cached. It will grow, unbounded, until you stop using it.
// While this is probably fine for testing or building up other caches,
// you probably should not use this in production
func NewUnbounded(valueFactory ValueMapper) GetInvalidater {
	return newUnbounded(valueFactory)
}

func newUnbounded(valueFactory ValueMapper) *unbounded {
	return &unbounded{
		cache:        make(stringKeyCache),
		valueFactory: valueFactory,
	}
}

func (u *unbounded) Get(ctx context.Context, key interface{}) (value interface{}, err error) {
	var ok bool
	if value, ok = u.cache[key]; ok {
		return
	}
	value, err = u.valueFactory(ctx, key)
	if err != nil {
		return
	}
	u.cache[key] = value
	return
}

func (u *unbounded) Invalidate(key interface{}) (ok bool) {
	if _, ok = u.cache[key]; ok {
		delete(u.cache, key)
	}
	return
}
