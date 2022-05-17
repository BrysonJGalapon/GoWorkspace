package set

type Set[T comparable] interface {
	// Add adds the given element to this set
	Add(T)
	// Remove removes the given element from this set. Does nothing if the set does not already contain this element
	Remove(T)
	// Contains checks if the given element is in this set
	Contains(T) bool
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
	ok, value := s.data[element]
	return ok && value
}

func New[T comparable](elements ...T) Set[T] {
	s := &set[T]{data: make(map[T]bool)}

	for _, element := range elements {
		s.Add(element)
	}

	return s
}
