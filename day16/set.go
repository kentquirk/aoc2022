package main

type Set[T comparable] interface {
	Add(val T)
	Clone() Set[T]
	Contains(val T) bool
	Remove(val T)
	Each(func(T) bool)
}

type set[T comparable] map[T]struct{}

func NewSet[T comparable](vals ...T) Set[T] {
	s := make(set[T])
	for _, item := range vals {
		s.Add(item)
	}
	return &s
}

func (s *set[T]) Add(v T) {
	(*s)[v] = struct{}{}
}

func (s *set[T]) Clone() Set[T] {
	clone := make(set[T])
	for v := range *s {
		clone.Add(v)
	}
	return &clone
}

func (s *set[T]) Contains(v T) bool {
	_, ok := (*s)[v]
	return ok
}

func (s *set[T]) Remove(v T) {
	delete((*s), v)
}

func (s *set[T]) Each(f func(T) bool) {
	for v := range *s {
		if f(v) {
			break
		}
	}
}
