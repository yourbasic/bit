// +build go1.10

// Package bit provides a bit array implementation.
//
// Bit set
//
// A bit set, or bit array, is an efficient set data structure
// that consists of an array of 64-bit words. Because it uses
// bit-level parallelism, limits memory access, and efficiently uses
// the data cache, a bit set often outperforms other data structures.
//
// Tutorial
//
// The Basics example shows how to create, combine, compare and
// print bit sets.
//
// Primes contains a short and simple, but still efficient,
// implementation of a prime number sieve.
//
// Union is a more advanced example demonstrating how to build
// an efficient variadic Union function using the SetOr method.
//
package bit

import (
	"fmt"
	"math/bits"
	"strings"
)

const (
	bpw   = 64         // bits per word
	maxw  = 1<<bpw - 1 // maximum value of a word
	shift = 6
	mask  = 0x3f
)

// Set represents a mutable set of non-negative integers.
// The zero value is an empty set ready to use.
// A set occupies approximately n bits, where n is the maximum value
// that has been stored in the set.
type Set struct {
	// Invariants:
	//   • data[n>>shift] & (1<<(n&mask)) == 1 iff n belongs to set,
	//   • data[len(data)-1] != 0 if set is nonempty,
	//   • data[i] == 0 for all i such that len(data) ≤ i < cap(data).
	data []uint64
}

// New creates a new set with the given elements.
// Negative numbers are not included in the set.
func New(n ...int) *Set {
	if len(n) == 0 {
		return new(Set)
	}
	max := n[0]
	for _, e := range n {
		if e > max {
			max = e
		}
	}
	if max < 0 {
		return new(Set)
	}
	s := &Set{
		data: make([]uint64, max>>shift+1),
	}
	for _, e := range n {
		if e >= 0 {
			s.data[e>>shift] |= 1 << uint(e&mask)
		}
	}
	return s
}

// Contains tells if n is an element of the set.
func (s *Set) Contains(n int) bool {
	if n < 0 {
		return false
	}
	d := s.data
	i := n >> shift
	if i >= len(d) {
		return false
	}
	return d[i]&(1<<uint(n&mask)) != 0
}

