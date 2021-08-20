package cache

import "context"

type Getter interface {
	/*
		Get a value from the cache. If missing, it will be loaded, then cached, then returned.
		if it exists it will be returned without attempting to load it again.

		ctx is passed because the underlying provider may need to make external calls
		key of the value to look up in the cache
		value is the mapped value that is cached
		err is non-null if there was a failure to cache or look up the value from the key
	*/
	Get(ctx context.Context, key interface{}) (value interface{}, err error)
}

type Invalidater interface {
	/*
			Invalidate marks a cached key as invalid. The next request for this key is guaranteed to be a fresh load
		however, implementers are under no obligation to clear the cached item immediately, it can be deferred
	*/
	Invalidate(key interface{})
}

type GetInvalidater interface {
	Getter
	Invalidater
}

/*
ValueMapper is the method that allows the cache to obtain uncached values

ctx: the context passed to Get calls by the caller
key: passed to Get calls by the caller
*/
type ValueMapper func(ctx context.Context, key interface{}) (value interface{}, err error)
