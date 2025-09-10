package cache

// A simple cache
type Cache struct {
    data  map[string][]byte
    order []string
    cap   int
}

func NewCache(capacity int) *Cache {
    return &Cache{
        data:  make(map[string][]byte),
        order: make([]string, 0, capacity),
        cap:   capacity,
    }
}

func (c *Cache) Get(key string) ([]byte, bool) {
    val, exists := c.data[key]
    return val, exists
}

func (c *Cache) Put(key string, value []byte) {
    if _, exists := c.data[key]; exists {
        c.data[key] = value
        return
    }
    if len(c.order) >= c.cap {
        oldest := c.order[0]
        c.order = c.order[1:]
        delete(c.data, oldest)
    }
    c.data[key] = value
    c.order = append(c.order, key)
}
