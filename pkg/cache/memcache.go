package cache

import (
	"container/list"
	"sync"
	"sync/atomic"
	//"fmt"
)

// An AtomicInt is an int64 to be accessed atomically.
type AtomicInt int64

// MemCache is an LRU cache. It is safe for concurrent access.
type MemCache struct {
	mutex       sync.RWMutex
	maxItemSize int
	cacheList   *list.List
	cache       map[interface{}]*list.Element
	hits, gets  AtomicInt
}

// MemCacheStatus return status of chache
type MemCacheStatus struct {
	Gets        int64
	Hits        int64
	MaxItemSize int
	CurrentSize int
}

type entry struct {
	key   interface{}
	value interface{}
}

// If maxItemSize is zero, the cache has no limit.
//if maxItemSize is not zero, when cache's size beyond maxItemSize,start to swap
func NewMemCache(maxItemSize int) *MemCache {
	return &MemCache{
		maxItemSize: maxItemSize,
		cacheList:   list.New(),
		cache:       make(map[interface{}]*list.Element),
	}
}

//return the status of cache
func (c *MemCache) Status() interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return &MemCacheStatus{
		MaxItemSize: c.maxItemSize,
		CurrentSize: c.cacheList.Len(),
		Gets:        c.gets.Get(),
		Hits:        c.hits.Get(),
	}
}

//get value with key
func (c *MemCache) Get(key interface{}) (value interface{}, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.gets.Add(1)
	if ele, hit := c.cache[key]; hit {
		c.hits.Add(1)
		return ele.Value.(*entry).value, true
	}
	return
}

//set a value with key
func (c *MemCache) Set(key, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.cacheList = list.New()
	}

	if ele, ok := c.cache[key]; ok {
		c.cacheList.MoveToFront(ele)
		ele.Value.(*entry).value = value
		return
	}

	ele := c.cacheList.PushFront(&entry{key: key, value: value})
	c.cache[key] = ele
	if c.maxItemSize != 0 && c.cacheList.Len() > c.maxItemSize {
		c.RemoveOldest()
	}
}

func (c *MemCache) Delete(key interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		c.cacheList.Remove(ele)
		key := ele.Value.(*entry).key
		delete(c.cache, key)
		return
	}
}

func (c *MemCache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.cacheList.Back()
	if ele != nil {
		c.cacheList.Remove(ele)
		key := ele.Value.(*entry).key
		delete(c.cache, key)
	}
}

// Add atomically adds n to i.
func (i *AtomicInt) Add(n int64) {
	atomic.AddInt64((*int64)(i), n)
}

// Get atomically gets the value of i.
func (i *AtomicInt) Get() int64 {
	return atomic.LoadInt64((*int64)(i))
}