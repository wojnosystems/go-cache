# Overview

Provides a set of caches for GoLang. Cache interface is extremely simple, essentially allowing you to build your own custom cache types. These are intended to be used under the hood in your own structs to cache any item type you want.

The most basic cache is the `Unbounded` cache, which will cache everything without deleting items. The more sophisticated `LRU` cache will evict items if you store too many of them.

# Example

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

## The ValueMapper

The second argument to the NewUnbounded is the value mapper. This function looks up values given the provided key. It is only called if there is a cache miss. So you can put any cache miss tracking into this method.

Normally, when you see a cache, you're used to seeing a "put" and a "get". By only supporting get, you remove having to handle a missing value.

# LRU Bounded cache

You can also get a bounded cache, which has a maximum capacity:

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

# How do I clear the cache?

Just make a new one :).

# Is this thread safe?

No.
