package bigcache

import (
	"github.com/ghorges/mybigcache/queue"
	"sync"
	"time"
)

type cacheShard struct {
	hashMap     map[uint64]uint32
	cache       queue.Queue
	lock        sync.Mutex
	entryBuffer []byte
	lifeWindow  int64
}

func (s *cacheShard) set(key string, hashedKey uint64, entry []byte) error {
	currentTimestamp := time.Now().Unix()

	s.lock.Lock()
	defer func() {
		s.lock.Unlock()
	}()

	// if have same key,discard old key
	if index := s.hashMap[hashedKey]; index != 0 {
		if data, err := s.cache.Get(int(index)); err != nil {
			hashKeyToZero(data)
		}
	}


	wrapEntry(currentTimestamp, hashedKey, key, entry, &s.entryBuffer)
}
