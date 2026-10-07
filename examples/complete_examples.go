package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/oarkflow/date"
)

func printHeader(title string) {
	banner := strings.Repeat("=", 70)
	fmt.Printf("\n%s\n  %s\n%s\n\n", banner, title, banner)
}

func printSubHeader(subtitle string) {
	fmt.Printf("\n--- %s ---\n", subtitle)
}

// RunAllDemos runs every feature demonstration sequentially.
func RunAllDemos() {
	DemoDateArithmetic()
	DemoDateRanges()
	DemoTimezoneServices()
	DemoDurationUtilities()
	DemoCalendarServices()
	DemoFormattingAndLocalization()
	DemoTimestampUtilities()
	DemoAgeAndLeapYears()
	DemoNepaliCalendar()
	DemoNaturalDateParsing()
}

// 1. Date Arithmetic, Boundaries & Business Days
func DemoDateArithmetic() {
	printHeader("1. DATE ARITHMETIC, BOUNDARIES & BUSINESS DAYS")

	now := time.Date(2024, time.July, 17, 15, 30, 45, 500, time.UTC) // Wednesday
	fmt.Printf("Reference Date: %v (Wednesday)\n\n", now.Format(time.RFC3339))

	printSubHeader("Period Boundaries")
	fmt.Printf("StartOfDay:        %v\n", date.StartOfDay(now))
	fmt.Printf("EndOfDay:          %v\n", date.EndOfDay(now))
	fmt.Printf("StartOfWeek (Mon): %v\n", date.StartOfWeek(now))
	fmt.Printf("EndOfWeek (Mon):   %v\n", date.EndOfWeek(now))
	fmt.Printf("StartOfWeek (Sun): %v\n", date.StartOfWeek(now, time.Sunday))
	fmt.Printf("EndOfWeek (Sun):   %v\n", date.EndOfWeek(now, time.Sunday))
	fmt.Printf("StartOfMonth:      %v\n", date.StartOfMonth(now))
	fmt.Printf("EndOfMonth:        %v\n", date.EndOfMonth(now))
	fmt.Printf("StartOfQuarter:    %v\n", date.StartOfQuarter(now))
	fmt.Printf("EndOfQuarter:      %v\n", date.EndOfQuarter(now))
	fmt.Printf("StartOfYear:       %v\n", date.StartOfYear(now))
	fmt.Printf("EndOfYear:         %v\n", date.EndOfYear(now))

	printSubHeader("Period Properties & Queries")
	fmt.Printf("Quarter:           Q%d\n", date.Quarter(now))
	fmt.Printf("ISO Week Number:   Week %d\n", date.WeekNumber(now))
	fmt.Printf("Day of Year:       Day %d\n", date.DayOfYear(now))
	fmt.Printf("Days in Month:     %d days (July 2024)\n", date.DaysInMonth(2024, time.July))
	fmt.Printf("Days in Feb 2024:  %d days (leap year)\n", date.DaysInMonth(2024, time.February))
	fmt.Printf("Days in Feb 2023:  %d days (common year)\n", date.DaysInMonth(2023, time.February))
	fmt.Printf("IsWeekday:         %t\n", date.IsWeekday(now))
	fmt.Printf("IsWeekend:         %t\n", date.IsWeekend(now))

	printSubHeader("Business Days & Holidays")
	friday := time.Date(2024, 5, 17, 9, 0, 0, 0, time.UTC) // Friday
	holidays := []time.Time{
		time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC), // Monday is holiday
	}
	fmt.Printf("Starting from:               %s (Friday)\n", friday.Format("2006-01-02"))
	fmt.Printf("+1 Business Day (no hol):    %s\n", date.AddBusinessDays(friday, 1, nil).Format("2006-01-02 (Monday)"))
	fmt.Printf("+1 Business Day (with hol):  %s\n", date.AddBusinessDays(friday, 1, holidays).Format("2006-01-02 (Tuesday, skipping Monday holiday)"))
	fmt.Printf("+5 Business Days:            %s\n", date.AddBusinessDays(friday, 5, holidays).Format("2006-01-02"))
	target := time.Date(2024, 5, 28, 0, 0, 0, 0, time.UTC)
	fmt.Printf("BusinessDaysBetween:         %d days between May 17 and May 28\n", date.BusinessDaysBetween(friday, target, holidays))

	printSubHeader("Detailed Component Difference (DiffDetailed)")
	startEvent := time.Date(2020, 1, 15, 9, 30, 0, 0, time.UTC)
	endEvent := time.Date(2024, 7, 20, 14, 45, 30, 0, time.UTC)
	diff := date.DiffDetailed(startEvent, endEvent)
	fmt.Printf("Between %s and %s:\n", startEvent.Format("2006-01-02 15:04:05"), endEvent.Format("2006-01-02 15:04:05"))
	fmt.Printf("  -> %d years, %d months, %d weeks, %d days, %d hours, %d minutes, %d seconds\n",
		diff.Years, diff.Months, diff.Weeks, diff.Days, diff.Hours, diff.Minutes, diff.Seconds)
}

