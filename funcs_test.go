package bit

import (
	"testing"
)

func TestConstants(t *testing.T) {
	if bitsPerWord != 32 && bitsPerWord != 64 {
		t.Errorf("bitsPerWord = %v; want 32 or 64", bitsPerWord)
	}
	if bitsPerWord == 32 {
		if MaxUint != 1<<32-1 {
			t.Errorf("MaxUint = %#x; want 1<<32 - 1", uint64(MaxUint))
		}
		if MaxInt != 1<<31-1 {
			t.Errorf("MaxInt = %#x; want 1<<31 - 1", int64(MaxInt))
		}
		if MinInt != -1<<31 {
			t.Errorf("MinInt = %#x; want -1 << 31", int64(MinInt))
		}
	}
	if bitsPerWord == 64 {
		if MaxUint != 1<<64-1 {
			t.Errorf("MaxUint = %#x; want 1<<64 - 1", uint64(MaxUint))
		}
		if MaxInt != 1<<63-1 {
			t.Errorf("MaxInt = %#x; want 1<<63 - 1", int64(MaxInt))
		}
		if MinInt != -1<<63 {
			t.Errorf("MinInt = %#x; want -1 << 63", int64(MinInt))
		}
	}
}

// Checks all words with one nonzero bit.
func TestWordOneBit(t *testing.T) {
	for i := 0; i < 64; i++ {
		var w uint64 = 1 << uint(i)
		lead, trail, count := LeadingZeros(w), TrailingZeros(w), Count(w)
		if lead != 63-i {
			t.Errorf("LeadingZeros(%#x) = %d; want %d", w, lead, i)
		}
		if trail != i {
			t.Errorf("TrailingZeros(%#x) = %d; want %d", w, trail, i)
		}
		if count != 1 {
			t.Errorf("Count(%#x) = %d; want %d", w, count, 1)
		}
	}
}

func TestWordFuncs(t *testing.T) {
	for _, x := range []struct {
		w                  uint64
		lead, trail, count int
	}{
		{0x0, 64, 64, 0},
		{0xa, 60, 1, 2},
		{0xffffffffffffffff, 0, 0, 64},
		{0x7ffffffffffffffe, 1, 1, 62},
		{0x5555555555555555, 1, 0, 32},
		{0xaaaaaaaaaaaaaaaa, 0, 1, 32},
	} {
		w := x.w
		lead, trail, count := LeadingZeros(w), TrailingZeros(w), Count(w)
		if lead != x.lead {
			t.Errorf("LeadingZeros(%#x) = %v; want %v", w, lead, x.lead)
		}
		if trail != x.trail {
			t.Errorf("TrailingZeros(%#x) = %v; want %v", w, trail, x.trail)
		}
		if count != x.count {
			t.Errorf("Count(%#x) = %v; want %v", w, count, x.count)
		}
	}
}
