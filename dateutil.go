package date

import (
	"time"
)

// DateDiff holds the component-wise difference between two time.Time values.
type DateDiff struct {
	Years   int
	Months  int
	Weeks   int
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

// StartOfDay returns t truncated to the beginning of its day (00:00:00.000000000).
func StartOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the last instant of t's day (23:59:59.999999999).
func EndOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, 999999999, t.Location())
}

// StartOfWeek returns the start of the week that contains t.
// startDay specifies which weekday begins the week (default: time.Monday).
func StartOfWeek(t time.Time, startDay ...time.Weekday) time.Time {
	firstDay := time.Monday
	if len(startDay) > 0 {
		firstDay = startDay[0]
	}
	wd := t.Weekday()
	diff := int(wd) - int(firstDay)
	if diff < 0 {
		diff += 7
	}
	start := StartOfDay(t.AddDate(0, 0, -diff))
	return start
}

// EndOfWeek returns the last instant of the week that contains t.
// startDay specifies which weekday begins the week (default: time.Monday).
func EndOfWeek(t time.Time, startDay ...time.Weekday) time.Time {
	start := StartOfWeek(t, startDay...)
	return EndOfDay(start.AddDate(0, 0, 6))
}

// StartOfMonth returns the first moment of t's month.
func StartOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the last instant of t's month.
func EndOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	// Day 0 of the next month = last day of current month
	last := time.Date(y, m+1, 0, 0, 0, 0, 0, t.Location())
	return EndOfDay(last)
}

// StartOfQuarter returns the first moment of the quarter in which t falls.
func StartOfQuarter(t time.Time) time.Time {
	month := t.Month()
	// Quarter starts: Jan, Apr, Jul, Oct
	quarterStartMonth := month - ((month - 1) % 3)
	return time.Date(t.Year(), quarterStartMonth, 1, 0, 0, 0, 0, t.Location())
}

// EndOfQuarter returns the last instant of the quarter in which t falls.
func EndOfQuarter(t time.Time) time.Time {
	start := StartOfQuarter(t)
	// 3 months later, day 0 = last day of previous month
	endMonth := start.Month() + 3
	end := time.Date(start.Year(), endMonth, 0, 0, 0, 0, 0, t.Location())
	return EndOfDay(end)
}

// StartOfYear returns the first moment of t's year.
func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear returns the last instant of t's year.
func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.December, 31, 23, 59, 59, 999999999, t.Location())
}

// Quarter returns the quarter number (1–4) for t.
func Quarter(t time.Time) int {
	return (int(t.Month())-1)/3 + 1
}

// WeekNumber returns the ISO 8601 week number for t.
func WeekNumber(t time.Time) int {
	_, week := t.ISOWeek()
	return week
}

// DayOfYear returns the day of the year for t (1–365 or 1–366 for leap years).
func DayOfYear(t time.Time) int {
	return t.YearDay()
}

// DaysInMonth returns the number of days in the given month and year.
func DaysInMonth(year int, month time.Month) int {
	// Day 0 of the next month = last day of the given month
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// IsWeekend reports whether t falls on a Saturday or Sunday.
func IsWeekend(t time.Time) bool {
	wd := t.Weekday()
	return wd == time.Saturday || wd == time.Sunday
}

// IsWeekday reports whether t falls on a Monday–Friday.
func IsWeekday(t time.Time) bool {
	return !IsWeekend(t)
}

// IsSameDay reports whether a and b fall on the same calendar day
// (ignoring clock time and timezone differences).
func IsSameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

// IsSameWeek reports whether a and b fall in the same ISO week of the same year.
func IsSameWeek(a, b time.Time) bool {
	ay, aw := a.ISOWeek()
	by, bw := b.ISOWeek()
	return ay == by && aw == bw
}

// IsSameMonth reports whether a and b fall in the same month of the same year.
func IsSameMonth(a, b time.Time) bool {
	ay, am, _ := a.Date()
	by, bm, _ := b.Date()
	return ay == by && am == bm
}

// IsSameYear reports whether a and b fall in the same year.
func IsSameYear(a, b time.Time) bool {
	return a.Year() == b.Year()
}

// isHolidayTime returns true if t (day-level) is present in the holidays slice.
func isHolidayTime(t time.Time, holidays []time.Time) bool {
	ty, tm, td := t.Date()
	for _, h := range holidays {
		hy, hm, hd := h.Date()
		if ty == hy && tm == hm && td == hd {
			return true
		}
	}
	return false
}

// AddBusinessDays adds n business days to t, skipping weekends and any dates
// in the optional holidays slice. n may be negative to go backwards.
func AddBusinessDays(t time.Time, n int, holidays []time.Time) time.Time {
	if n == 0 {
		return t
	}
	step := 1
	if n < 0 {
		step = -1
		n = -n
	}
	cur := t
	for n > 0 {
		cur = cur.AddDate(0, 0, step)
		if IsWeekday(cur) && !isHolidayTime(cur, holidays) {
			n--
		}
	}
	return cur
}

// BusinessDaysBetween counts the number of business days between a and b
// (exclusive of a, inclusive of b direction), skipping weekends and holidays.
// Returns a negative number if b is before a.
func BusinessDaysBetween(a, b time.Time, holidays []time.Time) int {
	if a.Equal(b) {
		return 0
	}
	sign := 1
	if b.Before(a) {
		a, b = b, a
		sign = -1
	}
	count := 0
	cur := StartOfDay(a).AddDate(0, 0, 1)
	end := StartOfDay(b)
	for !cur.After(end) {
		if IsWeekday(cur) && !isHolidayTime(cur, holidays) {
			count++
		}
		cur = cur.AddDate(0, 0, 1)
	}
	return count * sign
}

// DiffDetailed returns the component-wise difference between from and to.
// All fields are non-negative; the caller can determine direction by comparing
// the original times.
func DiffDetailed(from, to time.Time) DateDiff {
	if to.Before(from) {
		from, to = to, from
	}

	y1, m1, d1 := from.Date()
	y2, m2, d2 := to.Date()
	h1, min1, s1 := from.Clock()
	h2, min2, s2 := to.Clock()

	years := y2 - y1
	months := int(m2) - int(m1)
	days := d2 - d1
	hours := h2 - h1
	minutes := min2 - min1
	seconds := s2 - s1

	if seconds < 0 {
		seconds += 60
		minutes--
	}
	if minutes < 0 {
		minutes += 60
		hours--
	}
	if hours < 0 {
		hours += 24
		days--
	}
	if days < 0 {
		// borrow from previous month
		prevMonth := time.Date(y2, m2, 0, 0, 0, 0, 0, to.Location())
		days += prevMonth.Day()
		months--
	}
	if months < 0 {
		months += 12
		years--
	}

	weeks := days / 7
	days = days % 7

	return DateDiff{
		Years:   years,
		Months:  months,
		Weeks:   weeks,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
	}
}

// Clamp returns t if it is within [min, max]; otherwise it returns min or max.
func Clamp(t, minT, maxT time.Time) time.Time {
	if t.Before(minT) {
		return minT
	}
	if t.After(maxT) {
		return maxT
	}
	return t
}
