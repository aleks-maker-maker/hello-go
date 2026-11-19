package greeting

import "fmt"

// Greet returns a greeting for name.
func Greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

// SumRange returns the sum of an inclusive range of integers.
// An empty range has sum zero.
func SumRange(from, to int) int {
	if from > to {
		return 0
	}

	sum := from
	for current := from; current < to; {
		current++
		sum += current
	}
	return sum
}
