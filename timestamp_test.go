package date

import (
	"testing"
	"time"
)

func TestUnixTimestampConversions(t *testing.T) {
	// 2024-01-01 00:00:00 UTC = 1704067200
	sec := int64(1704067200)
	ti := UnixToTime(sec)
	if ti.Year() != 2024 || ti.Month() != time.January || ti.Day() != 1 {
		t.Errorf("UnixToTime unexpected: %v", ti)
	}
	if ToUnix(ti) != sec {
		t.Errorf("ToUnix = %d, want %d", ToUnix(ti), sec)
	}

	ms := sec * 1000
	tiMs := UnixMilliToTime(ms)
	if ToUnixMilli(tiMs) != ms {
		t.Errorf("ToUnixMilli = %d, want %d", ToUnixMilli(tiMs), ms)
	}

	us := sec * 1000000
	tiUs := UnixMicroToTime(us)
	if ToUnixMicro(tiUs) != us {
		t.Errorf("ToUnixMicro = %d, want %d", ToUnixMicro(tiUs), us)
	}

	ns := sec * 1000000000
	tiNs := UnixNanoToTime(ns)
	if ToUnixNano(tiNs) != ns {
		t.Errorf("ToUnixNano = %d, want %d", ToUnixNano(tiNs), ns)
	}
}

func TestParseTimestampAutoDetect(t *testing.T) {
	sec := int64(1704067200)
	ms := sec * 1000
	us := sec * 1000000
	ns := sec * 1000000000

	tSec := ParseTimestamp(sec)
	tMs := ParseTimestamp(ms)
	tUs := ParseTimestamp(us)
	tNs := ParseTimestamp(ns)

	if !tSec.Equal(tMs) || !tMs.Equal(tUs) || !tUs.Equal(tNs) {
		t.Errorf("all timestamp precisions should parse to equal time: %v, %v, %v, %v", tSec, tMs, tUs, tNs)
	}

	if TimestampPrecisionOf(sec) != "seconds" {
		t.Errorf("expected seconds")
	}
	if TimestampPrecisionOf(ms) != "milliseconds" {
		t.Errorf("expected milliseconds")
	}
	if TimestampPrecisionOf(us) != "microseconds" {
		t.Errorf("expected microseconds")
	}
	if TimestampPrecisionOf(ns) != "nanoseconds" {
		t.Errorf("expected nanoseconds")
	}

	// From string
	tiFromStr, err := ParseTimestampString("1704067200")
	if err != nil || !tiFromStr.Equal(tSec) {
		t.Errorf("ParseTimestampString error or mismatch: %v, %v", tiFromStr, err)
	}
}
