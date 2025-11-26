package cache

import (
	"maps"
	"slices"
)

type simpleMap[V any] struct {
	m map[string]V
}

func NewsimpleMap[V any]() Map[V] {
	return simpleMap[V]{}
}

func (s simpleMap[V]) Set(key string, value V) {
	s.m[key] = value
}
func (s simpleMap[V]) Get(key string) (V, bool) {
	var v V
	if v, ok := s.m[key]; ok {
		return v, true
	}
	return v, false
}

func (s simpleMap[V]) Remove(key string) {
	delete(s.m, key)
}

func (s simpleMap[V]) Has(key string) bool {
	if _, ok := s.m[key]; ok {
		return true
	}
	return false
}

func (s simpleMap[V]) Keys() []string {
	k := slices.Collect(maps.Keys(s.m))
	return k
}

func (s simpleMap[V]) IterBuffered() <-chan Tuple[V] {
	ch := make(chan Tuple[V])
	close(ch)
	return ch
}

func (s simpleMap[V]) Count() int {
	return len(s.m)
}

func (s simpleMap[V]) Clear() {
	clear(s.m)
}

func (s simpleMap[V]) Pop(key string) (V, bool) {
	var v V
	if v, ok := s.m[key]; ok {
		delete(s.m, key)
		return v, true
	}
	return v, false
}

func (s simpleMap[V]) Upsert(key string, value V, cb UpsertFunc[V]) (V, error) {
	s.m[key] = value
	return value, nil
}

// NoOp
//
// Dummy function to use for No Operatioond
func NoOp() {
}
