package date

import (
	"fmt"
	"time"
)

// UnixToTime converts a Unix timestamp (seconds since epoch) to time.Time (UTC).
func UnixToTime(unix int64) time.Time {
	return time.Unix(unix, 0).UTC()
}

// UnixMilliToTime converts a Unix timestamp in milliseconds to time.Time (UTC).
func UnixMilliToTime(ms int64) time.Time {
	return time.UnixMilli(ms).UTC()
}

// UnixMicroToTime converts a Unix timestamp in microseconds to time.Time (UTC).
func UnixMicroToTime(us int64) time.Time {
	return time.UnixMicro(us).UTC()
}

// UnixNanoToTime converts a Unix timestamp in nanoseconds to time.Time (UTC).
func UnixNanoToTime(ns int64) time.Time {
	return time.Unix(0, ns).UTC()
}

// ToUnix returns the Unix time (seconds since epoch) of t.
func ToUnix(t time.Time) int64 {
	return t.Unix()
}

// ToUnixMilli returns the Unix time in milliseconds of t.
func ToUnixMilli(t time.Time) int64 {
	return t.UnixMilli()
}

// ToUnixMicro returns the Unix time in microseconds of t.
func ToUnixMicro(t time.Time) int64 {
	return t.UnixMicro()
}

// ToUnixNano returns the Unix time in nanoseconds of t.
func ToUnixNano(t time.Time) int64 {
	return t.UnixNano()
}

// timestampPrecision enumerates the detected precision of a numeric timestamp.
type timestampPrecision int

const (
	precisionSeconds timestampPrecision = iota
	precisionMillis
	precisionMicros
	precisionNanos
)

// detectPrecision heuristically determines the precision of a Unix timestamp
// by its magnitude.
//
// Heuristic boundaries (as of the year 2000–2100 window):
//
//	seconds  : 9–10 digits  (946684800 – 4102444800)
//	millis   : 12–13 digits
//	micros   : 15–16 digits
//	nanos    : 18–19 digits
func detectPrecision(v int64) timestampPrecision {
	abs := v
	if abs < 0 {
		abs = -abs
	}
	switch {
	case abs < 1e12:
		return precisionSeconds
	case abs < 1e15:
		return precisionMillis
	case abs < 1e18:
		return precisionMicros
	default:
		return precisionNanos
	}
}

// ParseTimestamp auto-detects whether v is a Unix timestamp in seconds,
// milliseconds, microseconds, or nanoseconds and returns the corresponding
// time.Time (UTC).
func ParseTimestamp(v int64) time.Time {
	switch detectPrecision(v) {
	case precisionMillis:
		return UnixMilliToTime(v)
	case precisionMicros:
		return UnixMicroToTime(v)
	case precisionNanos:
		return UnixNanoToTime(v)
	default:
		return UnixToTime(v)
	}
}

// ParseTimestampString parses a numeric string as a Unix timestamp, applying
// the same auto-detection as ParseTimestamp. Returns an error if the string is
// not a valid integer.
func ParseTimestampString(s string) (time.Time, error) {
	var v int64
	_, err := fmt.Sscanf(s, "%d", &v)
	if err != nil {
		return time.Time{}, fmt.Errorf("timestamp: cannot parse %q as integer: %w", s, err)
	}
	return ParseTimestamp(v), nil
}

// TimestampPrecisionOf returns a string describing the auto-detected precision
// of a numeric timestamp value ("seconds", "milliseconds", "microseconds",
// "nanoseconds").
func TimestampPrecisionOf(v int64) string {
	switch detectPrecision(v) {
	case precisionMillis:
		return "milliseconds"
	case precisionMicros:
		return "microseconds"
	case precisionNanos:
		return "nanoseconds"
	default:
		return "seconds"
	}
}
