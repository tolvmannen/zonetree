package cache

import (
	cmap "github.com/orcaman/concurrent-map/v2"
)

type concurrentMap[V any] struct {
	m cmap.ConcurrentMap[string, V]
}

// NewConcurrentMap creates a new cmap-backed Map.
func NewConcurrentMap[V any]() Map[V] {
	return &concurrentMap[V]{m: cmap.New[V]()}
}

func (c *concurrentMap[V]) Set(key string, value V) {
	c.m.Set(key, value)
}

func (c *concurrentMap[V]) Get(key string) (V, bool) {
	return c.m.Get(key)
}

func (c *concurrentMap[V]) Remove(key string) {
	c.m.Remove(key)
}

func (c *concurrentMap[V]) Has(key string) bool {
	return c.m.Has(key)
}

func (c *concurrentMap[V]) Keys() []string {
	return c.m.Keys()
}

func (c *concurrentMap[V]) IterBuffered() <-chan Tuple[V] {
	out := make(chan Tuple[V], len(c.m.Keys()))
	go func() {
		for kv := range c.m.IterBuffered() {
			out <- Tuple[V]{Key: kv.Key, Value: kv.Val}
		}
		close(out)
	}()
	return out
}

func (c *concurrentMap[V]) Count() int {
	return c.m.Count()
}

func (c *concurrentMap[V]) Clear() {
	c.m.Clear()
}

func (c *concurrentMap[V]) Pop(key string) (V, bool) {
	return c.m.Pop(key)
}

func (c *concurrentMap[V]) Upsert(key string, value V, cb UpsertFunc[V]) (V, error) {
	return c.m.Upsert(key, value, cmap.UpsertCb[V](cb)), nil
}
