package bit_test

import (
	"fmt"
	"github.com/yourbasic/bit"
)

// Union returns the union of the given sets.
func Union(s ...*bit.Set) *bit.Set {
	// Optimization: allocate initital set with adequate capacity.
	max := -1
	for _, x := range s {
		if x.Size() > 0 && x.Max() > max { // Max is undefined for the empty set.
			max = x.Max()
		}
	}
	res := bit.New(max) // A negative number is not included.

	for _, x := range s {
		res.SetOr(res, x) // Reuses memory.
	}
	return res
}

// Implement a variadic Union function using SetOr.
func Example_union() {
	a, b, c := bit.New(1, 2), bit.New(2, 3), bit.New(5)
	fmt.Println(Union(a, b, c))
	// Output: {1..3 5}
}
