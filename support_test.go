package timeserie_test

import (
	"math"
	"testing"
	"time"

	"github.com/etnz/timeserie"
)

// TestSupport_Len basically append data and verify that len is correct.
func TestSupport_Len(t *testing.T) {
	s := new(timeserie.Support)
	if s.Len() != 0 {
		t.Errorf("{}.Len() =%v want 0", s.Len())
	}

	s.Append(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), 1.0)
	if s.Len() != 1 {
		t.Errorf("{'2000-1-1', 1.0}.Len() =%v want 1", s.Len())
	}

	s.Append(time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC), 2.0)
	if s.Len() != 2 {
		t.Errorf("{'2000-1-1', 1.0; '2000-1-2', 2.0}.Len() =%v want 2", s.Len())
	}

}

// TestSupport_Append basically append data and verify that len is correct and that order
// is correct.
func TestSupport_Append(t *testing.T) {
	s := new(timeserie.Support)
	if s.Len() != 0 {
		t.Errorf("{}.Len() =%v want 0", s.Len())
	}

	s.Append(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), 1.0)
	if s.Len() != 1 {
		t.Errorf("{'2000-1-1', 1.0}.Len() =%v want 1", s.Len())
	}

	s.Append(time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC), 2.0)
	if s.Len() != 2 {
		t.Errorf("{'2000-1-1', 1.0; '2000-1-2', 2.0}.Len() =%v want 2", s.Len())
	}

	s.Append(time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC), math.NaN())
	if s.Len() != 2 {
		t.Errorf("{'2000-1-1', 1.0; '2000-1-2', 2.0, xxx,NaN}.Len() =%v want 2", s.Len())
	}

	d := time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)
	s.Append(d, 3.0)
	if s.Len() != 3 {
		t.Errorf("{'1999-1-1', 3.0 ; '2000-1-1', 1.0; '2000-1-2', 2.0}.Len() =%v want 3", s.Len())
	}
	// also checks that the order is correct
	g, v := s.At(0)
	if g != d || v != 3.0 {
		t.Errorf("{'1999-1-1', 3.0 ; '2000-1-1', 1.0; '2000-1-2', 2.0}.At(0) =%v,%f want %v,%f", g, v, d, 3.0)
	}
}

// TestSupport_Find tries most edge cases on a length 2 timeserie.
func TestSupport_Find(t *testing.T) {
	low := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	high := time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)

	s := new(timeserie.Support)
	s.Append(low, 1.0)
	s.Append(high, 2.0)

	x := s.Find(low.Add(-48 * time.Hour))
	if x != -1 {
		t.Errorf("Find(before low) = %v want -1", x)
	}
	x = s.Find(low.Add(+48 * time.Hour))
	if x != 0 {
		t.Errorf("Find(after low) = %v want 0", x)
	}
	x = s.Find(high.Add(-48 * time.Hour))
	if x != 0 {
		t.Errorf("Find(before high) = %v want 0", x)
	}
	x = s.Find(high.Add(48 * time.Hour))
	if x != 1 {
		t.Errorf("Find(after high) = %v want 1", x)
	}
}

// TestSupport_Values tries a simple case.
func TestSupport_Values(t *testing.T) {
	low := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	high := time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)

	s := new(timeserie.Support)
	s.Append(low, 1.0)
	s.Append(high, 2.0)

	i := 0
	for ti, v := range s.Values() {
		if i == 0 {
			if ti != low || v != 1.0 {
				t.Errorf("Values[0] = %v, %v want %v, 1.0", ti, v, low)
			}
		}
		if i == 1 {
			if ti != high || v != 2.0 {
				t.Errorf("Values[1] = %v, %v want %v, 2.0", ti, v, high)
			}
		}
		i++
	}
}

// TestSupport_Times tries a simple case.
func TestSupport_Times(t *testing.T) {
	low := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	high := time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)

	s := new(timeserie.Support)
	s.Append(low, 1.0)
	s.Append(high, 2.0)

	i := 0
	for ti := range s.Times() {
		if i == 0 {
			if ti != low {
				t.Errorf("Times[0] = %v want %v", ti, low)
			}
		}
		if i == 1 {
			if ti != high {
				t.Errorf("Times[1] = %v want %v", ti, high)
			}
		}
		i++
	}
}

// TestSupport_Delta tries Delta on simple cases
func TestSupport_Delta(t *testing.T) {
	s := new(timeserie.Support)
	d0, d1, d2 := timeserie.DayDate(2000, 1, 1), timeserie.DayDate(2000, 1, 2), timeserie.DayDate(2000, 1, 3)
	s.Append(d0, 1.0)
	s.Append(d1, 2.0)
	s.Append(d2, 4.0)

	x := s.Delta()
	if x.Len() != 2 {
		t.Errorf("{1,2,4}.Delta() Len=%v want 2", x.Len())
	}
	x0, v0 := x.At(0)
	if x0 != d1 || v0 != 1.0 {
		t.Errorf("{1,2,4}.Delta()[0] =%v,%v want %v, %v", x0, v0, d1, 1.0)
	}

	x1, v1 := x.At(1)
	if x1 != d2 || v1 != 2.0 {
		t.Errorf("{1,2,4}.Delta()[1] =%v,%v want %v, %v", x1, v1, d2, 2.0)
	}

}

// TestSupport_Scan test the Scan method on simple case.
func TestSupport_Scan(t *testing.T) {
	s := new(timeserie.Support)
	d0, d1 := timeserie.DayDate(2000, 1, 1), timeserie.DayDate(2000, 1, 2)
	s.Append(d0, 1.0)
	s.Append(d1, 2.0)

	x := s.Scan(0, timeserie.ScannerAcc)
	if x.Len() != 2 {
		t.Errorf("{1,2}.Scan(Acc).Len=%v want 2", x.Len())
	}
	x0, v0 := x.At(0)
	if x0 != d0 || v0 != 1.0 {
		t.Errorf("{1,2}.Scan(Acc)[0] =%v,%v want %v, %v", x0, v0, d0, 1.0)
	}

	x1, v1 := x.At(1)
	if x1 != d1 || v1 != 3.0 {
		t.Errorf("{1,2}.Scan(Acc)[1] =%v,%v want %v, %v", x1, v1, d1, 3.0)
	}
}