// 2. Date Ranges
func DemoDateRanges() {
	printHeader("2. DATE RANGE STRUCT & OPERATIONS")

	r1Start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	r1End := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	r1, _ := date.NewDateRange(r1Start, r1End)

	r2Start := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	r2End := time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC)
	r2, _ := date.NewDateRange(r2Start, r2End)

	fmt.Printf("Range 1: %s to %s (Days: %d, Duration: %v)\n",
		r1.Start.Format("2006-01-02"), r1.End.Format("2006-01-02"), r1.Days(), r1.Duration())
	fmt.Printf("Range 2: %s to %s (Days: %d, Duration: %v)\n\n",
		r2.Start.Format("2006-01-02"), r2.End.Format("2006-01-02"), r2.Days(), r2.Duration())

	printSubHeader("Range Set Operations")
	midPoint := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)
	fmt.Printf("r1.Contains(2024-01-05): %t\n", r1.Contains(midPoint))
	fmt.Printf("r1.Overlaps(r2):         %t\n", r1.Overlaps(r2))

	if inter, ok := r1.Intersect(r2); ok {
		fmt.Printf("r1.Intersect(r2):        %s to %s\n",
			inter.Start.Format("2006-01-02"), inter.End.Format("2006-01-02"))
	}
	union := r1.Union(r2)
	fmt.Printf("r1.Union(r2):            %s to %s\n",
		union.Start.Format("2006-01-02"), union.End.Format("2006-01-02"))

	printSubHeader("Range Splitting & Shifting")
	smallRange, _ := date.NewDateRange(
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	)
	chunks := smallRange.Split(24 * time.Hour)
	fmt.Printf("Split into 24h chunks (%d chunks):\n", len(chunks))
	for i, c := range chunks {
		fmt.Printf("  Chunk %d: %s -> %s\n", i+1, c.Start.Format("2006-01-02 15:04"), c.End.Format("2006-01-02 15:04"))
	}
	shifted := smallRange.Shift(7 * 24 * time.Hour)
	fmt.Printf("Shifted +1 week:         %s to %s\n",
		shifted.Start.Format("2006-01-02"), shifted.End.Format("2006-01-02"))

	printSubHeader("Business Days Within Range")
	bdays := r1.BusinessDays(nil)
	fmt.Printf("Working days in Range 1 (%d days):\n", len(bdays))
	var dayNames []string
	for _, d := range bdays {
		dayNames = append(dayNames, d.Format("Jan 02 (Mon)"))
	}
	fmt.Printf("  %s\n", strings.Join(dayNames, ", "))
}

// 3. Timezone Utilities & Conversion
func DemoTimezoneServices() {
	printHeader("3. TIMEZONE UTILITIES & CONVERSION")

	fmt.Printf("IsValidTimezone(\"Asia/Kathmandu\"): %t\n", date.IsValidTimezone("Asia/Kathmandu"))
	fmt.Printf("IsValidTimezone(\"Fake/Zone\"):      %t\n\n", date.IsValidTimezone("Fake/Zone"))

	utc := time.Date(2024, 7, 15, 12, 0, 0, 0, time.UTC)
	fmt.Printf("Base Time: %s UTC\n", utc.Format("2006-01-02 15:04:05"))

	targetZones := []string{
		"Asia/Kathmandu",
		"Asia/Kolkata",
		"Asia/Tokyo",
		"Europe/London",
		"Europe/Paris",
		"America/New_York",
		"America/Los_Angeles",
		"Australia/Sydney",
	}

	fmt.Println("\nWorld Clocks via ConvertTZ & TimezoneOffset:")
	for _, tz := range targetZones {
		converted, err := date.ConvertTZ(utc, tz)
		if err != nil {
			continue
		}
		offset, _ := date.TimezoneOffset(tz, utc)
		hours := float64(offset) / 3600.0
		fmt.Printf("  %-22s -> %s (UTC%+0.2f)\n", tz, converted.Format("2006-01-02 15:04:05 MST"), hours)
	}

	printSubHeader("Abbreviation Mapping (GuessTimezone)")
	abbrs := []string{"NPT", "IST", "PST", "EST", "CET", "JST", "BST", "AEDT"}
	for _, a := range abbrs {
		if iana, ok := date.GuessTimezone(a); ok {
			fmt.Printf("  %-5s => %s\n", a, iana)
		}
	}

	printSubHeader("Direct Parse in Target Timezone (ParseInTimezone)")
	raw := "2024-04-15 10:30:00"
	parsed, err := date.ParseInTimezone(raw, "Asia/Kathmandu")
	if err == nil {
		fmt.Printf("Parsed %q in Asia/Kathmandu: %v\n", raw, parsed)
	}
}

