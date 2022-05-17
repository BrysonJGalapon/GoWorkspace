package stack

import "fmt"

type Stack[T any] interface {
	// Push pushes the given element to the top of the stack
	Push(element T)
	// Pop removes and returns the top-most element of the stack. Returns an error if the stack is empty
	Pop() (T, error)
	// Peek returns the top-most element of the stack. Returns an error if the stack is empty
	Peek() (T, error)
	// Size returns the number of elements currently in the stack
	Size() int
	// IsEmpty checks if the size of the stack is 0
	IsEmpty() bool
}

type stack[T any] struct {
	data []T
}

func (s *stack[T]) Push(element T) {
	s.data = append(s.data, element)
}

func (s *stack[T]) Pop() (T, error) {
	element, err := s.Peek()
	if err != nil {
		return element, err
	}

	s.data = s.data[:s.Size()-1]

	return element, nil
}

func (s *stack[T]) Peek() (T, error) {
	var element T

	if s.IsEmpty() {
		return element, fmt.Errorf("can not pop from an empty stack")
	}

	return s.data[s.Size()-1], nil
}

func (s *stack[T]) Size() int {
	return len(s.data)
}

func (s *stack[T]) IsEmpty() bool {
	return s.Size() == 0
}

// New returns a new stack
func New[T any](data ...T) Stack[T] {
	return &stack[T]{data: data}
}
