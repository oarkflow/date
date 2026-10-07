package date

import (
	"testing"
	"time"
)

func TestNthWeekdayOfMonth(t *testing.T) {
	// 4th Thursday of November 2024 (US Thanksgiving) is Nov 28, 2024
	thanksgiving := NthWeekdayOfMonth(2024, time.November, time.Thursday, 4)
	if thanksgiving.Year() != 2024 || thanksgiving.Month() != time.November || thanksgiving.Day() != 28 {
		t.Errorf("expected 2024-11-28 for 4th Thursday of Nov, got %v", thanksgiving)
	}

	// 1st Monday of September 2024 (US Labor Day) is Sep 2, 2024
	laborDay := NthWeekdayOfMonth(2024, time.September, time.Monday, 1)
	if laborDay.Day() != 2 {
		t.Errorf("expected Sep 2 for 1st Monday of Sep, got %v", laborDay)
	}

	// Non-existent 5th occurrence
	invalid := NthWeekdayOfMonth(2024, time.February, time.Monday, 5)
	if !invalid.IsZero() {
		t.Errorf("expected zero time for non-existent 5th Monday of Feb 2024, got %v", invalid)
	}
}

func TestLastWeekdayOfMonth(t *testing.T) {
	// Last Friday of May 2024 is May 31
	lastFri := LastWeekdayOfMonth(2024, time.May, time.Friday)
	if lastFri.Day() != 31 || lastFri.Weekday() != time.Friday {
		t.Errorf("expected Friday May 31, got %v", lastFri)
	}

	// Last Sunday of Feb 2024 (leap year: ends on Thursday Feb 29) is Feb 25
	lastSun := LastWeekdayOfMonth(2024, time.February, time.Sunday)
	if lastSun.Day() != 25 || lastSun.Weekday() != time.Sunday {
		t.Errorf("expected Sunday Feb 25, got %v", lastSun)
	}
}

func TestEasterDate(t *testing.T) {
	// Easter 2024 was March 31
	e2024 := EasterDate(2024)
	if e2024.Month() != time.March || e2024.Day() != 31 {
		t.Errorf("Easter 2024 should be March 31, got %v", e2024)
	}

	// Easter 2025 will be April 20
	e2025 := EasterDate(2025)
	if e2025.Month() != time.April || e2025.Day() != 20 {
		t.Errorf("Easter 2025 should be April 20, got %v", e2025)
	}
}

func TestCalendarWeeks(t *testing.T) {
	weeks := CalendarWeeks(2024, time.July)
	if len(weeks) < 5 {
		t.Errorf("July 2024 should have at least 5 weeks in grid, got %d", len(weeks))
	}
	// Check first day in month is present in first week
	foundFirst := false
	for _, day := range weeks[0] {
		if !day.IsZero() && day.Day() == 1 {
			foundFirst = true
			break
		}
	}
	if !foundFirst {
		t.Errorf("expected Day 1 to be in first week")
	}
}

func TestDaysBetweenHelpers(t *testing.T) {
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	t3 := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)

	if weeks := WeeksBetween(t1, t2); weeks != 2 {
		t.Errorf("expected 2 weeks between Jan 1 and Jan 15, got %d", weeks)
	}

	if months := MonthsBetween(t1, t3); months != 6 {
		t.Errorf("expected 6 months between Jan 1 and Jul 1, got %d", months)
	}
}
