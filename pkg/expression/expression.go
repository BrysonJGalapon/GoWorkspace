package expression

// Expression is a mathematical expression only consisting of integers, operators (+, -, /, *), parenthesis, and whitespace
type Expression string

func New(expression string) (Expression, error) {
	// TODO sanity checks

	return Expression(expression), nil
}
