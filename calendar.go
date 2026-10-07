package date

import (
	"time"
)

// NthWeekdayOfMonth returns the nth occurrence of weekday in the given year and
// month. For example, NthWeekdayOfMonth(2024, time.November, time.Thursday, 4)
// returns the 4th Thursday of November 2024 (US Thanksgiving).
// n must be in [1, 5]. If the month does not have an nth occurrence, the zero
// Time is returned.
func NthWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, n int) time.Time {
	if n < 1 || n > 5 {
		return time.Time{}
	}
	// First day of month.
	first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	// Days until the target weekday from the first of the month.
	offset := int(weekday) - int(first.Weekday())
	if offset < 0 {
		offset += 7
	}
	day := 1 + offset + (n-1)*7
	result := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	if result.Month() != month {
		return time.Time{} // n-th occurrence doesn't exist in this month
	}
	return result
}

// LastWeekdayOfMonth returns the last occurrence of weekday in the given year
// and month.
func LastWeekdayOfMonth(year int, month time.Month, weekday time.Weekday) time.Time {
	// Start from last day of month and walk backwards.
	last := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC) // last day
	offset := int(last.Weekday()) - int(weekday)
	if offset < 0 {
		offset += 7
	}
	return last.AddDate(0, 0, -offset)
}

// EasterDate returns the date of Easter Sunday for the given Gregorian year
// using the Anonymous Gregorian algorithm.
func EasterDate(year int) time.Time {
	a := year % 19
	b := year / 100
	c := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	month := (h + l - 7*m + 114) / 31
	day := ((h + l - 7*m + 114) % 31) + 1
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// IsHoliday reports whether t (at day granularity) matches any date in the
// holidays slice.
func IsHoliday(t time.Time, holidays []time.Time) bool {
	return isHolidayTime(t, holidays)
}

// CalendarWeeks returns the weeks of the given year and month as a 2-D slice.
// Each inner slice is a 7-element array [Sunday … Saturday]; days that fall
// outside the month are represented as the zero time.Time.
func CalendarWeeks(year int, month time.Month) [][]time.Time {
	first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := DaysInMonth(year, month)

	// Find the Sunday on or before the first of the month.
	startOffset := int(first.Weekday()) // Sunday=0
	calStart := first.AddDate(0, 0, -startOffset)

	// Build weeks until we've covered all days of the month.
	var weeks [][]time.Time
	cur := calStart
	for {
		week := make([]time.Time, 7)
		for i := 0; i < 7; i++ {
			if cur.Month() == month && cur.Day() >= 1 && cur.Day() <= lastDay {
				week[i] = cur
			}
			cur = cur.AddDate(0, 0, 1)
		}
		weeks = append(weeks, week)
		if cur.After(time.Date(year, month, lastDay, 0, 0, 0, 0, time.UTC)) {
			break
		}
	}
	return weeks
}

// DaysUntil returns the number of whole calendar days from now until t.
// Returns a negative number if t is in the past.
func DaysUntil(t time.Time) int {
	now := StartOfDay(time.Now())
	target := StartOfDay(t)
	diff := target.Sub(now)
	return int(diff.Hours() / 24)
}

// DaysSince returns the number of whole calendar days that have elapsed since t.
// Returns a negative number if t is in the future.
func DaysSince(t time.Time) int {
	return -DaysUntil(t)
}

// WeeksBetween returns the number of complete weeks between a and b.
func WeeksBetween(a, b time.Time) int {
	if b.Before(a) {
		a, b = b, a
	}
	return int(b.Sub(a).Hours() / (24 * 7))
}

// MonthsBetween returns the approximate number of whole calendar months between
// a and b (based on year/month arithmetic, ignoring day-of-month).
func MonthsBetween(a, b time.Time) int {
	if b.Before(a) {
		a, b = b, a
	}
	years := b.Year() - a.Year()
	months := int(b.Month()) - int(a.Month())
	total := years*12 + months
	if b.Day() < a.Day() {
		total--
	}
	if total < 0 {
		total = 0
	}
	return total
}