// Equal tells if s1 and s2 contain the same elements.
func (s1 *Set) Equal(s2 *Set) bool {
	if s1 == s2 {
		return true
	}
	a, b := s1.data, s2.data
	la := len(a)
	if la != len(b) {
		return false
	}
	for i := 0; i < la; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Subset tells if s1 is a subset of s2.
func (s1 *Set) Subset(s2 *Set) bool {
	if s1 == s2 {
		return true
	}
	a, b := s1.data, s2.data
	la := len(a)
	if la > len(b) {
		return false
	}
	for i := 0; i < la; i++ {
		if a[i]&^b[i] != 0 {
			return false
		}
	}
	return true
}

// Max returns the maximum element of the set;
// it panics if the set is empty.
func (s *Set) Max() int {
	if len(s.data) == 0 {
		panic("max not defined for empty set")
	}
	d := s.data
	i := len(d) - 1
	return i<<shift + bits.Len64(d[i]) - 1
}

// Size returns the number of elements in the set.
// This method scans the set; to check if a set is empty,
// consider using the more efficient Empty method.
func (s *Set) Size() int {
	d := s.data
	n := 0
	for i, len := 0, len(d); i < len; i++ {
		if w := d[i]; w != 0 {
			n += bits.OnesCount64(w)
		}
	}
	return n
}

// Empty tells if the set is empty.
func (s *Set) Empty() bool {
	return len(s.data) == 0
}

// Next returns the next element n, n > m, in the set,
// or -1 if there is no such element.
func (s *Set) Next(m int) int {
	d := s.data
	len := len(d)
	if len == 0 {
		return -1
	}
	if m < 0 {
		if d[0]&1 != 0 {
			return 0
		}
		m = 0
	}
	i := m >> shift
	if i >= len {
		return -1
	}
	t := 1 + uint(m&mask)
	w := d[i] >> t << t // Zero out bits for numbers ≤ m.
	for i < len-1 && w == 0 {
		i++
		w = d[i]
	}
	if w == 0 {
		return -1
	}
	return i<<shift + bits.TrailingZeros64(w)
}

// Prev returns the previous element n, n < m, in the set,
// or -1 if there is no such element.
func (s *Set) Prev(m int) int {
	d := s.data
	len := len(d)
	if len == 0 || m <= 0 {
		return -1
	}
	i := len - 1
	if max := i<<shift + bits.Len64(d[i]) - 1; m > max {
		return max
	}
	i = m >> shift
	t := bpw - uint(m&mask)
	w := d[i] << t >> t // Zero out bits for numbers ≥ m.
	for i > 0 && w == 0 {
		i--
		w = d[i]
	}
	if w == 0 {
		return -1
	}
	return i<<shift + bits.Len64(w) - 1
}

// Visit calls the do function for each element of s in numerical order.
// If do returns true, Visit returns immediately, skipping any remaining
// elements, and returns true. It is safe for do to add or delete
// elements e, e ≤ n. The behavior of Visit is undefined if do changes
// the set in any other way.
func (s *Set) Visit(do func(n int) (skip bool)) (aborted bool) {
	d := s.data
	for i, len := 0, len(d); i < len; i++ {
		w := d[i]
		if w == 0 {
			continue
		}
		n := i << shift // element represented by w&1
		for w != 0 {
			b := bits.TrailingZeros64(w)
			n += b
			if do(n) {
				return true
			}
			n++
			w >>= uint(b + 1)
			for w&1 != 0 { // common case
				if do(n) {
					return true
				}
				n++
				w >>= 1
			}
		}
	}
	return false
}

// String returns a string representation of the set. The elements
// are listed in ascending order. Runs of at least three consecutive
// elements from a to b are given as a..b.
func (s *Set) String() string {
	buf := new(strings.Builder)
	buf.WriteByte('{')
	a, b := -1, -2 // Keep track of a range a..b of elements.
	first := true
	s.Visit(func(n int) (skip bool) {
		if n == b+1 {
			b++ // Increase current range from a..b to a..b+1.
			return
		}
		if first && a <= b {
			first = false
		} else if a <= b {
			buf.WriteByte(' ')
		}
		writeRange(buf, a, b)
		a, b = n, n // Start new range.
		return
	})
	if !first && a <= b {
		buf.WriteByte(' ')
	}
	writeRange(buf, a, b)
	buf.WriteByte('}')
	return buf.String()
}

// writeRange appends either "", "a", "a b" or "a..b" to buf.
func writeRange(buf *strings.Builder, a, b int) {
	switch {
	case a > b:
		return // Append nothing.
	case a == b:
		fmt.Fprintf(buf, "%d", a)
	case a+1 == b:
		fmt.Fprintf(buf, "%d %d", a, b)
	default:
		fmt.Fprintf(buf, "%d..%d", a, b)
	}
}

// Add adds n to s and returns a pointer to the updated set.
// A negative n will not be added.
func (s *Set) Add(n int) *Set {
	if n < 0 {
		return s
	}
	i := n >> shift
	if i >= len(s.data) {
		s.resize(i + 1)
	}
	s.data[i] |= 1 << uint(n&mask)
	return s
}

// Delete removes n from s and returns a pointer to the updated set.
func (s *Set) Delete(n int) *Set {
	if n < 0 {
		return s
	}
	i := n >> shift
	if i >= len(s.data) {
		return s
	}
	s.data[i] &^= 1 << uint(n&mask)
	s.trim()
	return s
}

// AddRange adds all integers from m to n-1 to s
// and returns a pointer to the updated set.
// Negative numbers will not be added.
func (s *Set) AddRange(m, n int) *Set {
	if n < 1 || m >= n {
		return s
	}
	m = max(0, m)
	n--
	low, high := m>>shift, n>>shift
	if high >= len(s.data) {
		s.resize(high + 1)
	}
	d := s.data
	// Range fits in one word.
	if low == high {
		d[low] |= bitMask(m&mask, n&mask)
		return s
	}
	// Range spans at least two words.
	d[low] |= bitMask(m&mask, bpw-1)
	for i := low + 1; i < high; i++ {
		d[i] = maxw
	}
	d[high] |= bitMask(0, n&mask)
	return s
}

// DeleteRange removes all integers from m to n-1 from s
// and returns a pointer to the updated set.
func (s *Set) DeleteRange(m, n int) *Set {
	if n < 1 || m >= n {
		return s
	}
	m = max(0, m)
	n--
	d := s.data
	low, high := m>>shift, n>>shift
	// Range does not intersect set.
	if low >= len(d) {
		return s
	}
	// Top of range overshoots set.
	if len(d) <= high {
		high = len(d) - 1 // low ≤ high still holds, since low < len(d).
		n = bpw - 1       // To assure that n&mask == bpw-1 below.
	}
	// Range fits in one word.
	if low == high {
		d[low] &^= bitMask(m&mask, n&mask)
		s.trim()
		return s
	}
	// Range spans at least two words.
	d[low] &^= bitMask(m&mask, bpw-1)
	for i := low + 1; i < high; i++ {
		d[i] = 0
	}
	d[high] &^= bitMask(0, n&mask)
	s.trim()
	return s
}

// And creates a new set that consists of all elements that belong
// to both s1 and s2.
func (s1 *Set) And(s2 *Set) *Set {
	return new(Set).SetAnd(s1, s2)
}

// Or creates a new set that contains all elements that belong
// to either s1 or s2.
func (s1 *Set) Or(s2 *Set) *Set {
	return new(Set).SetOr(s1, s2)
}

// Xor creates a new set that contains all elements that belong
// to either s1 or s2, but not to both.
func (s1 *Set) Xor(s2 *Set) *Set {
	return new(Set).SetXor(s1, s2)
}

// AndNot creates a new set that consists of all elements that belong
// to s1, but not to s2.
func (s1 *Set) AndNot(s2 *Set) *Set {
	return new(Set).SetAndNot(s1, s2)
}

// Set sets s to s1 and then returns a pointer to the updated set s.
func (s *Set) Set(s1 *Set) *Set {
	s.realloc(len(s1.data))
	copy(s.data, s1.data)
	return s
}

// SetAnd sets s to the intersection s1 ∩ s2 and then returns a pointer to s.
func (s *Set) SetAnd(s1, s2 *Set) *Set {
	a, b := s1.data, s2.data
	// Find last nonzero word in result.
	n := min(len(a), len(b)) - 1
	for n >= 0 && a[n]&b[n] == 0 {
		n--
	}
	if s == s1 || s == s2 {
		s.resize(n + 1)
	} else {
		s.realloc(n + 1)
	}
	for i := 0; i <= n; i++ {
		s.data[i] = a[i] & b[i]
	}
	return s
}

// SetAndNot sets s to the set difference s1 ∖ s2 and then returns a pointer to s.
func (s *Set) SetAndNot(s1, s2 *Set) *Set {
	a, b := s1.data, s2.data
	la, lb := len(a), len(b)
	// Result requires len(a) words if len(a) > len(b),
	// otherwise find last nonzero word in result.
	n := la - 1
	if la <= lb {
		for n >= 0 && a[n]&^b[n] == 0 {
			n--
		}
	}
	if s == s1 || s == s2 {
		s.resize(n + 1)
	} else {
		s.realloc(n + 1)
	}
	d := s.data
	if m := lb; m <= n {
		copy(d[m:n+1], a[m:n+1])
		n = m - 1
	}
	for i := 0; i <= n; i++ {
		d[i] = a[i] &^ b[i]
	}
	return s
}

// SetOr sets s to the union s1 ∪ s2 and then returns a pointer to s.
func (s *Set) SetOr(s1, s2 *Set) *Set {
	// Swap, if necessary, to make s1 shorter than s2.
	if len(s1.data) > len(s2.data) {
		s1, s2 = s2, s1
	}
	a, b := s1.data, s2.data
	la := len(a)
	n := len(b) - 1
	if s == s1 || s == s2 {
		s.resize(n + 1)
	} else {
		s.realloc(n + 1)
	}
	d := s.data
	copy(d[la:n+1], b[la:n+1])
	for i := 0; i < la; i++ {
		d[i] = a[i] | b[i]
	}
	return s
}

// SetXor sets s to the  symmetric difference A ∆ B = (A ∪ B) ∖ (A ∩ B)
// and then returns a pointer to s.
func (s *Set) SetXor(s1, s2 *Set) *Set {
	// Swap, if necessary, to make s1 shorter than s2.
	if len(s1.data) > len(s2.data) {
		s1, s2 = s2, s1
	}
	a, b := s1.data, s2.data
	la, lb := len(a), len(b)
	n := lb - 1
	if la == lb { // The only case where result may be shorter than len(b).
		for n >= 0 && a[n]^b[n] == 0 {
			n--
		}
		if n == -1 { // No elements left.
			s.realloc(0)
			return s
		}
	}
	if s == s1 || s == s2 {
		s.resize(n + 1)
	} else {
		s.realloc(n + 1)
	}
	d := s.data
	if la <= n {
		copy(d[la:n+1], b[la:n+1])
		n = la - 1
	}
	for i := 0; i <= n; i++ {
		d[i] = a[i] ^ b[i]
	}
	return s
}

// resize changes the length of s.data to n, keeping old values.
// It preserves the invariant s.data[i] = 0, n ≤ i < cap(data).
func (s *Set) resize(n int) {
	d := s.data
	if s.realloc(n) {
		copy(s.data, d)
	}
}

// realloc creates a slice s.data of length n, possibly zeroing out old values.
// It preserves the invariant s.data[i] = 0, n ≤ i < cap(data).
// It returns true if new memory has been allocated.
func (s *Set) realloc(n int) (didAlloc bool) {
	if c := cap(s.data); c < n {
		s.data = make([]uint64, n, newCap(n, c))
		return true
	}
	// Add zeroes if shrinking.
	d := s.data
	for i := len(d) - 1; i >= n; i-- {
		d[i] = 0
	}
	s.data = d[:n]
	return false
}

// newCap suggests a new increased capacity, favoring powers of two,
// when growing a slice to length n. The suggested capacities guarantee
// linear amortized cost for repeated memory allocations.
func newCap(n, prevCap int) int {
	return max(n, nextPow2(prevCap))
}

// nextPow2 returns the smallest p = 1, 2, 4, ..., 2^k such that p > n,
// or MaxInt if p > MaxInt.
func nextPow2(n int) (p int) {
	if n <= 0 {
		return 1
	}
	if k := bits.Len64(uint64(n)); k < bitsPerWord-1 {
		return 1 << uint(k)
	}
	return MaxInt
}

// trim slices s.data by removing all trailing words equal to zero.
func (s *Set) trim() {
	d := s.data
	n := len(d) - 1
	for n >= 0 && d[n] == 0 {
		n--
	}
	s.data = d[:n+1]
}

// bitMask returns a bit mask with nonzero bits from m to n, 0 ≤ m ≤ n < bpw.
func bitMask(m, n int) uint64 {
	return maxw >> uint(bpw-1-(n-m)) << uint(m)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
