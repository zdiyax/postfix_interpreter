package cmd

import (
	"slices"
)

var validOperators = []string{"+", "-", "*", "/", "="}

// IsValidOperator checks if the operator is in the list of supported operators
func IsValidOperator(operator string) bool {
	return slices.Contains(validOperators, operator)
}