// 4. Duration Parsing, Formatting & Humanization
func DemoDurationUtilities() {
	printHeader("4. DURATION PARSING, FORMATTING & HUMANIZING")

	inputs := []string{
		"2d 3h 15m",
		"1 week 2 days",
		"90 seconds",
		"1.5h",
		"500ms",
		"3y 2mo 5d",
	}

	fmt.Println("Human Duration Parsing:")
	for _, in := range inputs {
		d, err := date.ParseDuration(in)
		if err != nil {
			fmt.Printf("  %-15s => error: %v\n", in, err)
			continue
		}
		fmt.Printf("  %-15s => %-22v (Short: %s, Long: %s)\n",
			in, d, date.FormatDuration(d, true), date.FormatDuration(d, false))
	}

	printSubHeader("HumanizeDuration (Fuzzy descriptions)")
	durations := []time.Duration{
		25 * time.Second,
		75 * time.Second,
		18 * time.Minute,
		75 * time.Minute,
		5 * time.Hour,
		36 * time.Hour,
		12 * 24 * time.Hour,
		45 * 24 * time.Hour,
		180 * 24 * time.Hour,
		500 * 24 * time.Hour,
	}

	for _, d := range durations {
		fmt.Printf("  %-20v => %s\n", d, date.HumanizeDuration(d))
	}

	printSubHeader("Rounding & Truncating Durations")
	rawDur := 3*time.Hour + 37*time.Minute + 42*time.Second
	fmt.Printf("Raw Duration:       %v\n", rawDur)
	fmt.Printf("Round to Hour:      %v\n", date.RoundDuration(rawDur, time.Hour))
	fmt.Printf("Truncate to Hour:   %v\n", date.TruncateDuration(rawDur, time.Hour))
	fmt.Printf("Round to 15m:       %v\n", date.RoundDuration(rawDur, 15*time.Minute))
}

// 5. Calendar Services & Event Calculations
func DemoCalendarServices() {
	printHeader("5. CALENDAR UTILITIES & EVENT CALCULATIONS")

	printSubHeader("Nth & Last Weekday of Month")
	thanksgiving := date.NthWeekdayOfMonth(2024, time.November, time.Thursday, 4)
	fmt.Printf("US Thanksgiving 2024 (4th Thursday Nov): %s\n", thanksgiving.Format("Monday, Jan 02, 2006"))

	mothersDay := date.NthWeekdayOfMonth(2024, time.May, time.Sunday, 2)
	fmt.Printf("Mother's Day 2024 (2nd Sunday May):     %s\n", mothersDay.Format("Monday, Jan 02, 2006"))

	lastFri := date.LastWeekdayOfMonth(2024, time.May, time.Friday)
	fmt.Printf("Last Friday of May 2024:                 %s\n", lastFri.Format("Monday, Jan 02, 2006"))

	easter2024 := date.EasterDate(2024)
	easter2025 := date.EasterDate(2025)
	fmt.Printf("Easter Sunday 2024:                      %s\n", easter2024.Format("Jan 02, 2006"))
	fmt.Printf("Easter Sunday 2025:                      %s\n", easter2025.Format("Jan 02, 2006"))

	printSubHeader("Month Calendar Grid (CalendarWeeks for July 2024)")
	weeks := date.CalendarWeeks(2024, time.July)
	fmt.Println(" Su  Mo  Tu  We  Th  Fr  Sa")
	for _, w := range weeks {
		var row []string
		for _, day := range w {
			if day.IsZero() {
				row = append(row, "   ")
			} else {
				row = append(row, fmt.Sprintf("%3d", day.Day()))
			}
		}
		fmt.Println(strings.Join(row, " "))
	}

	printSubHeader("Interval Helpers")
	d1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC)
	fmt.Printf("WeeksBetween(Jan 1, Oct 15):   %d complete weeks\n", date.WeeksBetween(d1, d2))
	fmt.Printf("MonthsBetween(Jan 1, Oct 15):  %d complete months\n", date.MonthsBetween(d1, d2))
}

