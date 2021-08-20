package cache

import "context"

// ByteGetInvalidator is just like GetInvalidator, but specific for byte array values
type ByteGetInvalidator interface {
	// Get will obtain and cache the byte array for the key, err is non-nil if
	// it could not be fetched or cached
	Get(ctx context.Context, key interface{}) (value []byte, err error)

	Invalidater
}

func byteLenFromInterface(byteSlice interface{}) uint {
	return uint(cap(byteSlice.([]byte)))
}

type lruByte struct {
	*lruBase
}

type ByteMapper func(ctx context.Context, key interface{}) (value []byte, err error)

// NewLRUByte stores byte slices in a bounded LRU. Oldest items are removed to make
// space for new items. byte slice capacity is used to determine the size an entry takes up
func NewLRUByte(maxBytes uint, valueMapper ByteMapper) ByteGetInvalidator {
	l := &lruByte{
		lruBase: newLRUBase(maxBytes,
			byteLenFromInterface,
			func(ctx context.Context, key interface{}) (value interface{}, err error) {

				value, err = valueMapper(ctx, key)
				return
			}),
	}
	return l
}

// Get only exists to convert the interface to an explicit byte-array as GoLang lacks generics
// this is just a convenience function
func (b *lruByte) Get(ctx context.Context, key interface{}) (value []byte, err error) {
	var iVal interface{}
	iVal, err = b.lruBase.Get(ctx, key)
	if err != nil {
		return
	}
	return iVal.([]byte), err
}
