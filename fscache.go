package fscache

import (
	"encoding/gob"
	"io"
	"os"
	"sort"
	"sync"
)

type Cache struct {
	mu       sync.RWMutex
	filepath string
	Items    map[string]any
}

type Item struct {
	Key   string
	Value any
}

func init() {
	gob.Register(&Item{})
}

// Register registers types to be used with the cache (gob).
func Register(values ...any) {
	for _, a := range values {
		gob.Register(a)
	}
}

// NewCache returns a new cache instances.
func NewCache(path string) *Cache {
	c := &Cache{
		mu:       sync.RWMutex{},
		filepath: path,
		Items:    make(map[string]any),
	}

	err := c.Save()
	if err != nil {
		panic(err)
	}

	return c
}

// NewItem returns a new cache item.
func NewItem(key string, value any) *Item {
	return &Item{
		Key:   key,
		Value: value,
	}
}

// Get returns a value from the cache.
func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.Items[key]
	if !ok {
		return nil, false
	}
	return item.(*Item).Value, true
}

// Set sets a value in the cache.
func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Items[key] = NewItem(key, value)
}

// Delete deletes a value from the cache.
func (c *Cache) Delete(key string) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	delete(c.Items, key)
}

// Clear clears the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Items = make(map[string]any)
}

// Size returns the size of the cache.
func (c *Cache) Size() int {
	return len(c.Items)
}

// Keys returns the keys of the cache.
func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0)
	for k := range c.Items {
		keys = append(keys, k)
	}
	return keys
}

// Sort returns the sorted keys of the cache.
func (c *Cache) Sort() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := c.Keys()
	sort.Strings(keys)
	return keys
}

// Save saves the cache to a file.
func (c *Cache) Save() error {
	f, err := os.OpenFile(c.filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	return encode(f, c.Items)
}

// Load loads the cache from a file.
func Load(path string) (*Cache, error) {
	c := &Cache{filepath: path}

	f, err := os.Open(path)
	if err != nil {
		return c, err
	}
	defer f.Close()

	err = decode(f, &c.Items)
	return c, err
}

func encode(w io.Writer, items map[string]any) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(items)
}

func decode(r io.Reader, items *map[string]any) error {
	dec := gob.NewDecoder(r)
	return dec.Decode(items)
}
