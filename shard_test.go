package bigcache

import (
	"fmt"
	"testing"
)

func TestShardSet(t *testing.T) {
	q := initShard(100, 2, remove)
	data := "hello"
	q.set("111", 1, []byte(data))
	q.set("111", 1, []byte(data))
	q.set("111", 1, []byte(data))
	q.set("111", 1, []byte(data))
	q.set("111", 1, []byte(data))
	data = "hello222"
	q.set("111", 1, []byte(data))

	data1, err := q.get(1, "111")
	if err == nil {
		fmt.Println(data1)
	}

}

func TestShardDel(t *testing.T) {
	q := initShard(100, 2, remove)
	data := "hello"
	q.set("111", 1, []byte(data))
	err := q.del(1, "111")

	data1, err := q.get(1, "111")
	if err == nil {
		fmt.Println(data1)
	} else {
		fmt.Println(err)
	}
}

func TestRemoveOldestEntry(t *testing.T) {
	q := initShard(100, 2 * 60, remove)
	data := "hello"
	q.set("111", 1, []byte(data))
	// q.set("111", 2, []byte(data))
	q.removeOldestEntry("test")

	data1, err := q.get(1, "111")
	if err == nil {
		fmt.Println(data1)
	} else {
		fmt.Println(err)
	}
}

func remove(wrappedEntry []byte, reason string) {
	fmt.Println(wrappedEntry, reason)
}
