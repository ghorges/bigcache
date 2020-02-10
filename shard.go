package bigcache

import (
	"errors"
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
	onRemove    onRemoveCallback
}

func (s *cacheShard) set(key string, hashedKey uint64, entry []byte) error {
	currentTimestamp := time.Now().Unix()

	s.lock.Lock()
	defer func() {
		s.lock.Unlock()
	}()

	// if have same key,discard old key
	if index := s.hashMap[hashedKey]; index != 0 {
		if data, err := s.cache.Get(int(index)); err == nil {
			hashKeyToZero(data)
		}
	}

	if data, err := s.cache.Peek(); err == nil {
		s.onEvict(data, uint64(currentTimestamp), s.removeOldestEntry)
	}

	wrapData := wrapEntry(currentTimestamp, hashedKey, key, entry, &s.entryBuffer)

	index := s.cache.Push(wrapData)
	s.hashMap[hashedKey] = uint32(index)

	return nil
}

func (s *cacheShard) get(hashedKey uint64, key string) ([]byte, error) {
	s.lock.Lock()
	defer func() {
		s.lock.Unlock()
	}()

	index := s.hashMap[hashedKey]

	if index == 0 {
		return nil, errors.New("Entry not found")
	}

	data, err := s.cache.Get(int(index))

	if err != nil {
		return nil, err
	}

	if hashedKey != getHashKey(data) {
		return nil, errors.New("Entry not found")
	}

	if key != string(getKey(data)) {
		return nil, err
	}

	return readEntry(data), nil
}

func (s *cacheShard) onEvict(oldestEntry []byte, currentTimestamp uint64, evict func(reason string) error) bool {
	olderTime := readTimestamp(oldestEntry)

	if currentTimestamp-olderTime < uint64(s.lifeWindow) {
		return false
	}

	evict("evict,time out")
	return true

}

func (s *cacheShard) del(hashedKey uint64, key string) error {
	s.lock.Lock()
	defer func() {
		s.lock.Unlock()
	}()

	index := s.hashMap[hashedKey]

	if index == 0 {
		return errors.New("Entry not found")
	}
	data, err := s.cache.Get(int(index))
	if err != nil {
		return err
	}
	if getHashKey(data) != hashedKey {
		return errors.New("Entry not found")
	}
	if string(getKey(data)) != key {
		return errors.New("Entry not found")
	}

	delete(s.hashMap, hashedKey)
	s.onRemove(data, "del")

	hashKeyToZero(data)
	return nil
}

func (s *cacheShard) removeOldestEntry(reason string) error {
	data, err := s.cache.Pop()

	if err != nil {
		return err
	}

	hashKey := getHashKey(data)
	delete(s.hashMap, hashKey)
	s.onRemove(data, reason)

	return nil
}

func (s *cacheShard) cleanUp(currentTimestamp int64) {
	s.lock.Lock()
	defer func() {
		s.lock.Unlock()
	}()

	for {
		data, err := s.cache.Peek()
		if err != nil {
			return
		}

		if s.onEvict(data, uint64(currentTimestamp), s.removeOldestEntry) == false {
			return
		}
	}
}

func initShard(initSize int, lifeWindow int64, callback onRemoveCallback) *cacheShard {
	return &cacheShard{
		hashMap:     make(map[uint64]uint32, initSize),
		cache:       *queue.NewBytesQueue(initSize),
		lock:        sync.Mutex{},
		entryBuffer: make([]byte, 0),
		lifeWindow:  lifeWindow,
		onRemove:    callback,
	}
}
