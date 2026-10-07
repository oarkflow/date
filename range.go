package date

import (
	"errors"
	"time"
)

// DateRange represents an inclusive range between two points in time.
type DateRange struct {
	Start time.Time
	End   time.Time
}

// NewDateRange creates a DateRange. It returns an error if end is before start.
func NewDateRange(start, end time.Time) (DateRange, error) {
	if end.Before(start) {
		return DateRange{}, errors.New("date range end must not be before start")
	}
	return DateRange{Start: start, End: end}, nil
}

// NewDateRangeFromDuration creates a DateRange starting at start with the given duration.
func NewDateRangeFromDuration(start time.Time, d time.Duration) DateRange {
	return DateRange{Start: start, End: start.Add(d)}
}

// Duration returns the duration of the range (End − Start).
func (r DateRange) Duration() time.Duration {
	return r.End.Sub(r.Start)
}

// Days returns the number of whole calendar days spanned by the range.
// A range within a single day returns 1.
func (r DateRange) Days() int {
	startDay := StartOfDay(r.Start)
	endDay := StartOfDay(r.End)
	diff := endDay.Sub(startDay)
	n := int(diff.Hours()/24) + 1
	if n < 1 {
		n = 1
	}
	return n
}

// IsEmpty reports whether the range has zero duration.
func (r DateRange) IsEmpty() bool {
	return !r.End.After(r.Start)
}

// Contains reports whether t falls within [Start, End] (inclusive).
func (r DateRange) Contains(t time.Time) bool {
	return !t.Before(r.Start) && !t.After(r.End)
}

// Overlaps reports whether r and other share any point in time.
func (r DateRange) Overlaps(other DateRange) bool {
	return r.Start.Before(other.End) && other.Start.Before(r.End)
}

// Intersect returns the intersection of r and other, plus a boolean indicating
// whether the intersection is non-empty.
func (r DateRange) Intersect(other DateRange) (DateRange, bool) {
	start := r.Start
	if other.Start.After(start) {
		start = other.Start
	}
	end := r.End
	if other.End.Before(end) {
		end = other.End
	}
	if end.Before(start) {
		return DateRange{}, false
	}
	return DateRange{Start: start, End: end}, true
}

// Union returns the smallest DateRange that contains both r and other.
func (r DateRange) Union(other DateRange) DateRange {
	start := r.Start
	if other.Start.Before(start) {
		start = other.Start
	}
	end := r.End
	if other.End.After(end) {
		end = other.End
	}
	return DateRange{Start: start, End: end}
}

// Each calls fn once for each calendar day in the range (inclusive), passing
// the start-of-day time. Iteration stops early if fn returns false.
func (r DateRange) Each(fn func(time.Time) bool) {
	cur := StartOfDay(r.Start)
	end := StartOfDay(r.End)
	for !cur.After(end) {
		if !fn(cur) {
			return
		}
		cur = cur.AddDate(0, 0, 1)
	}
}

// EachMonth calls fn once for the first day of each calendar month that overlaps
// the range. Iteration stops early if fn returns false.
func (r DateRange) EachMonth(fn func(time.Time) bool) {
	cur := StartOfMonth(r.Start)
	for !cur.After(r.End) {
		if !fn(cur) {
			return
		}
		cur = cur.AddDate(0, 1, 0)
	}
}

// Split divides the range into sub-ranges each of duration d.
// The last sub-range may be shorter than d.
func (r DateRange) Split(d time.Duration) []DateRange {
	if d <= 0 || r.IsEmpty() {
		return nil
	}
	var parts []DateRange
	cur := r.Start
	for cur.Before(r.End) {
		next := cur.Add(d)
		if next.After(r.End) {
			next = r.End
		}
		parts = append(parts, DateRange{Start: cur, End: next})
		cur = next
	}
	return parts
}

// BusinessDays returns a slice of all business days (Mon–Fri, excluding holidays)
// within the range, one entry per day at start-of-day.
func (r DateRange) BusinessDays(holidays []time.Time) []time.Time {
	var days []time.Time
	r.Each(func(t time.Time) bool {
		if IsWeekday(t) && !isHolidayTime(t, holidays) {
			days = append(days, t)
		}
		return true
	})
	return days
}

// Extend returns a new DateRange with End advanced by d.
func (r DateRange) Extend(d time.Duration) DateRange {
	return DateRange{Start: r.Start, End: r.End.Add(d)}
}

// Shift returns a new DateRange shifted forward (positive d) or backward (negative d).
func (r DateRange) Shift(d time.Duration) DateRange {
	return DateRange{Start: r.Start.Add(d), End: r.End.Add(d)}
}
