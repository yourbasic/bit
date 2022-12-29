package bit_test

import (
	"fmt"
	"math"

	"github.com/yourbasic/bit"
)

// Create, combine, compare and print bit sets.
func Example_basics() {
	// Add all elements in the range [0, 100) to the empty set.
	A := new(bit.Set).AddRange(0, 100) // {0..99}

	// Create a new set containing the two elements 0 and 200,
	// and then add all elements in the range [50, 150) to the set.
	B := bit.New(0, 200).AddRange(50, 150) // {0 50..149 200}

	// Compute the symmetric difference A △ B.
	X := A.Xor(B)

	// Compute A △ B as (A ∖ B) ∪ (B ∖ A).
	Y := A.AndNot(B).Or(B.AndNot(A))

	// Compare the results.
	if X.Equal(Y) {
		fmt.Println(X)
	}

	// Compute A ∩ B
	Z := A.And(B)
	fmt.Println(Z)
	// Output: {1..49 100..149 200}
	// {0 50..99}
}

// Create the set of all primes less than n in O(n log log n) time.
// Try the code with n equal to a few hundred millions and be pleasantly surprised.
func Example_primes() {
	// Sieve of Eratosthenes
	const n = 50
	sieve := bit.New().AddRange(2, n)
	sqrtN := int(math.Sqrt(n))
	for p := 2; p <= sqrtN; p = sieve.Next(p) {
		for k := p * p; k < n; k += p {
			sieve.Delete(k)
		}
	}
	fmt.Println(sieve)
	// Output: {2 3 5 7 11 13 17 19 23 29 31 37 41 43 47}
}

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

func ExampleNew() {
	fmt.Println(bit.New(0, 1, 10, 10, -1))
	// Output: {0 1 10}
}

func ExampleSet_String() {
	fmt.Println(bit.New(1, 2, 6, 5, 3))
	// Output: {1..3 5 6}
}
