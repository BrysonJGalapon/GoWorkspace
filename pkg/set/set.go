package set

import "fmt"

type Set[T comparable] interface {
	// Add adds the given element to this set
	Add(T)
	// Remove removes the given element from this set. Does nothing if the set does not already contain this element
	Remove(T)
	// Contains checks if the given element is in this set
	Contains(T) bool
	// Elements returns the elements in this set
	Elements() []T
	// Union unions this set with another set, and returns the unioned set
	Union(Set[T]) Set[T]
}

type set[T comparable] struct {
	data map[T]bool
}

func (s *set[T]) Add(element T) {
	s.data[element] = true
}

func (s *set[T]) Remove(element T) {
	delete(s.data, element)
}

func (s *set[T]) Contains(element T) bool {
	_, ok := s.data[element]
	return ok
}

func (s *set[T]) Elements() []T {
	elements := make([]T, len(s.data))

	i := 0
	for element := range s.data {
		elements[i] = element
		i += 1
	}

	return elements
}

func (s *set[T]) String() string {
	ret := "{"
	i := 0
	for element := range s.data {
		ret += fmt.Sprint(element)
		if i != len(s.data)-1 {
			ret += ", "
		}
		i += 1
	}
	ret += "}"
	return ret
}

func (s *set[T]) Union(o Set[T]) Set[T] {
	ret := newSet[T]()

	for element := range s.data {
		ret.Add(element)
	}

	for _, element := range o.Elements() {
		ret.Add(element)
	}

	return ret
}

func newSet[T comparable](elements ...T) *set[T] {
	s := &set[T]{data: make(map[T]bool)}

	for _, element := range elements {
		s.Add(element)
	}

	return s
}

func New[T comparable](elements ...T) Set[T] {
	return newSet(elements...)
}
