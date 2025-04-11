package imissphpgo

// A pacakge for common functions that I cannot find in the standard library

import (
	"unicode"

	"golang.org/x/exp/constraints"
)

// A general Number constraint for use in generic functions
type Number interface {
	constraints.Integer | constraints.Float
}

// Capitalizes the first letter in a string
func UcFirst(s string) string {
	if len(s) <= 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

// Check if item is in array
func InArray[T comparable](val T, list []T) bool {
	for _, item := range list {
		if item == val {
			return true
		}
	}

	return false
}
