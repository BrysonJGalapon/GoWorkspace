package expression

import (
	"math"
	"testing"
)

type test struct {
	expression      string
	expectedValue   float64
	isErrorExpected bool
}

func runTest(tst test, t *testing.T) {
	expression, err := New(tst.expression)
	expectedValue := tst.expectedValue
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	evaluator := NewEvaluator()

	actualValue, err := evaluator.Evaluate(expression)

	if tst.isErrorExpected && (err == nil) {
		t.Fatalf("Expression: %s, error expected, but got nil", tst.expression)
	}

	if !tst.isErrorExpected && (err != nil) {
		t.Fatalf("Expression: %s, error not expected, but got %s", tst.expression, err)
	}

	if !tst.isErrorExpected && !isFloat64Equal(actualValue, expectedValue) {
		t.Fatalf("Expression: %s, Expected: %f, but got: %f", tst.expression, expectedValue, actualValue)
	}
}

func isFloat64Equal(x, y float64) bool {
	tolerance := 0.001
	return math.Abs(x-y) <= tolerance
}

func TestEvaluateTestCase1(t *testing.T) {
	runTest(test{expression: "1+1", expectedValue: 2.0}, t)
}

func TestEvaluateTestCase2(t *testing.T) {
	runTest(test{expression: "(3+3)*7", expectedValue: 42.0}, t)
}

func TestEvaluateTestCase3(t *testing.T) {
	runTest(test{expression: "(3- 2)*2* 1+3.5", expectedValue: 5.5}, t)
}

func TestEvaluateTestCase4(t *testing.T) {
	runTest(test{expression: "(3/2)", expectedValue: 1.5}, t)
}

func TestEvaluateTestCase5(t *testing.T) {
	runTest(test{expression: "(-3-2)*2", expectedValue: -10.0}, t)
}

func TestEvaluateTestCase6(t *testing.T) {
	runTest(test{expression: "(-3-2)*-2", expectedValue: 10.0}, t)
}
