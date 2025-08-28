package cache

type dummyMap[V any] struct{}

func NewDummyMap[V any]() Map[V] {
	return dummyMap[V]{}
}

func (d dummyMap[V]) Set(key string, value V)                                 {}
func (d dummyMap[V]) Get(key string) (V, bool)                                { var zero V; return zero, false }
func (d dummyMap[V]) Remove(key string)                                       {}
func (d dummyMap[V]) Has(key string) bool                                     { return false }
func (d dummyMap[V]) Keys() []string                                          { return nil }
func (d dummyMap[V]) IterBuffered() <-chan Tuple[V]                           { ch := make(chan Tuple[V]); close(ch); return ch }
func (d dummyMap[V]) Count() int                                              { return 0 }
func (d dummyMap[V]) Clear()                                                  {}
func (d dummyMap[V]) Pop(key string) (V, bool)                                { var zero V; return zero, false }
func (d dummyMap[V]) Upsert(key string, value V, cb UpsertFunc[V]) (V, error) { return value, nil }
