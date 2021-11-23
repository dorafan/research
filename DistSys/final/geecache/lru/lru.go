package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	//Maximum memory allowed
	maxBytes int64
	//Currently used memory
	nbytes   int64
	//Doubly linked list
	ll       *list.List
	//The key is a string, and the value is a pointer to the corresponding node in the doubly linked list
	cache    map[string]*list.Element
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

//Data type of doubly linked list node
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	//If the key exists
	if ele, ok := c.cache[key]; ok {
		//move the node to the end of the queue
		c.ll.MoveToFront(ele)
		//update the value of the corresponding node
		kv := ele.Value.(*entry)
		//update current used memory
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {//If it does not exist, add a new node
		//Add a new node &entry{key, value} at the end of the team,
		ele := c.ll.PushFront(&entry{key, value})
		//add the mapping relationship between key and node in the dictionary.
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	//Update c.nbytes, if it exceeds the set maximum value c.maxBytes, remove the least visited node.
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Get look ups a key's value
//find the node of the corresponding doubly linked list from the dictionary.
//then move the node to the end of the queue.
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item
//i.e remove the least recently visited node
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		//Get the head node of the team and delete it from the linked list
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		//Delete the mapping relationship of the node from the dictionary c.cache
		delete(c.cache, kv.key)
		//Update the currently used memory
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
