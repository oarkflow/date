package date

import (
	"testing"
	"time"
)

func TestFormatStrftime(t *testing.T) {
	ti := time.Date(2024, 8, 25, 14, 35, 10, 0, time.UTC)

	tests := []struct {
		layout string
		want   string
	}{
		{"%Y-%m-%d", "2024-08-25"},
		{"%H:%M:%S", "14:35:10"},
		{"%F %T", "2024-08-25 14:35:10"},
		{"%A, %B %d, %Y", "Sunday, August 25, 2024"},
		{"%a, %b %d", "Sun, Aug 25"},
		{"%I:%M %p", "02:35 PM"},
		{"Century: %C", "Century: 20"},
		{"Day of year: %j", "Day of year: 238"},
		{"Weekday: %w", "Weekday: 0"},
		{"ISO Week: %V", "ISO Week: 34"},
	}

	for _, tt := range tests {
		got := Format(ti, tt.layout)
		if got != tt.want {
			t.Errorf("Format(%q) = %q, want %q", tt.layout, got, tt.want)
		}
	}
}

func TestRelativeTime(t *testing.T) {
	ref := time.Date(2024, 5, 10, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		ti   time.Time
		want string
	}{
		{ref.Add(-20 * time.Second), "just now"},
		{ref.Add(-5 * time.Minute), "5 minutes ago"},
		{ref.Add(-3 * time.Hour), "3 hours ago"},
		{ref.Add(-4 * 24 * time.Hour), "4 days ago"},
		{ref.Add(5 * time.Minute), "in 5 minutes"},
		{ref.Add(2 * time.Hour), "in 2 hours"},
		{ref.Add(3 * 24 * time.Hour), "in 3 days"},
	}

	for _, tt := range tests {
		got := RelativeTime(tt.ti, ref)
		if got != tt.want {
			t.Errorf("RelativeTime(%v) = %q, want %q", tt.ti, got, tt.want)
		}
	}
}

func TestOrdinal(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{1, "1st"},
		{2, "2nd"},
		{3, "3rd"},
		{4, "4th"},
		{11, "11th"},
		{12, "12th"},
		{13, "13th"},
		{21, "21st"},
		{22, "22nd"},
		{23, "23rd"},
		{101, "101st"},
	}

	for _, tt := range tests {
		if got := Ordinal(tt.n); got != tt.want {
			t.Errorf("Ordinal(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestNameAndFormattingShortcuts(t *testing.T) {
	if MonthName(time.March) != "March" {
		t.Errorf("expected March")
	}
	if ShortMonthName(time.March) != "Mar" {
		t.Errorf("expected Mar")
	}
	if WeekdayName(time.Wednesday) != "Wednesday" {
		t.Errorf("expected Wednesday")
	}
	if ShortWeekdayName(time.Wednesday) != "Wed" {
		t.Errorf("expected Wed")
	}

	ti := time.Date(2024, 3, 15, 9, 30, 0, 0, time.UTC)
	if FormatDate(ti) != "2024-03-15" {
		t.Errorf("expected 2024-03-15")
	}
	if FormatDateTime(ti) != "2024-03-15 09:30:00" {
		t.Errorf("expected 2024-03-15 09:30:00")
	}
}
