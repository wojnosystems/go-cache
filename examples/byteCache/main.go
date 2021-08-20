package main

import (
	"context"
	"fmt"
	"github.com/wojnosystems/go-cache"
)

const (
	Byte uint = 1
	KiB       = Byte * 1_024
)

func main() {
	ctx := context.Background()
	twoKibibytes := 2 * KiB
	echoLru := cache.NewLRUByte(twoKibibytes, func(ctx context.Context, key interface{}) (value []byte, err error) {
		fmt.Printf("looked up: '%s'\n", key)
		return make([]byte, 800*Byte), nil
	})

	// 1 will be looked up
	_, _ = echoLru.Get(ctx, "1")

	// 2 will be looked up
	_, _ = echoLru.Get(ctx, "2")

	// this will evict 1 and look up 3
	_, _ = echoLru.Get(ctx, "3")

	// 2 is still cached, so no look up
	_, _ = echoLru.Get(ctx, "2")

	// 1 will need to be looked up again
	_, _ = echoLru.Get(ctx, "1")
}
