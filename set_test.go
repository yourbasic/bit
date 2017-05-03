package bit

import (
	"reflect"
	"strconv"
	"testing"
)

// CheckInvariants checks that the invariants for s.data hold.
func CheckInvariants(t *testing.T, msg string, s *Set) {
	len := len(s.data)
	cap := cap(s.data)
	data := s.data[:cap]
	m := "Invariant for "
	if len > 0 && data[len-1] == 0 {
		t.Errorf("%s%s: data = %v, data[%d] = 0; want non-zero", m, msg, data, len-1)
	}
	for i := len; i < cap; i++ {
		if data[i] != 0 {
			t.Errorf("%s%s: data = %v, data[%d] = %#x; want 0", m, msg, data, i, data[i])
			break
		}
	}
}

// Panics tells if function f panics with parameters p.
func Panics(f interface{}, p ...interface{}) bool {
	fv := reflect.ValueOf(f)
	ft := reflect.TypeOf(f)
	if ft.NumIn() != len(p) {
		panic("wrong argument count")
	}
	pv := make([]reflect.Value, len(p))
	for i, v := range p {
		if reflect.TypeOf(v) != ft.In(i) {
			panic("wrong argument type")
		}
		pv[i] = reflect.ValueOf(v)
	}
	return call(fv, pv)
}

func call(fv reflect.Value, pv []reflect.Value) (b bool) {
	defer func() {
		if err := recover(); err != nil {
			b = true
		}
	}()
	fv.Call(pv)
	return
}

func TestNew(t *testing.T) {
	for _, s := range []*Set{
		New(),
		New(-1),
		New(1),
		New(1, 1),
		New(65),
		New(1, 2, 3),
		New(100, 200, 300),
	} {
		CheckInvariants(t, "New", s)
	}
}

func TestContains(t *testing.T) {
	for _, x := range []struct {
		s        *Set
		n        int
		contains bool
	}{
		{New(), -1, false},
		{New(), 1, false},
		{New(), 100, false},
		{New(-1), 1, false},
		{New(-1), -1, false},
		{New(0), 0, true},
		{New(1), 0, false},
		{New(1), 1, true},
		{New(1), 100, false},
		{New(65), 0, false},
		{New(65), 1, false},
		{New(65), 65, true},
		{New(65), 100, false},

		{New(1, 2, 3), 0, false},
		{New(1, 2, 3), 1, true},
		{New(1, 2, 3), 2, true},
		{New(1, 2, 3), 3, true},
		{New(1, 2, 3), 4, false},

		{New(100, 200, 300), 0, false},
		{New(100, 200, 300), 100, true},
		{New(100, 200, 300), 200, true},
		{New(100, 200, 300), 300, true},
		{New(100, 200, 300), 400, false},
	} {
		s, n := x.s, x.n
		contains := s.Contains(n)
		if contains != x.contains {
			t.Errorf("%v.Contains(%d) = %t; want %t", s, n, contains, x.contains)
		}
	}
}

func TestEqual(t *testing.T) {
	s1, s2 := New(), New()
	if !s1.Equal(s1) {
		t.Errorf("Equal not equal to self.")
	}
	if !s1.Equal(s2) {
		t.Errorf("%v.Equal(%v) false; want true.", s1, s2)
	}
	s1 = New().AddRange(1, 100)
	if s1.Equal(s2) {
		t.Errorf("%v.Equal(%v) true; want false.", s1, s2)
	}
	s2 = New().AddRange(1, 100)
	if !s1.Equal(s2) {
		t.Errorf("%v.Equal(%v) false; want true.", s1, s2)
	}
	s2 = s2.Delete(65)
	if s1.Equal(s2) {
		t.Errorf("%v.Equal(%v) true; want false.", s1, s2)
	}
}

func TestMax(t *testing.T) {
	for _, x := range []struct {
		s   *Set
		max int
	}{
		{New(0), 0},
		{New(65), 65},
		{New(1, 2, 3), 3},
		{New(100, 200, 300), 300},
	} {
		s := x.s
		max := s.Max()
		if max != x.max {
			t.Errorf("%v.Max() = %d; want %d", s, max, x.max)
		}
	}

	s := New()
	if !Panics((*Set).Max, s) {
		t.Errorf("Max() should panic for empty set.")
	}
	CheckInvariants(t, "Max() panic", s)
}

