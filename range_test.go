package date

import (
	"testing"
	"time"
)

func TestDateRangeBasics(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

	r, err := NewDateRange(start, end)
	if err != nil {
		t.Fatalf("unexpected error creating DateRange: %v", err)
	}

	if r.Duration() != 9*24*time.Hour {
		t.Errorf("expected 9 days duration, got %v", r.Duration())
	}

	if r.Days() != 10 {
		t.Errorf("expected 10 days spanned, got %d", r.Days())
	}

	if r.IsEmpty() {
		t.Errorf("range should not be empty")
	}

	// Invalid range
	_, err = NewDateRange(end, start)
	if err == nil {
		t.Errorf("expected error when end is before start")
	}
}

func TestDateRangeContainsAndOverlaps(t *testing.T) {
	r1, _ := NewDateRange(
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	)

	inside := time.Date(2024, 1, 5, 12, 0, 0, 0, time.UTC)
	outside := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	if !r1.Contains(inside) {
		t.Errorf("expected r1 to contain inside")
	}
	if r1.Contains(outside) {
		t.Errorf("expected r1 to not contain outside")
	}

	r2, _ := NewDateRange(
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	)
	r3, _ := NewDateRange(
		time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	)

	if !r1.Overlaps(r2) {
		t.Errorf("r1 and r2 should overlap")
	}
	if r1.Overlaps(r3) {
		t.Errorf("r1 and r3 should not overlap")
	}
}

func TestDateRangeIntersectAndUnion(t *testing.T) {
	r1, _ := NewDateRange(
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	)
	r2, _ := NewDateRange(
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	)

	inter, ok := r1.Intersect(r2)
	if !ok {
		t.Fatalf("expected intersection to exist")
	}
	if !inter.Start.Equal(r2.Start) || !inter.End.Equal(r1.End) {
		t.Errorf("unexpected intersection: %v to %v", inter.Start, inter.End)
	}

	union := r1.Union(r2)
	if !union.Start.Equal(r1.Start) || !union.End.Equal(r2.End) {
		t.Errorf("unexpected union: %v to %v", union.Start, union.End)
	}
}

func TestDateRangeIteration(t *testing.T) {
	r, _ := NewDateRange(
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	)

	days := 0
	r.Each(func(ti time.Time) bool {
		days++
		return true
	})
	if days != 5 {
		t.Errorf("expected 5 days iterated, got %d", days)
	}

	// Early stop
	days = 0
	r.Each(func(ti time.Time) bool {
		days++
		return days < 3
	})
	if days != 3 {
		t.Errorf("expected 3 days before break, got %d", days)
	}
}

func TestDateRangeSplitAndShift(t *testing.T) {
	r, _ := NewDateRange(
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	)

	chunks := r.Split(24 * time.Hour)
	if len(chunks) != 4 {
		t.Errorf("expected 4 chunks of 24h, got %d", len(chunks))
	}

	shifted := r.Shift(24 * time.Hour)
	if shifted.Start.Day() != 2 || shifted.End.Day() != 6 {
		t.Errorf("unexpected shifted range: %v to %v", shifted.Start, shifted.End)
	}
}
