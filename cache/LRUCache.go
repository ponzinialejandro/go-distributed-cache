package cache

import (
	contianer "container/list"
)

type LRUCache struct {
	size int
	hash map[string]string
	dll  *contianer.List
}

func NewLRUCache(size int) LRUCache {
	return LRUCache{size: size, hash: make(map[string]string), dll: contianer.New()}
}

func (c *LRUCache) Put(key string, value string) {
	if _, found := c.hash[key]; !found {
		if len(c.hash) == c.size {
			delete(c.hash, c.dll.Back().Value.(string))
			c.dll.Remove(c.dll.Back())
		}
		c.dll.PushFront(key)
	}

	c.hash[key] = value
}

func (c *LRUCache) Get(key string) (string, bool) {
	if value, ok := c.hash[key]; ok {
		c.deleteKeyFromList(key)
		c.dll.PushFront(key)
		return value, ok
	}
	return "", false
}

func (c *LRUCache) deleteKeyFromList(key string) {
	ele := c.dll.Front()
	for ; ele != nil && ele.Value.(string) != key; ele = ele.Next() {
	}

	if ele != nil {
		c.dll.Remove(ele)
	}
}
