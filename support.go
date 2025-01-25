package timeserie

import (
	"iter"
	"math"
	"slices"
	"sort"
	"time"
)

// Support struct contains the time based finite [support] for real-valued functions.
//
// [support]: https://en.wikipedia.org/wiki/Support_(mathematics)
type Support struct {
	times  []time.Time
	values []float64
}

// make the support sortable by times.

type sortable struct{ *Support }

func (s sortable) Less(i, j int) bool { return s.times[i].Before(s.times[j]) }

func (s sortable) Swap(i, j int) {
	s.times[i], s.times[j] = s.times[j], s.times[i]
	s.values[i], s.values[j] = s.values[j], s.values[i]
}

// sort timeserie in chronollogical order.
func (s *Support) sort() { sort.Sort(sortable{s}) }

// Len returns the timeserie support's length.
func (s Support) Len() int { return len(s.times) }

// At return the point at given position in the support.
func (s Support) At(i int) (time.Time, float64) { return s.times[i], s.values[i] }

// Append a point to this support.
func (s *Support) Append(on time.Time, q float64) {
	if math.IsNaN(q) {
		return
	}
	s.times, s.values = append(s.times, on), append(s.values, q)
	s.sort()

}

// Find returns the index of the closest value before 't'.
func (s Support) Find(t time.Time) int {
	next := slices.IndexFunc(s.times, func(i time.Time) bool { return i.After(t) })
	if next < 0 {
		return s.Len() - 1
	}
	return next - 1
}

// Values return an iterator over all values in the support.
func (s *Support) Values() iter.Seq2[time.Time, float64] {
	return func(yield func(time.Time, float64) bool) {
		for i, on := range s.times {
			if !yield(on, s.values[i]) {
				return
			}
		}
	}
}

// Iterate over dates in the support.
func (s *Support) Times() iter.Seq[time.Time] {
	return func(yield func(time.Time) bool) {
		for _, on := range s.times {
			if !yield(on) {
				return
			}
		}
	}
}

// Delta loop over all interval in this support
// and returns a new support defined at the end of each interval
// with the delta on this interval.
func (s *Support) Delta() *Support {
	result := new(Support)
	for i := 1; i < len(s.times); i++ {
		result.Append(s.times[i], s.values[i]-s.values[i-1])
	}
	return result
}

// Scan computes a new Support by cumulating values, so that
//
//	c:= initial
//
// for value 'v' in support 's' do c <- f(c, v) and append c to the new support
func (s *Support) Scan(initial float64, scanner Scanner) *Support {
	res := new(Support)
	c := initial
	for t, v := range s.Values() {
		c = scanner(c, v)
		res.Append(t, c)
	}
	return res
}
