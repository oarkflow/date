package date

import (
	"testing"
	"time"
)

func TestStartAndEndOfDay(t *testing.T) {
	loc := time.UTC
	ti := time.Date(2024, 7, 15, 14, 30, 45, 123456789, loc)

	start := StartOfDay(ti)
	if start.Hour() != 0 || start.Minute() != 0 || start.Second() != 0 || start.Nanosecond() != 0 {
		t.Errorf("StartOfDay expected 00:00:00.0, got %v", start)
	}
	if start.Day() != 15 || start.Month() != 7 || start.Year() != 2024 {
		t.Errorf("StartOfDay date changed: got %v", start)
	}

	end := EndOfDay(ti)
	if end.Hour() != 23 || end.Minute() != 59 || end.Second() != 59 || end.Nanosecond() != 999999999 {
		t.Errorf("EndOfDay expected 23:59:59.999999999, got %v", end)
	}
	if end.Day() != 15 {
		t.Errorf("EndOfDay day changed: got %v", end)
	}
}

func TestStartAndEndOfWeek(t *testing.T) {
	// 2024-07-17 is Wednesday
	wed := time.Date(2024, 7, 17, 15, 0, 0, 0, time.UTC)

	// Default starts on Monday (2024-07-15)
	startMon := StartOfWeek(wed)
	if startMon.Weekday() != time.Monday || startMon.Day() != 15 {
		t.Errorf("expected Monday July 15, got %v", startMon)
	}

	endMon := EndOfWeek(wed)
	if endMon.Weekday() != time.Sunday || endMon.Day() != 21 {
		t.Errorf("expected Sunday July 21, got %v", endMon)
	}

	// Custom: starts on Sunday (2024-07-14)
	startSun := StartOfWeek(wed, time.Sunday)
	if startSun.Weekday() != time.Sunday || startSun.Day() != 14 {
		t.Errorf("expected Sunday July 14, got %v", startSun)
	}

	endSun := EndOfWeek(wed, time.Sunday)
	if endSun.Weekday() != time.Saturday || endSun.Day() != 20 {
		t.Errorf("expected Saturday July 20, got %v", endSun)
	}
}

func TestStartAndEndOfMonth(t *testing.T) {
	ti := time.Date(2024, 2, 14, 10, 0, 0, 0, time.UTC)

	start := StartOfMonth(ti)
	if start.Day() != 1 || start.Month() != time.February {
		t.Errorf("expected Feb 1, got %v", start)
	}

	end := EndOfMonth(ti) // 2024 is leap year, 29 days
	if end.Day() != 29 || end.Month() != time.February {
		t.Errorf("expected Feb 29, got %v", end)
	}

	// Non-leap year
	ti2023 := time.Date(2023, 2, 14, 10, 0, 0, 0, time.UTC)
	end2023 := EndOfMonth(ti2023)
	if end2023.Day() != 28 {
		t.Errorf("expected Feb 28, got %v", end2023)
	}
}

func TestQuarter(t *testing.T) {
	tests := []struct {
		m       time.Month
		quarter int
		qStart  time.Month
		qEnd    time.Month
	}{
		{time.January, 1, time.January, time.March},
		{time.March, 1, time.January, time.March},
		{time.April, 2, time.April, time.June},
		{time.August, 3, time.July, time.September},
		{time.December, 4, time.October, time.December},
	}

	for _, tt := range tests {
		ti := time.Date(2024, tt.m, 15, 0, 0, 0, 0, time.UTC)
		if q := Quarter(ti); q != tt.quarter {
			t.Errorf("Quarter(%v) = %d, want %d", tt.m, q, tt.quarter)
		}
		start := StartOfQuarter(ti)
		if start.Month() != tt.qStart || start.Day() != 1 {
			t.Errorf("StartOfQuarter(%v) = %v, want month %v day 1", tt.m, start, tt.qStart)
		}
		end := EndOfQuarter(ti)
		if end.Month() != tt.qEnd {
			t.Errorf("EndOfQuarter(%v) = %v, want month %v", tt.m, end, tt.qEnd)
		}
	}
}

func TestStartAndEndOfYear(t *testing.T) {
	ti := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	start := StartOfYear(ti)
	if start.Month() != time.January || start.Day() != 1 {
		t.Errorf("StartOfYear expected Jan 1, got %v", start)
	}
	end := EndOfYear(ti)
	if end.Month() != time.December || end.Day() != 31 {
		t.Errorf("EndOfYear expected Dec 31, got %v", end)
	}
}

