package cache

type utestMap[V any] struct {
	m map[string]V
}

func NewutestMap[V any]() Map[V] {
	return utestMap[V]{}
}

func (u utestMap[V]) Set(key string, value V) {
	u.m[key] = value
}
func (u utestMap[V]) Get(key string) (V, bool) {
	var v V
	if v, ok := u.m[key]; ok {
		return v, true
	}
	return v, false
}

// Functions below are not used

func (u utestMap[V]) Remove(key string)   {}
func (u utestMap[V]) Has(key string) bool { return false }
func (u utestMap[V]) Keys() []string      { return nil }
func (u utestMap[V]) IterBuffered() <-chan Tuple[V] {
	ch := make(chan Tuple[V])
	close(ch)
	return ch
}

func (u utestMap[V]) Count() int {
	return len(u.m)
}

func (u utestMap[V]) Clear()                                                  {}
func (u utestMap[V]) Pop(key string) (V, bool)                                { var zero V; return zero, false }
func (u utestMap[V]) Upsert(key string, value V, cb UpsertFunc[V]) (V, error) { return value, nil }
