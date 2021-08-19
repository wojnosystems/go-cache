package cache

import "context"

type stringKeyCache map[interface{}]interface{}

type unbounded struct {
	cache        stringKeyCache
	valueFactory ValueMapper
}

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

func (u *unbounded) Invalidate(key interface{}) {
	delete(u.cache, key)
}