func TestDateHelpers(t *testing.T) {
	ti := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // Monday
	if !IsWeekday(ti) || IsWeekend(ti) {
		t.Errorf("Monday should be weekday, not weekend")
	}

	sat := time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)
	if !IsWeekend(sat) || IsWeekday(sat) {
		t.Errorf("Saturday should be weekend, not weekday")
	}

	if DaysInMonth(2024, time.February) != 29 {
		t.Errorf("Feb 2024 should have 29 days")
	}
	if DaysInMonth(2023, time.February) != 28 {
		t.Errorf("Feb 2023 should have 28 days")
	}

	if DayOfYear(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)) != 32 {
		t.Errorf("DayOfYear for Feb 1 should be 32")
	}
}

func TestIsSameChecks(t *testing.T) {
	t1 := time.Date(2024, 5, 10, 8, 30, 0, 0, time.UTC)
	t2 := time.Date(2024, 5, 10, 18, 45, 0, 0, time.UTC)
	t3 := time.Date(2024, 5, 11, 8, 30, 0, 0, time.UTC)
	t4 := time.Date(2025, 5, 10, 8, 30, 0, 0, time.UTC)

	if !IsSameDay(t1, t2) {
		t.Errorf("t1 and t2 should be same day")
	}
	if IsSameDay(t1, t3) {
		t.Errorf("t1 and t3 should not be same day")
	}
	if !IsSameMonth(t1, t3) {
		t.Errorf("t1 and t3 should be same month")
	}
	if IsSameMonth(t1, t4) {
		t.Errorf("t1 and t4 should not be same month")
	}
	if !IsSameYear(t1, t3) {
		t.Errorf("t1 and t3 should be same year")
	}
	if IsSameYear(t1, t4) {
		t.Errorf("t1 and t4 should not be same year")
	}
}

func TestBusinessDays(t *testing.T) {
	// 2024-05-17 is Friday
	fri := time.Date(2024, 5, 17, 0, 0, 0, 0, time.UTC)
	// Add 1 business day -> Monday May 20
	mon := AddBusinessDays(fri, 1, nil)
	if mon.Weekday() != time.Monday || mon.Day() != 20 {
		t.Errorf("expected Monday May 20, got %v", mon)
	}

	// With holiday on Monday May 20 -> Tuesday May 21
	holidays := []time.Time{time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC)}
	tue := AddBusinessDays(fri, 1, holidays)
	if tue.Weekday() != time.Tuesday || tue.Day() != 21 {
		t.Errorf("expected Tuesday May 21 with holiday, got %v", tue)
	}

	// BusinessDaysBetween
	count := BusinessDaysBetween(fri, mon, nil)
	if count != 1 {
		t.Errorf("expected 1 business day between fri and mon, got %d", count)
	}

	countWithHol := BusinessDaysBetween(fri, tue, holidays)
	if countWithHol != 1 {
		t.Errorf("expected 1 business day with holiday, got %d", countWithHol)
	}

	// Negative business days
	prevThu := AddBusinessDays(fri, -1, nil)
	if prevThu.Weekday() != time.Thursday || prevThu.Day() != 16 {
		t.Errorf("expected Thursday May 16, got %v", prevThu)
	}
}

func TestDiffDetailed(t *testing.T) {
	t1 := time.Date(2020, 1, 15, 10, 30, 0, 0, time.UTC)
	t2 := time.Date(2023, 3, 20, 14, 45, 30, 0, time.UTC)

	diff := DiffDetailed(t1, t2)
	if diff.Years != 3 || diff.Months != 2 || diff.Days != 5 || diff.Hours != 4 || diff.Minutes != 15 || diff.Seconds != 30 {
		t.Errorf("unexpected DiffDetailed result: %+v", diff)
	}
}

func TestClamp(t *testing.T) {
	minT := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	maxT := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	before := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	after := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	inside := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	if Clamp(before, minT, maxT) != minT {
		t.Errorf("expected minT for before")
	}
	if Clamp(after, minT, maxT) != maxT {
		t.Errorf("expected maxT for after")
	}
	if Clamp(inside, minT, maxT) != inside {
		t.Errorf("expected inside for inside")
	}
}
