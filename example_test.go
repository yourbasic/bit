package bit_test

import (
	"fmt"
	"github.com/yourbasic/bit"
)

// Compute the sum of all elements in a set.
func ExampleSet_Visit() {
	s := bit.New(1, 2, 3, 4)
	sum := 0
	s.Visit(func(n int) (skip bool) {
		sum += n
		return
	})
	fmt.Println("sum", s, "=", sum)
	// Output: sum {1..4} = 10
}

// Abort an iteration in mid-flight.
func ExampleSet_Visit_abort() {
	s := bit.New(2, 3, 5, 7, 11, 13)

	// Print all single digit numbers in s.
	s.Visit(func(n int) (skip bool) {
		if n >= 10 {
			skip = true
			return
		}
		fmt.Print(n, " ")
		return
	})
	// Output: 2 3 5 7
}

// How to create, combine, compare and print bitsets.
func Example_basics() {
	// Add all elements in the range [0, 100) to the empty set.
	A := new(bit.Set).AddRange(0, 100) // {0..99}

	// Create a new set with the elements 0 and 200, and then add [50, 150).
	B := bit.New(0, 200).AddRange(50, 150) // {0 50..149 200}

	// Compute the symmetric difference X = A △ B, also known as XOR.
	X := A.AndNot(B).Or(B.AndNot(A)) // (A ∖ B) ∪ (B ∖ A)

	// Compute A △ B in a different way.
	Y := A.Or(B).AndNot(A.And(B)) // (A ∪ B) ∖ (A ∩ B)

	// Compare the results.
	if X.Equal(Y) {
		fmt.Println(X)
	}
	// Output: {1..49 100..149 200}
	//
}

func ExampleNew() {
	fmt.Println(bit.New(0, 1, 10, 10, -1))
	// Output: {0 1 10}
}

func ExampleSet_String() {
	fmt.Println(bit.New(1, 2, 6, 5, 3))
	// Output: {1..3 5 6}
}
