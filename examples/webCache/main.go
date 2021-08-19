// +build skip_doc_example

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