func TestSize(t *testing.T) {
	for _, x := range []struct {
		s    *Set
		size int
	}{
		{New(), 0},
		{New(-1), 0},
		{New(1), 1},
		{New(64), 1},
		{New(65), 1},
		{New(1, 2, 3), 3},
		{New(100, 200, 300), 3},
		{New().AddRange(0, 64), 64},
		{New().AddRange(1, 64), 63},
		{New().AddRange(0, 63), 63},
	} {
		s := x.s
		size := s.Size()
		if size != x.size {
			t.Errorf("%v.Size() = %d; want %d", s, size, x.size)
		}
	}
}

func TestEmpty(t *testing.T) {
	for _, x := range []struct {
		s     *Set
		empty bool
	}{
		{New(), true},
		{New(-1), true},
		{New(1), false},
		{New(65), false},
		{New(1, 2, 3), false},
		{New(100, 200, 300), false},
	} {
		s := x.s
		empty := s.Empty()
		if empty != x.empty {
			t.Errorf("%v.Empty() = %v; want %v", s, empty, x.empty)
		}
	}
}

func TestVisit(t *testing.T) {
	for _, x := range []struct {
		s   *Set
		res string
	}{
		{New(), ""},
		{New(0), "0"},
		{New(1, 2, 3, 62, 63, 64), "123626364"},
		{New(1, 22, 333, 4444), "1223334444"},
	} {
		s := x.s
		res := ""

		s.Visit(func(n int) (skip bool) {
			res += strconv.Itoa(n)
			return
		})
		if res != x.res {
			t.Errorf("%v.Visit(func(n int) { s += Itoa(n) }) -> s=%q; want %q", s, res, x.res)
		}

		s = x.s
		res = ""
		s.Visit(func(n int) (skip bool) {
			s.DeleteRange(0, n+1)
			res += strconv.Itoa(n)
			return
		})
		if res != x.res {
			t.Errorf("%v.Visit(func(n int) { s.DeleteRange(0, n+1); s += Itoa(n) }) -> s=%q; want %q", s, res, x.res)
		}
	}
	s := New(1, 2)
	count := 0
	aborted := s.Visit(func(n int) (skip bool) {
		count++
		if n == 1 {
			skip = true
			return
		}
		return
	})
	if aborted == false {
		t.Errorf("Visit returned false when aborted.")
	}
	if count > 1 {
		t.Errorf("Visit didn't abort.")
	}
	count = 0
	aborted = s.Visit(func(n int) (skip bool) {
		count++
		return
	})
	if aborted == true {
		t.Errorf("Visit returned true when not aborted.")
	}
	if count != 2 {
		t.Errorf("Visit aborted.")
	}
}

func TestString(t *testing.T) {
	for _, x := range []struct {
		s   *Set
		res string
	}{
		{New(), "{}"},
		{New(-1), "{}"},
		{New(1), "{1}"},
		{New(1, -1), "{1}"},
		{New(1, 2), "{1 2}"},
		{New(1, 3), "{1 3}"},
		{New(0, 2, 3), "{0 2 3}"},
		{New(0, 1, 3), "{0 1 3}"},
		{New(0, 2, 3, 5), "{0 2 3 5}"},
		{New(0, 1, 2, 4, 5), "{0..2 4 5}"},
		{New(0, 1, 2, 3, 5, 7, 8, 9), "{0..3 5 7..9}"},
		{New(65), "{65}"},
		{New(100, 200, 300), "{100 200 300}"},
	} {
		res := x.s.String()
		if res != x.res {
			t.Errorf("S.String() = %q; want %q", res, x.res)
		}
	}
}

func TestAdd(t *testing.T) {
	for _, x := range []struct {
		s   *Set
		res string
	}{
		{New().Add(-1), "{}"},
		{New().Add(1), "{1}"},
		{New(1).Add(1), "{1}"},
		{New(1).Add(2), "{1 2}"},
		{New().Add(65), "{65}"},
		{New().Add(100).Add(200).Add(300), "{100 200 300}"},
	} {
		res := x.s.String()
		if res != x.res {
			t.Errorf("s.Add() = %q; want %q", res, x.res)
		}
		CheckInvariants(t, "Add", x.s)
	}
}

