package timeserie_test

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/etnz/timeserie"
)

var (
	// three dates for easier testing.
	d0 = timeserie.DayDate(2000, 1, 1)
	d1 = timeserie.DayDate(2000, 1, 2)
	d2 = timeserie.DayDate(2000, 1, 3)
)

// TestFunction_F tests the method on some simple cases.
func TestFunction_F(t *testing.T) {
	s := new(timeserie.Support)
	s.Append(d0, 1.0)
	s.Append(d1, 2.0)

	f := timeserie.New(s, timeserie.ModeNullset)

	{
		x := f.F(d0)
		if x != 1.0 {
			t.Errorf("nullset.F(d0)=%v want %v", x, 1.0)
		}
	}

	{
		x := f.F(d1)
		if x != 2.0 {
			t.Errorf("nullset.F(d1)=%v want %v", x, 2.0)
		}
	}

	// any other date should return NaN
	{
		x := f.F(d0.Add(-timeserie.Day))
		if !math.IsNaN(x) {
			t.Errorf("nullset.F(d0-)=%v want %v", x, math.NaN())
		}
	}

	{
		x := f.F(d0.Add(time.Hour))
		if !math.IsNaN(x) {
			t.Errorf("nullset.F(d0+)=%v want %v", x, math.NaN())
		}
	}

	{
		x := f.F(d1.Add(time.Hour))
		if !math.IsNaN(x) {
			t.Errorf("nullset.F(d1+)=%v want %v", x, math.NaN())
		}
	}
	// now with f as ModeStep
	f = timeserie.New(s, timeserie.ModeStep)

	{
		x := f.F(d0)
		if x != 1.0 {
			t.Errorf("nullset.F(d0)=%v want %v", x, 1.0)
		}
	}

	{
		x := f.F(d1)
		if x != 2.0 {
			t.Errorf("nullset.F(d1)=%v want %v", x, 2.0)
		}
	}

	// any other date should return previous value
	{
		x := f.F(d0.Add(-timeserie.Day))
		if x != 0 {
			t.Errorf("nullset.F(d0-)=%v want %v", x, 0.0)
		}
	}

	{
		x := f.F(d0.Add(time.Hour))
		if x != 1.0 {
			t.Errorf("nullset.F(d0+)=%v want %v", x, 1.0)
		}
	}

	{
		x := f.F(d1.Add(time.Hour))
		if x != 2.0 {
			t.Errorf("nullset.F(d1+)=%v want %v", x, 2.0)
		}
	}

}

// TestIterate checks on simple case that the function works.
func TestIterate(t *testing.T) {
	s1 := new(timeserie.Support)
	s1.Append(d0, 1.0)
	s1.Append(d1, 2.0)
	f1 := timeserie.New(s1, timeserie.ModeNullset)

	s2 := new(timeserie.Support)
	s2.Append(d0, 1.0)
	s2.Append(d2, 3.0)
	f2 := timeserie.New(s2, timeserie.ModeNullset)

	x := slices.AppendSeq([]time.Time{}, timeserie.Iterate(f1, f2))

	if len(x) != 3 {
		t.Errorf("Iterate({0,1}, {0,2}).Len() = %v want 3", len(x))
	}
	if x[0] != d0 || x[1] != d1 || x[2] != d2 {
		t.Errorf("Iterate({0,1}, {0,2}) = %v want [%v, %v, %v]", x, d0, d1, d2)
	}
}

// TestAdd checks on simple case that the result is as expected.
func TestAdd(t *testing.T) {
	s1 := new(timeserie.Support)
	s1.Append(d0, 1.0)
	s1.Append(d1, 2.0)

	s2 := new(timeserie.Support)
	s2.Append(d0, 1.0)
	s2.Append(d2, 3.0)

	{
		x := timeserie.Add(timeserie.New(s1, timeserie.ModeNullset), timeserie.New(s2, timeserie.ModeNullset))
		if x.Len() != 1 {
			t.Errorf("{1,2}+{1,3}.Len() =%v want %v", x.Len(), 1)
		}
		{
			d, val := x.At(0)
			if d != d0 || val != 2.0 {
				t.Errorf("{1,2}+{1,3}[0]=%v,%v want %v,%v", d, val, d0, 2.0)
			}
		}
	}
}
