package date

import (
	"testing"
	"time"
)

func TestTimezoneValidation(t *testing.T) {
	if !IsValidTimezone("UTC") {
		t.Errorf("expected UTC to be valid")
	}
	if !IsValidTimezone("Asia/Kathmandu") {
		t.Errorf("expected Asia/Kathmandu to be valid")
	}
	if !IsValidTimezone("America/New_York") {
		t.Errorf("expected America/New_York to be valid")
	}
	if IsValidTimezone("NonExistent/Place") {
		t.Errorf("expected NonExistent/Place to be invalid")
	}
}

func TestConvertTZ(t *testing.T) {
	utcTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	ktmTime, err := ConvertTZ(utcTime, "Asia/Kathmandu")
	if err != nil {
		t.Fatalf("unexpected error converting to Kathmandu: %v", err)
	}

	// Kathmandu is UTC+5:45
	if ktmTime.Hour() != 17 || ktmTime.Minute() != 45 {
		t.Errorf("expected 17:45 in Kathmandu, got %02d:%02d", ktmTime.Hour(), ktmTime.Minute())
	}
}

func TestTimezoneOffset(t *testing.T) {
	utcTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	offset, err := TimezoneOffset("Asia/Kathmandu", utcTime)
	if err != nil {
		t.Fatalf("unexpected error getting offset: %v", err)
	}

	// 5 hours 45 mins = 20700 seconds
	expected := 5*3600 + 45*60
	if offset != expected {
		t.Errorf("expected offset %d, got %d", expected, offset)
	}
}

func TestGuessTimezone(t *testing.T) {
	tests := []struct {
		abbr string
		iana string
	}{
		{"NPT", "Asia/Kathmandu"},
		{"IST", "Asia/Kolkata"},
		{"PST", "America/Los_Angeles"},
		{"UTC", "UTC"},
	}

	for _, tt := range tests {
		iana, ok := GuessTimezone(tt.abbr)
		if !ok || iana != tt.iana {
			t.Errorf("GuessTimezone(%q) = %q, want %q", tt.abbr, iana, tt.iana)
		}
	}

	_, ok := GuessTimezone("UNKNOWN_XYZ")
	if ok {
		t.Errorf("expected false for unknown abbr")
	}
}

func TestLocalAndUTCConversions(t *testing.T) {
	loc := time.Date(2024, 5, 1, 12, 0, 0, 0, time.Local)
	utc := LocalToUTC(loc)
	if utc.Location() != time.UTC {
		t.Errorf("expected UTC location, got %v", utc.Location())
	}

	backToLoc := UTCToLocal(utc)
	if backToLoc.Location() != time.Local {
		t.Errorf("expected Local location, got %v", backToLoc.Location())
	}
}

func TestParseInTimezone(t *testing.T) {
	ti, err := ParseInTimezone("2024-05-01 12:00:00", "Asia/Kathmandu")
	if err != nil {
		t.Fatalf("unexpected error parsing: %v", err)
	}

	if ti.Location().String() != "Asia/Kathmandu" {
		t.Errorf("expected Asia/Kathmandu location, got %v", ti.Location())
	}
}