// 6. Formatting & Localization
func DemoFormattingAndLocalization() {
	printHeader("6. FORMATTING & LOCALIZATION (STRFTIME & SHORTCUTS)")

	t := time.Date(2024, time.August, 25, 14, 35, 10, 0, time.UTC)

	printSubHeader("Strftime-Style Formatting")
	layouts := []string{
		"%Y-%m-%d %H:%M:%S",
		"%A, %B %d, %Y",
		"%a, %b %d (%r)",
		"Century: %C | Day of year: %j | ISO Week: %V | Weekday: %w",
		"%F %T",
	}

	for _, l := range layouts {
		fmt.Printf("  %-55s => %s\n", l, date.Format(t, l))
	}

	printSubHeader("Relative Time (RelativeTime)")
	ref := time.Now()
	testTimes := []struct {
		name string
		val  time.Time
	}{
		{"30s ago", ref.Add(-30 * time.Second)},
		{"5m ago", ref.Add(-5 * time.Minute)},
		{"2h ago", ref.Add(-2 * time.Hour)},
		{"Yesterday", ref.AddDate(0, 0, -1)},
		{"5 days ago", ref.AddDate(0, 0, -5)},
		{"3 months ago", ref.AddDate(0, -3, 0)},
		{"In 10m", ref.Add(10 * time.Minute)},
		{"Tomorrow", ref.AddDate(0, 0, 1)},
		{"In 2 weeks", ref.AddDate(0, 0, 14)},
	}

	for _, item := range testTimes {
		fmt.Printf("  %-15s => %s\n", item.name, date.RelativeTime(item.val, ref))
	}

	printSubHeader("English Ordinals")
	nums := []int{1, 2, 3, 4, 11, 12, 13, 21, 22, 23, 31, 100, 102}
	var ords []string
	for _, n := range nums {
		ords = append(ords, date.Ordinal(n))
	}
	fmt.Println(" ", strings.Join(ords, ", "))

	printSubHeader("Built-in Shortcut Formatters")
	fmt.Printf("FormatDate:        %s\n", date.FormatDate(t))
	fmt.Printf("FormatDateTime:    %s\n", date.FormatDateTime(t))
	fmt.Printf("FormatDateTimeMs:  %s\n", date.FormatDateTimeMs(t))
	fmt.Printf("FormatRFC3339:     %s\n", date.FormatRFC3339(t))
	fmt.Printf("FormatHuman:       %s\n", date.FormatHuman(t))
	fmt.Printf("FormatHumanTime:   %s\n", date.FormatHumanTime(t))
	fmt.Printf("Names:             Month=%s (%s), Weekday=%s (%s)\n",
		date.MonthName(t.Month()), date.ShortMonthName(t.Month()),
		date.WeekdayName(t.Weekday()), date.ShortWeekdayName(t.Weekday()))
}

// 7. Unix Timestamp Utilities & Auto-Detection
func DemoTimestampUtilities() {
	printHeader("7. UNIX TIMESTAMP UTILITIES & AUTO-DETECTION")

	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	s := date.ToUnix(base)
	ms := date.ToUnixMilli(base)
	us := date.ToUnixMicro(base)
	ns := date.ToUnixNano(base)

	fmt.Printf("Base Time: %s\n", base.Format(time.RFC3339))
	fmt.Printf("  Seconds:      %19d  (detected as %s)\n", s, date.TimestampPrecisionOf(s))
	fmt.Printf("  Milliseconds: %19d  (detected as %s)\n", ms, date.TimestampPrecisionOf(ms))
	fmt.Printf("  Microseconds: %19d  (detected as %s)\n", us, date.TimestampPrecisionOf(us))
	fmt.Printf("  Nanoseconds:  %19d  (detected as %s)\n", ns, date.TimestampPrecisionOf(ns))

	printSubHeader("Auto-Detection Parsing (ParseTimestamp & ParseTimestampString)")
	parsedS := date.ParseTimestamp(s)
	parsedMs := date.ParseTimestamp(ms)
	parsedUs := date.ParseTimestamp(us)
	parsedNs := date.ParseTimestamp(ns)
	fmt.Printf("ParseTimestamp(seconds):      %v\n", parsedS)
	fmt.Printf("ParseTimestamp(milliseconds): %v\n", parsedMs)
	fmt.Printf("ParseTimestamp(microseconds): %v\n", parsedUs)
	fmt.Printf("ParseTimestamp(nanoseconds):  %v\n", parsedNs)

	fromStr, _ := date.ParseTimestampString("1704067200000")
	fmt.Printf("ParseTimestampString(\"1704067200000\"): %v\n", fromStr)
}

