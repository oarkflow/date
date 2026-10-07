package date

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"1h30m", 90 * time.Minute},
		{"2d 3h 15m", 2*24*time.Hour + 3*time.Hour + 15*time.Minute},
		{"1 week 2 days", 9 * 24 * time.Hour},
		{"90 seconds", 90 * time.Second},
		{"500ms", 500 * time.Millisecond},
		{"2 hours", 2 * time.Hour},
		{"1.5h", 90 * time.Minute},
	}

	for _, tt := range tests {
		d, err := ParseDuration(tt.input)
		if err != nil {
			t.Errorf("ParseDuration(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if d != tt.expected {
			t.Errorf("ParseDuration(%q) = %v, want %v", tt.input, d, tt.expected)
		}
	}

	// Error cases
	invalid := []string{"", "abc", "10xyz", "not a duration"}
	for _, in := range invalid {
		_, err := ParseDuration(in)
		if err == nil {
			t.Errorf("ParseDuration(%q) expected error, got nil", in)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	d := 2*24*time.Hour + 3*time.Hour + 15*time.Minute + 5*time.Second

	short := FormatDuration(d, true)
	if short != "2d 3h 15m 5s" {
		t.Errorf("FormatDuration short = %q, want %q", short, "2d 3h 15m 5s")
	}

	long := FormatDuration(d, false)
	if long != "2 days 3 hours 15 minutes 5 seconds" {
		t.Errorf("FormatDuration long = %q, want %q", long, "2 days 3 hours 15 minutes 5 seconds")
	}

	// Zero
	if FormatDuration(0, true) != "0s" {
		t.Errorf("expected 0s, got %q", FormatDuration(0, true))
	}
	if FormatDuration(0, false) != "0 seconds" {
		t.Errorf("expected 0 seconds, got %q", FormatDuration(0, false))
	}

	// Negative
	neg := -1 * time.Hour
	if FormatDuration(neg, true) != "-1h" {
		t.Errorf("expected -1h, got %q", FormatDuration(neg, true))
	}
}

func TestHumanizeDuration(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{20 * time.Second, "just now"},
		{60 * time.Second, "about a minute"},
		{15 * time.Minute, "about 15 minutes"},
		{75 * time.Minute, "about an hour"},
		{5 * time.Hour, "about 5 hours"},
		{30 * time.Hour, "a day"},
		{5 * 24 * time.Hour, "5 days"},
		{45 * 24 * time.Hour, "about a month"},
		{120 * 24 * time.Hour, "about 4 months"},
		{400 * 24 * time.Hour, "about a year"},
		{800 * 24 * time.Hour, "about 2 years"},
	}

	for _, tt := range tests {
		got := HumanizeDuration(tt.d)
		if got != tt.want {
			t.Errorf("HumanizeDuration(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}

func TestRoundAndTruncateDuration(t *testing.T) {
	d := 3*time.Minute + 35*time.Second
	rounded := RoundDuration(d, time.Minute)
	if rounded != 4*time.Minute {
		t.Errorf("RoundDuration = %v, want 4m", rounded)
	}

	truncated := TruncateDuration(d, time.Minute)
	if truncated != 3*time.Minute {
		t.Errorf("TruncateDuration = %v, want 3m", truncated)
	}
}