func TestDelete(t *testing.T) {
	for _, x := range []struct {
		s   *Set
		res string
	}{
		{New(1).Delete(1), "{}"},
		{New(1).Delete(-1), "{1}"},
		{New(1).Delete(2), "{1}"},
		{New(65).Delete(64), "{65}"},
		{New(100, 200, 300).Delete(200), "{100 300}"},
		{New(100, 200, 300).Delete(300), "{100 200}"},
	} {
		res := x.s.String()
		if res != x.res {
			t.Errorf("s.Delete() = %q; want %q", res, x.res)
		}
		CheckInvariants(t, "Delete", x.s)
	}
}

type rangeFunc struct {
	fInt   func(S *Set, n int) *Set
	fRange func(S *Set, m, n int) *Set
	name   string
}

func TestRange(t *testing.T) {
	rangeFuncs := []rangeFunc{
		{(*Set).Add, (*Set).AddRange, "AddRange"},
		{(*Set).Delete, (*Set).DeleteRange, "DeleteRange"},
	}

	for _, x := range []struct {
		s    *Set
		m, n int
	}{
		{New(), 0, 0},
		{New(), 2, 1},
		{New(), -2, -1},
		{New(), -1, 0},
		{New(), -1, -1},
		{New(), 1, 10},
		{New(), 64, 66},
		{New(), 1, 1000},

		{New(1, 2, 3), 0, 1},
		{New(1, 2, 3), 0, 2},
		{New(1, 2, 3), 0, 3},
		{New(1, 2, 3), 0, 4},
		{New(1, 2, 3), 1, 2},
		{New(1, 2, 3), 1, 4},
		{New(1, 2, 3), 1, 5},
		{New(1, 2, 3), 1, 1000},

		{New(100, 200, 300), 50, 100},
		{New(100, 200, 300), 50, 101},
		{New(100, 200, 300), 50, 250},
		{New(100, 200, 300), 50, 350},
		{New(100, 200, 300), 250, 350},
		{New(100, 200, 300), 300, 350},
		{New(100, 200, 300), 350, 400},
		{New(100, 200, 300), 1, 1000},
	} {
		for _, o := range rangeFuncs {
			fInt, fRange, name := o.fInt, o.fRange, o.name
			s := x.s
			m, n := x.m, x.n

			res := fRange(new(Set).Set(s), m, n)
			exp := new(Set).Set(s)
			for i := m; i < n; i++ {
				fInt(exp, i)
			}
			if !res.Equal(exp) {
				t.Errorf("%v.%v(%d, %d) = %v; want %v", s, name, m, n, res, exp)
			}
			CheckInvariants(t, name, res)
		}
	}
}

func TestSet(t *testing.T) {
	for _, x := range []struct {
		s, a *Set
	}{
		{New(), New()},
		{New(), New(1)},
		{New(), New(65)},
		{New(), New(1, 2, 3)},
		{New(), New(100, 200, 300)},

		{New(1, 2, 3), New()},
		{New(1, 2, 3), New(1)},
		{New(1, 2, 3), New(65)},
		{New(1, 2, 3), New(1, 2, 3)},
		{New(1, 2, 3), New(100, 200, 300)},

		{New(100, 200, 300), New()},
		{New(100, 200, 300), New(1)},
		{New(100, 200, 300), New(65)},
		{New(100, 200, 300), New(1, 2, 3)},
		{New(100, 300, 300), New(100, 200, 300)},
	} {
		s := x.s

		ss := s.Set(x.a)
		if ss != s {
			t.Errorf("&(s.set(%v)) = %p, &S = %p; want same", x.a, ss, s)
		}
		if !ss.Equal(x.a) {
			t.Errorf("s.set(%v) = %v; want %v", x.a, ss, x.a)
		}
		CheckInvariants(t, "set", ss)
	}
}

type binOp struct {
	f    func(s *Set, a, b *Set) *Set
	name string
}

