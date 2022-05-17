package expression

import (
	"bytes"
	"fmt"
	"galapb/goworkspace/pkg/set"
	"galapb/goworkspace/pkg/stack"
	"galapb/goworkspace/pkg/tree"
	"strconv"
	"unicode"
)

// Evaluator evaluates an Expression and returns its equivalent value
type Evaluator interface {
	Evaluate(Expression) (value float64)
}

type evaluator struct{}

func (e *evaluator) Evaluate(expression Expression) (value float64) {
	// handle nested expressions with a stack
	stack := stack.New[bytes.Buffer]()

	// maintain the current expression as a byte buffer that can be appended to
	var currExpression bytes.Buffer = bytes.Buffer{}

	for _, char := range expression {
		// ignore whitespace characters
		if unicode.IsSpace(char) {
			continue
		}

		// an open parenthesis signals a nested expression, so:
		//	1. push current buffer onto the stack
		//	2. reset the current buffer
		if char == '(' {
			stack.Push(currExpression)
			currExpression = bytes.Buffer{}
			continue
		}

		// a closed parenthesis signals the end of the nested expression, so:
		//	1. evaluate the nested expression
		//	2. append the value to the most recently saved buffer
		//	3. set the current buffer as that buffer
		if char == ')' {
			expr := Expression(currExpression.String())
			v := evaluateNonNestedExpression(expr)
			currExpression, _ = stack.Pop()
			currExpression.WriteString(fmt.Sprintf("%f", v))
			continue
		}

		// append any non-parenthesis characters to the current buffer
		currExpression.WriteRune(char)
	}

	// all nested expressions have been evaluated, so evaluate the root-level expression
	expr := Expression(currExpression.String())
	return evaluateNonNestedExpression(expr)
}

// evaluateNonNestedExpression evaluates the given expression, assuming that it does not contain any open or closed parenthesis
func evaluateNonNestedExpression(expression Expression) (value float64) {
	if expression == "" {
		panic("can't evaluate an empty expression")
	}

	// initialize AST as a single node containing the full expression
	tokens := []string{string(expression)}
	ast := buildAst(tokens)

	// expand AST in reverse order of operation priority, so that least-priority operations get evaluated last
	ast = expandLeaves(ast, ADDITION_SUBTRACTION_OPERATIONS)
	ast = expandLeaves(ast, MULTIPLICATION_DIVISION_OPERATIONS)

	// evaluate the fully-expanded AST
	return evaluateAst(ast)
}

// evaluateAst evaluates the AST in post-order fashion
func evaluateAst(ast *tree.Node[string]) (value float64) {
	// ast is a leaf node, so evaluate it as a constant
	if len(ast.Children) == 0 {
		return evaluateConstant(ast.Data)
	}

	// ast is an inner node, so evaluate all of its children ...
	values := make([]float64, len(ast.Children))
	for i, child := range ast.Children {
		values[i] = evaluateAst(child)
	}

	// ... and then perform the operation on all the children
	operation := rune(ast.Data[0])
	return reduceValuesByOperation(values, operation)
}

func reduceValuesByOperation(values []float64, operation rune) float64 {
	switch operation {
	case '+':
		return reduce(add, values)
	case '-':
		return reduce(subtract, values)
	case '*':
		return reduce(multiply, values)
	case '/':
		return reduce(divide, values)
	default:
		panic("unhandled operation: " + string(operation))
	}
}

func evaluateConstant(str string) float64 {
	value, _ := strconv.ParseFloat(str, 64)
	return value
}

// expandLeaves replaces each leaf node in the given AST with an equivalent AST that was generated from the tokens
// created by splitting the data at the leaf node by the given operations, and returns a pointer to the root node of the new tree
func expandLeaves(ast *tree.Node[string], operations set.Set[rune]) (root *tree.Node[string]) {
	// initially assume that the root node of the returned AST will not be different than the original AST
	root = ast

	// stackItem maintains parent-child relationships to enable replacement of nodes within the AST
	type stackItem struct {
		node       *tree.Node[string]
		parent     *tree.Node[string]
		childIndex int
	}

	// initialize the stack with the original AST
	stack := stack.New(stackItem{node: ast, parent: nil, childIndex: 0})

	for !stack.IsEmpty() {
		item, _ := stack.Pop()

		node := item.node

		// node is not a leaf, so drill down further into the tree
		if len(node.Children) != 0 {
			for childIndex, child := range node.Children {
				stack.Push(stackItem{node: child, parent: node, childIndex: childIndex})
			}
			continue
		}

		// node is a leaf, so build an equivalent AST using the data in the node
		tokens := splitStringByOperations(node.Data, operations)
		leafAst := buildAst(tokens)

		// replace the node in the tree
		if node == root {
			root = leafAst // edge case in which the original AST is itself a leaf node, so replace the entire tree
		} else {
			parent := item.parent
			childIndex := item.childIndex
			parent.Children[childIndex] = leafAst
		}
	}

	return root
}

// splitStringByOperations splits the given string into a slice of tokens, delimited by the given operations, while preserving the
// the operations in the returned token slice. In other words, concatenating all the tokens in the returned token slice will
// be equivalent to the given string
func splitStringByOperations(str string, operations set.Set[rune]) (tokens []string) {
	tokens = make([]string, 0)

	token := bytes.Buffer{}
	for _, char := range str {
		if operations.Contains(char) {
			tokens = append(tokens, token.String())
			token = bytes.Buffer{}

			tokens = append(tokens, string(char))
			continue
		}

		token.WriteRune(char)
	}

	tokens = append(tokens, token.String())

	return tokens
}

// buildAst builds a abstract syntax tree (AST) from the given tokens. Assumes the length of tokens is odd
func buildAst(tokens []string) *tree.Node[string] {
	var root *tree.Node[string] = tree.NewNode(tokens[0])

	var (
		leftChild  *tree.Node[string]
		rightChild *tree.Node[string]
	)

	// evaluation executes from left to right, so build AST with a "heavy" left side
	for i := 1; i < len(tokens); i += 2 {
		leftChild = root
		rightChild = tree.NewNode(tokens[i+1])

		root = tree.NewNode(tokens[i], leftChild, rightChild)
	}

	return root
}

func NewEvaluator() Evaluator {
	return &evaluator{}
}
