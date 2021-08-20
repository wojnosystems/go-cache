// +build skip_doc_example

package main

import (
	"context"
	"fmt"
	"github.com/wojnosystems/go-cache"
)

func main() {
	ctx := context.Background()
	echoLru := cache.NewLRUItem(2, func(ctx context.Context, key interface{}) (value interface{}, err error) {
		fmt.Printf("looked up: '%s'\n", key)
		return key, nil
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
