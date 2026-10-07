package date

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// strftimeDirectives maps strftime-style % directives to Go format tokens or
// custom handler functions.
var strftimeMap = map[byte]string{
	'A': "Monday",
	'a': "Mon",
	'B': "January",
	'b': "Jan",
	'C': "", // handled separately: century
	'd': "02",
	'e': "_2",
	'F': "2006-01-02",
	'G': "",  // ISO year — handled separately
	'H': "15",
	'I': "03",
	'j': "",  // day of year — handled separately
	'k': " 3", // hour (0-23) space-padded
	'l': " 3", // hour (1-12) space-padded — Go doesn't distinguish, same token
	'm': "01",
	'M': "04",
	'n': "\n",
	'p': "PM",
	'P': "pm",
	'R': "15:04",
	'r': "03:04:05 PM",
	'S': "05",
	'T': "15:04:05",
	't': "\t",
	'u': "", // weekday as number (1=Mon) — handled separately
	'V': "", // ISO week number — handled separately
	'w': "", // weekday as number (0=Sun) — handled separately
	'X': "15:04:05",
	'x': "01/02/06",
	'Y': "2006",
	'y': "06",
	'Z': "MST",
	'z': "-0700",
	'%': "%",
}

// Format formats t using a strftime-style layout string.
// Most POSIX strftime directives are supported.
func Format(t time.Time, layout string) string {
	var b strings.Builder
	i := 0
	for i < len(layout) {
		if layout[i] != '%' {
			b.WriteByte(layout[i])
			i++
			continue
		}
		i++
		if i >= len(layout) {
			b.WriteByte('%')
			break
		}
		dir := layout[i]
		i++

		switch dir {
		case 'C':
			fmt.Fprintf(&b, "%02d", t.Year()/100)
		case 'G':
			year, _ := t.ISOWeek()
			fmt.Fprintf(&b, "%04d", year)
		case 'j':
			fmt.Fprintf(&b, "%03d", t.YearDay())
		case 'u':
			wd := int(t.Weekday())
			if wd == 0 {
				wd = 7
			}
			fmt.Fprintf(&b, "%d", wd)
		case 'V':
			_, week := t.ISOWeek()
			fmt.Fprintf(&b, "%02d", week)
		case 'w':
			fmt.Fprintf(&b, "%d", int(t.Weekday()))
		default:
			if goFmt, ok := strftimeMap[dir]; ok {
				b.WriteString(t.Format(goFmt))
			} else {
				b.WriteByte('%')
				b.WriteByte(dir)
			}
		}
	}
	return b.String()
}

// RelativeTime returns a human-readable string describing t relative to ref,
// such as "3 days ago", "in 2 hours", "just now", "yesterday", "tomorrow".
func RelativeTime(t, ref time.Time) string {
	diff := t.Sub(ref)
	abs := diff
	if abs < 0 {
		abs = -abs
	}

	past := diff < 0

	wrap := func(s string) string {
		if past {
			return s + " ago"
		}
		return "in " + s
	}

	seconds := math.Abs(diff.Seconds())

	switch {
	case seconds < 45:
		return "just now"
	case seconds < 90:
		return wrap("a minute")
	case seconds < 45*60:
		return wrap(fmt.Sprintf("%d minutes", int(math.Round(seconds/60))))
	case seconds < 90*60:
		return wrap("an hour")
	case seconds < 36*3600:
		hours := int(math.Round(seconds / 3600))
		if hours == 1 {
			return wrap("an hour")
		}
		return wrap(fmt.Sprintf("%d hours", hours))
	case IsSameDay(t, ref.AddDate(0, 0, 1)) && past:
		return "tomorrow" // ref is in the past relative to t
	case IsSameDay(t, ref.AddDate(0, 0, -1)) && !past:
		return "yesterday"
	case IsSameDay(t, ref):
		return "today"
	case seconds < 48*3600:
		if past {
			return "yesterday"
		}
		return "tomorrow"
	case seconds < 30*24*3600:
		days := int(math.Round(seconds / (24 * 3600)))
		return wrap(fmt.Sprintf("%d days", days))
	case seconds < 60*24*3600:
		return wrap("a month")
	case seconds < 365*24*3600:
		months := int(math.Round(seconds / (30 * 24 * 3600)))
		return wrap(fmt.Sprintf("%d months", months))
	case seconds < 2*365*24*3600:
		return wrap("a year")
	default:
		years := int(math.Round(seconds / (365 * 24 * 3600)))
		return wrap(fmt.Sprintf("%d years", years))
	}
}

// Ordinal returns the English ordinal suffix string for n (e.g. 1→"1st",
// 2→"2nd", 11→"11th", 23→"23rd").
func Ordinal(n int) string {
	abs := n
	if abs < 0 {
		abs = -abs
	}
	suffix := "th"
	switch abs % 100 {
	case 11, 12, 13:
		// teens always use "th"
	default:
		switch abs % 10 {
		case 1:
			suffix = "st"
		case 2:
			suffix = "nd"
		case 3:
			suffix = "rd"
		}
	}
	return fmt.Sprintf("%d%s", n, suffix)
}

// MonthName returns the full English name of the given month.
func MonthName(m time.Month) string {
	return m.String()
}

// ShortMonthName returns the 3-letter English abbreviation of the given month.
func ShortMonthName(m time.Month) string {
	return m.String()[:3]
}

// WeekdayName returns the full English name of the given weekday.
func WeekdayName(w time.Weekday) string {
	return w.String()
}

// ShortWeekdayName returns the 3-letter English abbreviation of the given weekday.
func ShortWeekdayName(w time.Weekday) string {
	return w.String()[:3]
}

// FormatDate formats t as "YYYY-MM-DD".
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateTime formats t as "YYYY-MM-DD HH:MM:SS".
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatDateTimeMs formats t as "YYYY-MM-DD HH:MM:SS.mmm".
func FormatDateTimeMs(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000")
}

// FormatRFC3339 formats t as RFC 3339 (a profile of ISO 8601).
func FormatRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FormatRFC3339Nano formats t as RFC 3339 with nanosecond precision.
func FormatRFC3339Nano(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

// FormatHuman returns a friendly date string like "Wednesday, January 15, 2025".
func FormatHuman(t time.Time) string {
	return t.Format("Monday, January 2, 2006")
}

// FormatHumanTime returns a friendly datetime string like "Wednesday, January 15, 2025 at 3:04 PM".
func FormatHumanTime(t time.Time) string {
	return t.Format("Monday, January 2, 2006 at 3:04 PM")
}
