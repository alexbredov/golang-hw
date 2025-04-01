package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
	Len() int
}

type lruCache struct {
	sync.RWMutex
	capacity int
	queue    List
	items    map[Key]*ListItem
	keys     map[*ListItem]Key
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	if cache == nil {
		panic("cache is nil")
	}
	cache.Lock()
	defer cache.Unlock()
	_, ok := cache.items[key]
	var litem *ListItem
	keyMatch := false
	if ok {
		litem = cache.items[key]
		litem.Value = value
		cache.queue.MoveToFront(litem)
		keyMatch = true
	} else {
		litem = cache.queue.PushFront(value)
		cache.items[key] = litem
		cache.keys[litem] = key
		if cache.queue.Len() > cache.capacity {
			litem = cache.queue.Back()
			key, ok = cache.keys[litem]
			if ok {
				delete(cache.items, key)
				delete(cache.keys, litem)
			}
			cache.queue.Remove(litem)
		}
	}
	return keyMatch
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	if cache == nil {
		panic("cache is nil")
	}
	cache.RLock()
	defer cache.RUnlock()
	_, ok := cache.items[key]
	var litem *ListItem
	var value any
	keyMatch := false
	if ok {
		litem = cache.items[key]
		value = litem.Value
		cache.queue.MoveToFront(litem)
		keyMatch = true
	}
	return value, keyMatch
}

func (cache *lruCache) Clear() {
	if cache == nil {
		panic("cache is nil")
	}
	cache.Lock()
	defer cache.Unlock()
	for key, value := range cache.items {
		cache.queue.Remove(value)
		delete(cache.items, key)
	}
	for key := range cache.keys {
		delete(cache.keys, key)
	}
	clear(cache.items)
}

func (cache *lruCache) Len() int {
	return cache.queue.Len()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		keys:     make(map[*ListItem]Key, capacity),
	}
}