// 8. Age & Leap Year Calculations
func DemoAgeAndLeapYears() {
	printHeader("8. AGE & LEAP YEAR CALCULATIONS")

	birth := time.Date(1995, 6, 15, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Birthdate: %s\n", birth.Format("2006-01-02"))
	fmt.Printf("Age (years to now):   %d years old\n", date.CalculateToNow(birth))

	years, months, days := date.AgeDetailedToNow(birth)
	fmt.Printf("AgeDetailedToNow:     %d years, %d months, and %d days old\n", years, months, days)

	printSubHeader("Leap Year Analysis")
	yearsToCheck := []int{1900, 2000, 2020, 2023, 2024, 2025}
	for _, y := range yearsToCheck {
		fmt.Printf("  Year %d: LeapYear = %t\n", y, date.IsLeapYear(y))
	}
	prev, _ := date.PrevLeapYear(2023)
	next, _ := date.NextLeapYear(2023)
	fmt.Printf("Leap year before 2023: %d | Leap year after 2023: %d\n", prev, next)
}

// 9. Nepali Calendar (Bikram Sambat ↔ Gregorian)
func DemoNepaliCalendar() {
	printHeader("9. NEPALI (BIKRAM SAMBAT) CALENDAR CONVERSIONS")

	// Current Nepali Time
	npNow := date.Now()
	fmt.Printf("Current Nepali Time: %s\n", npNow.String())
	fmt.Printf("  Year: %d, Month: %d, Day: %d\n", npNow.Year(), npNow.Month(), npNow.Day())
	fmt.Printf("  Corresponding English Time: %s\n", npNow.GetEnglishTime().Format("2006-01-02 15:04:05"))

	printSubHeader("Gregorian (AD) -> Nepali (BS)")
	enY, enM, enD := 2024, 4, 13 // New year boundary in Nepal
	npDate, err := date.EnglishToNepali(enY, enM, enD)
	if err == nil {
		fmt.Printf("  AD %04d-%02d-%02d => BS %04d-%02d-%02d (Nepali Baishakh 1)\n",
			enY, enM, enD, npDate[0], npDate[1], npDate[2])
	}

	printSubHeader("Nepali (BS) -> Gregorian (AD)")
	bsY, bsM, bsD := 2081, 1, 1
	enRes, err := date.NepaliToEnglish(bsY, bsM, bsD)
	if err == nil {
		fmt.Printf("  BS %04d-%02d-%02d => AD %04d-%02d-%02d\n",
			bsY, bsM, bsD, enRes[0], enRes[1], enRes[2])
	}

	printSubHeader("Parsing Nepali Date String (ParseNP)")
	npDateStr := "2080-09-15 14:30:00"
	parsedNP, err := date.ParseNP(npDateStr, "%Y-%m-%d %H:%M:%S")
	if err == nil {
		fmt.Printf("  Parsed %q => %s (AD: %s)\n",
			npDateStr, parsedNP.String(), parsedNP.GetEnglishTime().Format("2006-01-02 15:04:05"))
	}
}

// 10. Natural Date Expressions
func DemoNaturalDateParsing() {
	printHeader("10. NATURAL DATE EXPRESSIONS")

	ref := time.Date(2024, 7, 17, 12, 0, 0, 0, time.UTC) // Wednesday
	fmt.Printf("Reference Date: %s (Wednesday)\n\n", ref.Format("2006-01-02 15:04:05"))

	expressions := []string{
		"yesterday",
		"tomorrow",
		"next friday",
		"last monday",
		"first day of this month",
		"last day of this month",
	}

	for _, expr := range expressions {
		t, exprType, err := date.ParseNaturalDate(expr, ref)
		if err != nil {
			fmt.Printf("  %-25s => error: %v\n", expr, err)
			continue
		}
		fmt.Printf("  %-25s => %-22s (ExprType: %d)\n", expr, t.Format("2006-01-02 15:04"), exprType)
	}

	printSubHeader("Time Elapsed / Duration Formatter")
	earlier := ref.AddDate(-2, -3, -15)
	fmt.Printf("From %s to %s:\n", earlier.Format("2006-01-02"), ref.Format("2006-01-02"))
	fmt.Printf("  Short: %s\n", date.TimeElapsed(ref, earlier, false))
	fmt.Printf("  Full:  %s\n", date.TimeElapsed(ref, earlier, true))
}
