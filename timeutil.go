package timeserie

import (
	"math"
	"time"
)

// Here goes functions used to play with time, and period of times related to timeserie.

const Day = 24 * time.Hour

// DayDate returns a comparable Time to identity a single day.
func DayDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// TimeCond is a function to filter in some events
type TimeCond func(day time.Time) bool

func CondWeekday(day time.Weekday) TimeCond {
	return func(d time.Time) bool { return d.Weekday() == day }
}

func CondMonthday(day int) TimeCond {
	return func(t time.Time) bool { _, _, d := t.Date(); return d == day }
}

// CondEndOfMonth return true if 't'  represent the last day of the month.
func CondEndOfMonth(t time.Time) bool { return t.Month() != t.Add(Day).Month() }

// CondQuarterly returns true if 't' represent the first day of a new quarter, jan or apr or jul or oct 1st.
func CondQuarterly(t time.Time) bool { _, m, d := t.Date(); return d == 1 && m%3 == 1 }

// Condyearly returns true if 't' represent the first day of a new year (a january first).
func CondYearly(t time.Time) bool { _, m, d := t.Date(); return d == 1 && m == 1 }

// Days returns a list of all days starting with 'from' (included) ends after 'end' and return only
// days accepted by time condition.
func Days(from, end time.Time, accept TimeCond) []time.Time {
	var result []time.Time
	y, m, day := from.Date()
	for d := DayDate(y, m, day); d.Before(end); d, day = DayDate(y, m, day+1), day+1 {
		if accept(d) {
			result = append(result, d)
		}
	}
	return result
}

// Scanner is a function that can be used in the Scan method.
type Scanner func(c, s float64) float64

// Acc is a Scanner to compute the cumulative value of a support.
var (
	// list all scanner function
	ScannerAcc = func(c, v float64) float64 { return c + v }
)

// ScannerIf returns a Scanner based on a condition. If condition is true the
// value is return, otherwise NaN is returned.
func ScannerIf(cond ValueCond) Scanner {
	return func(c, v float64) float64 {
		if cond(v) {
			return v
		} else {
			return math.NaN()
		}
	}
}

// ValueCond is a function that can be used to filter a support.
type ValueCond func(v float64) bool

var (
	CondPositive = func(v float64) bool { return v > 0 }
	CondNegative = func(v float64) bool { return v < 0 }
)

// If computes a new Support by keep only the value for a given condition.
func (s *Support) If(cond ValueCond) *Support { return s.Scan(0, ScannerIf(cond)) }
