package date

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// durUnit maps lowercase unit tokens to their time.Duration equivalent.
// "week" and "month" are approximations.
var durUnitMap = map[string]time.Duration{
	// nanoseconds
	"ns": time.Nanosecond, "nanosecond": time.Nanosecond, "nanoseconds": time.Nanosecond,
	// microseconds
	"us": time.Microsecond, "µs": time.Microsecond,
	"microsecond": time.Microsecond, "microseconds": time.Microsecond,
	// milliseconds
	"ms": time.Millisecond, "millisecond": time.Millisecond, "milliseconds": time.Millisecond,
	// seconds
	"s": time.Second, "sec": time.Second, "secs": time.Second,
	"second": time.Second, "seconds": time.Second,
	// minutes
	"m": time.Minute, "min": time.Minute, "mins": time.Minute,
	"minute": time.Minute, "minutes": time.Minute,
	// hours
	"h": time.Hour, "hr": time.Hour, "hrs": time.Hour,
	"hour": time.Hour, "hours": time.Hour,
	// days
	"d": 24 * time.Hour, "day": 24 * time.Hour, "days": 24 * time.Hour,
	// weeks
	"w": 7 * 24 * time.Hour, "wk": 7 * 24 * time.Hour, "wks": 7 * 24 * time.Hour,
	"week": 7 * 24 * time.Hour, "weeks": 7 * 24 * time.Hour,
	// months (approximate: 30 days)
	"mo": 30 * 24 * time.Hour, "month": 30 * 24 * time.Hour, "months": 30 * 24 * time.Hour,
	// years (approximate: 365 days)
	"y": 365 * 24 * time.Hour, "yr": 365 * 24 * time.Hour, "yrs": 365 * 24 * time.Hour,
	"year": 365 * 24 * time.Hour, "years": 365 * 24 * time.Hour,
}

// ParseDuration parses a human-readable duration string and returns the
// equivalent time.Duration. It supports Go's native format (e.g. "1h30m")
// as a fallback, plus extended human forms such as:
//
//	"2d 3h 15m"
//	"1 week 2 days"
//	"90 seconds"
//	"2y 3mo 1w 4d 5h 30m 15s"
func ParseDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty duration string")
	}

	// Try Go's native format first.
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	var total time.Duration
	remaining := strings.ToLower(s)

	for remaining != "" {
		remaining = strings.TrimLeftFunc(remaining, unicode.IsSpace)
		if remaining == "" {
			break
		}

		// Read the numeric part (may include a decimal point).
		numEnd := 0
		for numEnd < len(remaining) && (unicode.IsDigit(rune(remaining[numEnd])) || remaining[numEnd] == '.') {
			numEnd++
		}
		if numEnd == 0 {
			return 0, fmt.Errorf("invalid duration %q: expected number", s)
		}
		numStr := remaining[:numEnd]
		remaining = remaining[numEnd:]

		val, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid duration %q: %w", s, err)
		}

		// Skip optional whitespace between number and unit.
		remaining = strings.TrimLeftFunc(remaining, unicode.IsSpace)

		// Read the unit part.
		unitEnd := 0
		for unitEnd < len(remaining) && !unicode.IsDigit(rune(remaining[unitEnd])) && !unicode.IsSpace(rune(remaining[unitEnd])) {
			unitEnd++
		}
		if unitEnd == 0 {
			return 0, fmt.Errorf("invalid duration %q: missing unit after %q", s, numStr)
		}
		unit := remaining[:unitEnd]
		remaining = remaining[unitEnd:]

		mult, ok := durUnitMap[unit]
		if !ok {
			return 0, fmt.Errorf("invalid duration %q: unknown unit %q", s, unit)
		}

		total += time.Duration(math.Round(val * float64(mult)))
	}

	return total, nil
}

// FormatDuration formats d as a human-readable string.
// If short is true, abbreviated units are used (e.g. "2d 3h 15m 5s").
// If short is false, full unit names are used (e.g. "2 days 3 hours 15 minutes 5 seconds").
// Only non-zero components are included.
func FormatDuration(d time.Duration, short bool) string {
	if d == 0 {
		if short {
			return "0s"
		}
		return "0 seconds"
	}

	negative := d < 0
	if negative {
		d = -d
	}

	type unit struct {
		size      time.Duration
		short     string
		singular  string
		plural    string
	}

	units := []unit{
		{365 * 24 * time.Hour, "y", "year", "years"},
		{30 * 24 * time.Hour, "mo", "month", "months"},
		{7 * 24 * time.Hour, "w", "week", "weeks"},
		{24 * time.Hour, "d", "day", "days"},
		{time.Hour, "h", "hour", "hours"},
		{time.Minute, "m", "minute", "minutes"},
		{time.Second, "s", "second", "seconds"},
		{time.Millisecond, "ms", "millisecond", "milliseconds"},
	}

	var parts []string
	for _, u := range units {
		n := int64(d / u.size)
		if n == 0 {
			continue
		}
		d -= time.Duration(n) * u.size
		if short {
			parts = append(parts, fmt.Sprintf("%d%s", n, u.short))
		} else {
			label := u.plural
			if n == 1 {
				label = u.singular
			}
			parts = append(parts, fmt.Sprintf("%d %s", n, label))
		}
	}

	if len(parts) == 0 {
		if short {
			return "0s"
		}
		return "0 seconds"
	}

	result := strings.Join(parts, " ")
	if negative {
		result = "-" + result
	}
	return result
}

// HumanizeDuration returns a coarse, human-friendly description of d such as
// "just now", "about 3 minutes", "2 hours", "3 days", "about 2 months", etc.
func HumanizeDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}

	seconds := d.Seconds()

	switch {
	case seconds < 45:
		return "just now"
	case seconds < 90:
		return "about a minute"
	case seconds < 45*60:
		return fmt.Sprintf("about %d minutes", int(math.Round(seconds/60)))
	case seconds < 90*60:
		return "about an hour"
	case seconds < 24*3600:
		return fmt.Sprintf("about %d hours", int(math.Round(seconds/3600)))
	case seconds < 48*3600:
		return "a day"
	case seconds < 30*24*3600:
		return fmt.Sprintf("%d days", int(math.Round(seconds/(24*3600))))
	case seconds < 60*24*3600:
		return "about a month"
	case seconds < 365*24*3600:
		return fmt.Sprintf("about %d months", int(math.Round(seconds/(30*24*3600))))
	case seconds < 2*365*24*3600:
		return "about a year"
	default:
		return fmt.Sprintf("about %d years", int(math.Round(seconds/(365*24*3600))))
	}
}

// RoundDuration rounds d to the nearest multiple of unit.
func RoundDuration(d, unit time.Duration) time.Duration {
	if unit <= 0 {
		return d
	}
	return d.Round(unit)
}

// TruncateDuration truncates d to a multiple of unit (towards zero).
func TruncateDuration(d, unit time.Duration) time.Duration {
	if unit <= 0 {
		return d
	}
	return d.Truncate(unit)
}
