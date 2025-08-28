package cache

// Map is a full interface matching cmap.ConcurrentMap[V].
type Map[V any] interface {
	Set(key string, value V)
	Get(key string) (V, bool)
	Remove(key string)
	Has(key string) bool
	Keys() []string
	IterBuffered() <-chan Tuple[V]
	Count() int
	Clear()
	Pop(key string) (V, bool)
	Upsert(key string, value V, cb UpsertFunc[V]) (res V, err error)
}

// Tuple holds a key-value pair for iteration.
type Tuple[V any] struct {
	Key   string
	Value V
}

// UpsertFunc is the callback type for Upsert.
type UpsertFunc[V any] func(exist bool, valueInMap V, newValue V) V

func NewMapFromConfig[V any](dummy bool) Map[V] {
	if dummy == true {
		return NewDummyMap[V]()
	}
	return NewConcurrentMap[V]()
}
