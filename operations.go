package timeserie

import (
	"iter"
	"math"
	"slices"
	"time"
)

type Mode int

const (
	ModeNullset Mode = iota // function defined only on the support, NaN everywhere else.
	ModeStep                // Function's value between two support events is the value of the earliest.
	//ModeLinear              // Function's value between two support events is a linear interpolation between the two.
	LenMode // not a mode but the length of modes
)

// Function is the interface of all support-based functions.
type Function struct {
	Support
	mode Mode
}

// New creates a new function defined by its support and the interpolation mode.
func New(s *Support, mode Mode) *Function { return &Function{Support: *s, mode: mode} }

// F returns the function value at a given time. If not defined on that time, it returns NaN.
func (f *Function) F(t time.Time) float64 {
	switch f.mode {
	case ModeNullset:
		i := slices.Index(f.times, t)
		if i >= 0 {
			return f.values[i]
		} else {
			return math.NaN()
		}
	case ModeStep:
		// // below block is useless but I'll wait for test to remove it.
		// if f.Len() == 0 { // but it is not nullset
		// 	// no data
		// 	return 0.0 // in ModeStep value at -inf is 0
		// }
		if prev := f.Find(t); prev < 0 {
			return 0.0
		} else {
			return f.values[prev]
		}
	}
	return math.NaN()
}

// Iterate returns an iterator over all event time in chronological order, without repetition.
func Iterate(functions ...*Function) iter.Seq[time.Time] {
	return func(yield func(time.Time) bool) {
		indexes := make([]int, len(functions))
		// find the reached mins
		times := make([]time.Time, 0, len(functions))
		for {
			times = times[:0] //empty the slice again
			for i, index := range indexes {
				if index < functions[i].Len() {
					on, _ := functions[i].At(index)
					times = append(times, on)
				}
			}
			if len(times) == 0 {
				// All timeseries have been consumed, exit.
				return
			}
			// there are some remaining values:
			var m time.Time
			if len(times) > 0 {
				m = times[0]
				for _, t := range times {
					if t.Before(m) {
						m = t
					}
				}
			}
			// now extract the ones that are equals to the min
			for i, index := range indexes {
				if index >= functions[i].Len() {
					continue
				}
				if on, _ := functions[i].At(index); on.Equal(m) {
					// Updates and consume this value
					indexes[i]++
				}
			}
			if !yield(m) {
				return
			}
		}
	}
}

// Add returns a new function that is the result of adding all functions
func Add(functions ...*Function) *Function {
	s := new(Support)

	for on := range Iterate(functions...) {
		v := 0.0
		for _, f := range functions {
			v += f.F(on)
		}
		s.Append(on, v)
	}
	// compute the mode of the resulting function.
	// currently it is the min of all modes.
	mode := LenMode
	for _, f := range functions {
		mode = min(mode, f.mode)
	}
	return New(s, mode)
}

// Times returns a new function that is the result of multiplying all functions
func Times(functions ...*Function) *Function {
	s := new(Support)
	for on := range Iterate(functions...) {
		v := 1.0
		for _, f := range functions {
			v *= f.F(on)
		}
		s.Append(on, v)
	}
	// compute the mode of the resulting function.
	// currently it is the min of all modes.
	mode := LenMode
	for _, f := range functions {
		mode = min(mode, f.mode)
	}
	return New(s, mode)
}

// Sub returns a new function that is the result of a-b
func Sub(a, b *Function) *Function {
	s := new(Support)
	for on := range Iterate(a, b) {
		s.Append(on, a.F(on)-b.F(on))
	}
	// compute the mode of the resulting function.
	// currently it is the min of all modes.
	mode := min(a.mode, b.mode)
	return New(s, mode)
}

// Div returns a new function that is the result of a-b
func Div(a, b *Function) *Function {
	s := new(Support)
	for on := range Iterate(a, b) {
		s.Append(on, a.F(on)/b.F(on))
	}
	// compute the mode of the resulting function.
	// currently it is the min of all modes.
	mode := min(a.mode, b.mode)
	return New(s, mode)
}

// Sample resample functions on a daily basis.
func Sample(times []time.Time, f *Function) *Function {
	// Prepare the result Timeserie.
	result := new(Support)

	// Loop over the sampling interval computing the delta.
	for _, t := range times {
		result.Append(t, f.F(t))
	}
	return New(result, f.mode)
}
