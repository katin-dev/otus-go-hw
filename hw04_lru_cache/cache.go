package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   string
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if c.capacity == 0 {
		return false
	}

	ci := cacheItem{
		key:   string(key),
		value: value,
	}

	// Может, он уже там есть?
	if current, ok := c.items[key]; ok {
		c.queue.MoveToFront(current)
		current.Value = ci
		return true
	}

	if len(c.items) == c.capacity {
		// Удаляем самый неиспользуемый
		lastItem := c.queue.Back()
		lastCacheItem := lastItem.Value.(cacheItem)
		delete(c.items, Key(lastCacheItem.key))
		c.queue.Remove(lastItem)
	}

	c.items[key] = c.queue.PushFront(ci)

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if li, ok := c.items[key]; ok {
		ci := li.Value.(cacheItem)
		c.queue.MoveToFront(li)

		return ci.value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}
