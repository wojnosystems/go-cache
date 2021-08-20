# Overview

Provides a set of caches for GoLang. Cache interface is extremely simple, essentially allowing you to build your own custom cache types. These are intended to be used under the hood in your own structs to cache any item type you want.

The most basic cache is the `Unbounded` cache, which will cache everything without deleting items. The more sophisticated `LRU` caches will evict items if you store too many of them.

LRU caches come in 2 flavors:
* bounded by number of items stored
* bounded by the total size of items stored

## The ValueMapper

The ValueMapper is how you tell go-cache how to fetch cache-miss values. This function should look up values given the provided key. It is only called if there is a cache miss. So you can put any cache miss tracking into this method.

Normally, when you see a cache, you're used to seeing a "put" and a "get". By only supporting get, you remove having to handle a missing value. You will have to handle error values, but you would have done that anyway when you looked up the value before doing a traditional "put".

# Example: LRU Bounded cache

LRU caches have a maximum capacity and count each item in the cache as a single unit:

```go
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
```

This program will output:

```text
looked up: '1'
looked up: '2'
looked up: '3'
looked up: '1'
```

# Example: Byte slice LRU Bounded Cache

This is a special LRU cache that stores only byte slices. Each slice's capacity is used to limit the maximum amount of memory taken up by the values. It does not count the overhead used to store or track usage, however.

```go
package main

import (
	"context"
	"fmt"
	"github.com/wojnosystems/go-cache"
)

const (
	Byte uint = 1
	KiB       = Byte*1_024
)

func main() {
	ctx := context.Background()
	twoKibibytes := 2 * KiB
	echoLru := cache.NewLRUByte(twoKibibytes, func(ctx context.Context, key interface{}) (value []byte, err error) {
		fmt.Printf("looked up: '%s'\n", key)
		return make([]byte, 800 * Byte), nil
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
```

This program will output:

```text
looked up: '1'
looked up: '2'
looked up: '3'
looked up: '1'
```

# Example: unbounded cache (dangerous)

The foundational building block of all caches in this library is the Unbounded cache. It places no limits on the number of items it will store.

```go
package main

import (
	"context"
	"github.com/wojnosystems/go-cache"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	timesLookedUp := 0

	webCache := cache.NewUnbounded(func(ctx context.Context, key interface{}) (value interface{}, err error) {
        timesLookedUp++
        return loadWebPage(ctx, key.(string))
	})

	googlePage1, _ := webCache.Get(ctx, "https://www.google.com")
	googlePage2, _ := webCache.Get(ctx, "https://www.google.com")
	googlePage3, _ := webCache.Get(ctx, "https://www.google.com")
	wojnoPage1, _ := webCache.Get(ctx, "https://www.wojno.com")
	wojnoPage2, _ := webCache.Get(ctx, "https://www.wojno.com")

	log.Println("times looked up, should be 2: ", timesLookedUp)
	if googlePage1 == googlePage2 && googlePage2 == googlePage3 {
		log.Println("google pages are the same")
	}
	if wojnoPage1 == wojnoPage2 {
		log.Println("wojno pages are the same")
	}
}

func loadWebPage(ctx context.Context, url string) (content string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	value, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
    }
	return string(value), nil
}
```

The above program will look up google.com and wojno.com and output:

```text
2021/08/18 23:46:22 times looked up, should be 2:  2
2021/08/18 23:46:22 google pages are the same
2021/08/18 23:46:22 wojno pages are the same
```

Meaning it only looked up each page once and always returned the value in the cache when it was available.

# Interfaces and controlling usage

All caches support the "Getter" and "Invalidater" interfaces, with the LRUByte having a similar method that returns a byte array instead of an `interface{}` value.

## Getter

Allows developers to simply request items from the cache and to be auto-loaded into the cache if missing. It does not allow them to delete items from the cache.

Generally, you want to use this interface type when using caches to simplify the interface for developers using caches, but not maintaining or building new cache storage systems.

## Invalidater

Allows developers to remove items from the cache. For the unbounded cache, it will free items. For the LRU caches, it will also perform the tracking and house keeping.

Generally, you don't need to expose this to developers. This is exposed to you in case you wish to create your own sub-classes of caches and need to control this.

## Getter

# FAQ's

## How do I clear the cache?

Just make a new one :).

## Is this thread safe?

No.
