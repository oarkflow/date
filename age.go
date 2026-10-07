package date

import (
	"math"
	"time"
)

const (
	calculateLeapYearDay = 60
)

func calculate(startTime, endTime time.Time) int {
	startYear := startTime.Year()
	endYear := endTime.Year()

	age := endYear - startYear

	startYearIsLeapYear := IsLeapYear(startYear)
	endYearIsLeapYear := IsLeapYear(endYear)

	startYearDay := startTime.YearDay()
	endYearDay := endTime.YearDay()

	if startYearIsLeapYear && !endYearIsLeapYear && startYearDay >= calculateLeapYearDay {
		startYearDay--
	} else if endYearIsLeapYear && !startYearIsLeapYear && endYearDay >= calculateLeapYearDay {
		startYearDay++
	}

	if endYearDay < startYearDay {
		age--
	}

	return age
}

// Calculate returns an integer-value age based on the duration
// between the two times that are given as arguments.
func Calculate(startTime, endTime time.Time) int {
	switch endLocation := endTime.Location(); endLocation {
	case time.UTC, nil:
		startTime = startTime.UTC()
	default:
		startTime = startTime.In(endLocation)
	}

	return calculate(startTime, endTime)
}

// CalculateToNow returns an integer-value age based on the duration
// between the given time and the present time.
func CalculateToNow(givenTime time.Time) int {
	presentTime := time.Now().UTC()

	switch givenLocation := givenTime.Location(); givenLocation {
	case time.UTC, nil:
		presentTime = presentTime.UTC()
	default:
		presentTime = presentTime.In(givenLocation)
	}

	return calculate(givenTime, presentTime)
}

// IsLeapYear returns true only if the given year contains a leap day,
// meaning that the year is a leap year.
func IsLeapYear(givenYear int) bool {
	if givenYear%400 == 0 {
		return true
	} else if givenYear%100 == 0 {
		return false
	} else if givenYear%4 == 0 {
		return true
	}

	return false
}

// PrevLeapYear returns the previous leap year before a given year.
// It also returns a boolean value that is set to false if no
// such year can be found.
func PrevLeapYear(givenYear int) (foundYear int, found bool) {
	for foundYear = givenYear; foundYear > math.MinInt; {
		foundYear--

		if found = IsLeapYear(foundYear); found {
			return
		}
	}

	return
}

// NextLeapYear returns the next leap year after a given year.
// It also returns a boolean value that is set to false if no
// such year can be found.
func NextLeapYear(givenYear int) (foundYear int, found bool) {
	for foundYear = givenYear; foundYear < math.MaxInt; {
		foundYear++

		if found = IsLeapYear(foundYear); found {
			return
		}
	}

	return
}

// AgeDetailed returns the exact difference in years, months, and days between birthDate and targetDate.
// If targetDate is before birthDate, they are swapped.
func AgeDetailed(birthDate, targetDate time.Time) (years, months, days int) {
	if targetDate.Before(birthDate) {
		birthDate, targetDate = targetDate, birthDate
	}
	y1, m1, d1 := birthDate.Date()
	y2, m2, d2 := targetDate.Date()

	years = y2 - y1
	months = int(m2) - int(m1)
	days = d2 - d1

	if days < 0 {
		// Borrow days from previous month
		prevMonth := time.Date(y2, m2, 0, 0, 0, 0, 0, targetDate.Location())
		days += prevMonth.Day()
		months--
	}
	if months < 0 {
		months += 12
		years--
	}
	return years, months, days
}

// AgeDetailedToNow returns the exact difference in years, months, and days between birthDate and now.
func AgeDetailedToNow(birthDate time.Time) (years, months, days int) {
	return AgeDetailed(birthDate, time.Now())
}
