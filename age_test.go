package date

import (
	"testing"
	"time"
)

func TestAgeCalculation(t *testing.T) {
	birth := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	refBeforeBday := time.Date(2020, 5, 14, 0, 0, 0, 0, time.UTC)
	refAfterBday := time.Date(2020, 5, 16, 0, 0, 0, 0, time.UTC)

	if age := Calculate(birth, refBeforeBday); age != 29 {
		t.Errorf("age before bday = %d, want 29", age)
	}

	if age := Calculate(birth, refAfterBday); age != 30 {
		t.Errorf("age after bday = %d, want 30", age)
	}

	// Leap year test: someone born Feb 29 2000 is age 0 on Feb 27 2001, and turns 1 on Feb 28 2001
	leapBirth := time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC)
	leapRefBefore := time.Date(2001, 2, 27, 0, 0, 0, 0, time.UTC)
	leapRefAfter := time.Date(2001, 2, 28, 0, 0, 0, 0, time.UTC)
	if age := Calculate(leapBirth, leapRefBefore); age != 0 {
		t.Errorf("age on Feb 27 following leap year = %d, want 0", age)
	}
	if age := Calculate(leapBirth, leapRefAfter); age != 1 {
		t.Errorf("age on Feb 28 following leap year = %d, want 1", age)
	}
}

func TestAgeDetailed(t *testing.T) {
	birth := time.Date(2000, 1, 15, 0, 0, 0, 0, time.UTC)
	target := time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC)

	y, m, d := AgeDetailed(birth, target)
	if y != 24 || m != 1 || d != 24 {
		t.Errorf("AgeDetailed = %d years, %d months, %d days, want 24y, 1m, 24d", y, m, d)
	}

	// Swapped order
	y2, m2, d2 := AgeDetailed(target, birth)
	if y2 != y || m2 != m || d2 != d {
		t.Errorf("AgeDetailed should handle swapped dates identically")
	}
}

func TestLeapYearHelpers(t *testing.T) {
	if !IsLeapYear(2020) || !IsLeapYear(2024) || !IsLeapYear(2000) {
		t.Errorf("expected 2000, 2020, 2024 to be leap years")
	}
	if IsLeapYear(1900) || IsLeapYear(2023) {
		t.Errorf("expected 1900 and 2023 to not be leap years")
	}

	next, ok := NextLeapYear(2021)
	if !ok || next != 2024 {
		t.Errorf("NextLeapYear(2021) = %d, want 2024", next)
	}

	prev, ok := PrevLeapYear(2023)
	if !ok || prev != 2020 {
		t.Errorf("PrevLeapYear(2023) = %d, want 2020", prev)
	}
}