func TestBinOp(t *testing.T) {
	And := binOp{(*Set).SetAnd, "SetAnd"}
	AndNot := binOp{(*Set).SetAndNot, "SetAndNot"}
	Or := binOp{(*Set).SetOr, "SetOr"}
	for _, x := range []struct {
		op   binOp
		a, b *Set
		exp  *Set
	}{
		{And, New(), New(), New()},
		{And, New(1), New(), New()},
		{And, New(), New(1), New()},
		{And, New(1), New(1), New(1)},
		{And, New(1), New(2), New()},
		{And, New(1), New(1, 2), New(1)},
		{And, New(1, 2), New(2, 3), New(2)},
		{And, New(100), New(), New()},
		{And, New(), New(100), New()},
		{And, New(100), New(100), New(100)},
		{And, New(100), New(100, 200), New(100)},
		{And, New(200), New(100, 200), New(200)},
		{And, New(100, 200), New(200, 300), New(200)},

		{AndNot, New(), New(), New()},
		{AndNot, New(1), New(), New(1)},
		{AndNot, New(), New(1), New()},
		{AndNot, New(1), New(1), New()},
		{AndNot, New(1), New(2), New(1)},
		{AndNot, New(1), New(1, 2), New()},
		{AndNot, New(1, 2), New(2, 3), New(1)},
		{AndNot, New(100), New(), New(100)},
		{AndNot, New(), New(100), New()},
		{AndNot, New(100), New(100), New()},
		{AndNot, New(100), New(100, 200), New()},
		{AndNot, New(200), New(100, 200), New()},
		{AndNot, New(100, 200), New(200, 300), New(100)},

		{Or, New(), New(), New()},
		{Or, New(), New(1), New(1)},
		{Or, New(1), New(), New(1)},
		{Or, New(1), New(1), New(1)},
		{Or, New(1), New(2), New(1, 2)},
		{Or, New(1), New(1, 2), New(1, 2)},
		{Or, New(1, 2), New(2, 3), New(1, 2, 3)},
		{Or, New(100), New(), New(100)},
		{Or, New(), New(100), New(100)},
		{Or, New(100), New(100), New(100)},
		{Or, New(100), New(100, 200), New(100, 200)},
		{Or, New(200), New(100, 200), New(100, 200)},
		{Or, New(100, 200), New(200, 300), New(100, 200, 300)},
	} {
		op, name := x.op.f, x.op.name
		a, b := New().Set(x.a), New().Set(x.b)
		s := New()

		res := op(s, a, b)
		exp := x.exp
		if s != res {
			t.Errorf("&(s.%s(%v, %v)) = %p &s = %p; want same", name, a, b, s, res)
		}
		if !res.Equal(exp) {
			t.Errorf("s.%s(%v, %v) = %v; want %v", name, x.a, x.b, res, exp)
		}
		CheckInvariants(t, name, res)

		a.Set(x.a)
		b.Set(x.b)
		s = a
		res = op(s, a, b)
		if !res.Equal(exp) {
			t.Errorf("s.%s(%v, %v) = %v; want %v", name, x.a, x.b, res, exp)
		}
		CheckInvariants(t, name, res)

		a.Set(x.a)
		b.Set(x.b)
		s = b
		res = op(s, a, b)
		if !res.Equal(exp) {
			t.Errorf("s.%s(%v, %v) = %v; want %v", name, x.a, x.b, res, exp)
		}
		CheckInvariants(t, name, res)

		a.Set(x.a)
		b.Set(x.b)
		s = New().AddRange(150, 250)
		res = op(s, a, b)
		if !res.Equal(exp) {
			t.Errorf("s.%s(%v, %v) = %v; want %v", name, x.a, x.b, res, exp)
		}
		CheckInvariants(t, name, res)
	}
}

func TestNextPow2(t *testing.T) {
	for _, x := range []struct {
		n, p int
	}{
		{MinInt, 1},
		{-1, 1},
		{0, 1},
		{1, 2},
		{2, 4},
		{3, 4},
		{4, 8},
		{1<<19 - 1, 1 << 19},
		{1 << 19, 1 << 20},
		{MaxInt >> 1, MaxInt>>1 + 1},
		{MaxInt>>1 + 1, MaxInt},
		{MaxInt - 1, MaxInt},
		{MaxInt, MaxInt},
	} {
		n := x.n

		p := nextPow2(n)
		if p != x.p {
			t.Errorf("nextPow2(%#x) = %#x; want %#x", n, p, x.p)
		}
	}
}
