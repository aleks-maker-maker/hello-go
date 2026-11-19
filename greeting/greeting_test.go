package greeting

import (
	"math"
	"testing"
)

func TestGreet(t *testing.T) {
	got := Greet("Docker")
	want := "Hello, Docker!"
	if got != want {
		t.Errorf("Greet() = %q, want %q", got, want)
	}
}

func TestSumRange(t *testing.T) {
	tests := []struct {
		name string
		from int
		to   int
		want int
	}{
		{name: "one_to_ten", from: 1, to: 10, want: 55},
		{name: "crosses_zero", from: -3, to: 1, want: -5},
		{name: "single_number", from: 4, to: 4, want: 4},
		{name: "empty", from: 5, to: 2, want: 0},
		{name: "max_int", from: math.MaxInt, to: math.MaxInt, want: math.MaxInt},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := SumRange(test.from, test.to)
			if got != test.want {
				t.Errorf("SumRange(%d, %d) = %d, want %d", test.from, test.to, got, test.want)
			}
		})
	}
}
