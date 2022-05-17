package expression

import "galapb/goworkspace/pkg/set"

var ADDITION_SUBTRACTION_OPERATIONS set.Set[rune] = set.New('+', '-')
var MULTIPLICATION_DIVISION_OPERATIONS set.Set[rune] = set.New('*', '/')
