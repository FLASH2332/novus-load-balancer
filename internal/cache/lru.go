package cache

import "github.com/FLASH2332/novus-load-balancer/pkg"

// LRUCache struct
type LRUCache struct {
	capacity int
	data     map[string]*pkg.Node
	order    *pkg.DoublyLinkedList
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		data:     make(map[string]*pkg.Node),
		order:    pkg.NewDoublyLinkedList(),
	}
}

func (c *LRUCache) Get(key string) ([]byte, bool) {
	if node, exists := c.data[key]; exists {
		c.order.MoveToFront(node) // mark as recently used
		return node.Value, true
	}
	return nil, false
}

func (c *LRUCache) Put(key string, value []byte) {
	// If already exists, update & move to front
	if node, exists := c.data[key]; exists {
		node.Value = value
		c.order.MoveToFront(node)
		return
	}

	// If full, evict least recently used
	if len(c.data) >= c.capacity {
		evicted := c.order.PopBack()
		if evicted != nil {
			delete(c.data, evicted.Key)
		}
	}

	// Insert new node
	newNode := &pkg.Node{Key: key, Value: value}
	c.order.PushFront(newNode)
	c.data[key] = newNode
}
