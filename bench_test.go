package bit

import "testing"

// Number of words in test set.
const nw = 1 << 10

func BenchmarkSize(b *testing.B) {
	s := BuildTestSet(nw << 3) // Allocates nw<<3 bytes = nw words.
	b.ResetTimer()
	for i := 0; i < b.N/nw; i++ { // Measure time per word.
		s.Size()
	}
}

func BenchmarkNext(b *testing.B) {
	s := BuildTestSet(b.N)
	b.ResetTimer()
	for n := -2; n != -1; {
		n = s.Next(n)
	}
}

func BenchmarkPrev(b *testing.B) {
	s := BuildTestSet(b.N)
	b.ResetTimer()
	for n := MaxInt; n != -1; {
		n = s.Prev(n)
	}
}

func BenchmarkVisit(b *testing.B) {
	s := BuildTestSet(b.N) // As Visit is pretty fast, s can be pretty big.
	b.ResetTimer()
	s.Visit(func(n int) (skip bool) { return })
}

func BenchmarkSetAnd(b *testing.B) {
	s := New(64*nw - 1).Delete(64*nw - 1) // Allocates nw words.
	s1 := BuildTestSet(nw << 3)
	s2 := BuildTestSet(nw << 3)
	b.ResetTimer()
	for i := 0; i < b.N/nw; i++ { // Measure time per word.
		s.SetAnd(s1, s2)
	}
}

func BenchmarkString(b *testing.B) {
	s := BuildTestSet(b.N) // As Visit is pretty fast, s can be pretty big.
	b.ResetTimer()
	_ = s.String()
}

// Quickly builds a set of n somewhat random elements from 0..8n-1.
func BuildTestSet(n int) *Set {
	s := New()
	lfsr := uint16(0xace1) // linear feedback shift register
	for i := 0; i < n; i++ {
		bit := (lfsr>>0 ^ lfsr>>2 ^ lfsr>>3 ^ lfsr>>5) & 1
		lfsr = lfsr>>1 | bit<<15
		e := i<<3 + int(lfsr&0x7)
		s.Add(e) // Add a number from 8i..8i+7.
	}
	return s
}
