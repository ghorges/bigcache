package bigcache_test

import (
	"fmt"
	"testing"
	"time"

	bigcache "github.com/ghorges/mybigcache"
)

func TestBigCacheSet(t *testing.T) {
	cache := bigcache.NewBigCache(1024, 10*time.Minute, 2*time.Minute, removeCall, 100)

	data := "hello"
	cache.Set("123", []byte(data))

	// test Del
	// cache.Del("123")

	data1, err := cache.Get("123")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("data is ", data1)
}

func TestBigCacheTimeout(t *testing.T) {
	cache := bigcache.NewBigCache(1024, 1*time.Minute, 1*time.Minute, removeCall, 100)

	data := "hello"
	cache.Set("123", []byte(data))
	cache.Set("231", []byte(data))
	cache.Set("312", []byte(data))

	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			cache.Set("keys", []byte("bigCache"))
		}
	}

}

func removeCall(wrappedEntry []byte, reason string) {
	fmt.Println("oncall:", wrappedEntry, "reason is:", reason)
}
