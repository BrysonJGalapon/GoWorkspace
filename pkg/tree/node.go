package tree

type Node[T any] struct {
	Data     T
	Children []*Node[T]
}

func NewNode[T any](data T, children ...*Node[T]) *Node[T] {
	return &Node[T]{Data: data, Children: children}
}
