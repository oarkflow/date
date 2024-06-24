package date

import (
	"math"
	"strconv"
	"strings"
	"time"
)

// day duration.
var day = time.Hour * 24

// week duration.
var week = time.Hour * 24 * 7

// Direction is the direction used for ambiguous expressions.
type Direction int

// Directions available.
const (
	Past Direction = iota
	Future
)

type ExprType int

func (t ExprType) IsTimeOnly() bool {
	return t == ExprTypeTime
}

func (t ExprType) IsDateOnly() bool {
	return t == ExprTypeDate
}

const (
	ExprTypeInvalid = 0
	ExprTypeDate    = ExprType(1 << iota)
	ExprTypeTime
	ExprTypeNow
	ExprTypeRelativeMinutes
	ExprTypeRelativeHours
	ExprTypeRelativeDays
	ExprTypeRelativeWeeks
	ExprTypeRelativeWeekdays
	ExprTypeRelativeMonth
	ExprTypeRelativeYear
	ExprTypeClock12Hour
	ExprTypeClock24Hour
)

// Option function.
type Option func(*gparser)

// WithDirection sets the direction used for ambiguous expressions. By default
// the Past direction is used, so "sunday" will be the previous Sunday, rather
// than the next Sunday.
func WithDirection(d Direction) Option {
	return func(p *gparser) {
		switch d {
		case Past:
			p.direction = -1
		case Future:
			p.direction = 1
		default:
			panic("unhandled direction")
		}
	}
}

// ParseNaturalDate query string.
func ParseNaturalDate(s string, ref time.Time, options ...Option) (time.Time, ExprType, error) {
	p := &gparser{
		Buffer:    strings.ToLower(s),
		direction: -1,
		t:         ref,
	}
	if s == "first day of the month" || s == "first day of this month" {
		return BeginningOfMonth(ref), p.exprType, nil
	}
	if s == "last day of the month" || s == "last day of this month" {
		return EndOfMonth(ref), p.exprType, nil
	}

	for _, o := range options {
		o(p)
	}

	p.Init()

	if err := p.Parse(); err != nil {
		return time.Time{}, ExprTypeInvalid, err
	}

	p.Execute()

	// p.PrintSyntaxTree()
	return p.t, p.exprType, nil
}

// withDirection returns duration with direction.
func (p *gparser) withDirection(d time.Duration) time.Duration {
	return d * time.Duration(p.direction)
}

func (p *gparser) dateExprSet(t time.Time) {
	p.exprType |= ExprTypeDate
	p.t = t
}

func (p *gparser) timeExprSet(t time.Time) {
	p.exprType |= ExprTypeTime
	p.t = t
}

// prevWeekday returns the previous week day relative to time t.
func prevWeekday(t time.Time, day time.Weekday) time.Time {
	d := t.Weekday() - day
	if d <= 0 {
		d += 7
	}
	return t.Add(-time.Hour * 24 * time.Duration(d))
}

// nextWeekday returns the next week day relative to time t.
func nextWeekday(t time.Time, day time.Weekday) time.Time {
	d := day - t.Weekday()
	if d <= 0 {
		d += 7
	}
	return t.Add(time.Hour * 24 * time.Duration(d))
}

// nextMonth returns the next month relative to time t.
func nextMonth(t time.Time, month time.Month) time.Time {
	y := t.Year()
	if month-t.Month() <= 0 {
		y++
	}
	_, _, day := t.Date()
	hour, min, sec := t.Clock()
	return time.Date(y, month, day, hour, min, sec, 0, t.Location())
}

// prevMonth returns the next month relative to time t.
func prevMonth(t time.Time, month time.Month) time.Time {
	y := t.Year()
	if t.Month()-month <= 0 {
		y--
	}
	_, _, day := t.Date()
	hour, min, sec := t.Clock()
	return time.Date(y, month, day, hour, min, sec, 0, t.Location())
}

// truncateDay returns a date truncated to the day.
func truncateDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func s(x float64) string {
	if int(x) == 1 {
		return ""
	}
	return "s"
}

func TimeElapsed(now time.Time, then time.Time, full bool) string {
	var parts []string
	var text string

	year2, month2, day2 := now.Date()
	hour2, minute2, second2 := now.Clock()

	year1, month1, day1 := then.Date()
	hour1, minute1, second1 := then.Clock()

	year := math.Abs(float64(year2 - year1))
	month := math.Abs(float64(month2 - month1))
	day := math.Abs(float64(day2 - day1))
	hour := math.Abs(float64(hour2 - hour1))
	minute := math.Abs(float64(minute2 - minute1))
	second := math.Abs(float64(second2 - second1))

	week := math.Floor(day / 7)

	if year > 0 {
		parts = append(parts, strconv.Itoa(int(year))+" year"+s(year))
	}

	if month > 0 {
		parts = append(parts, strconv.Itoa(int(month))+" month"+s(month))
	}

	if week > 0 {
		parts = append(parts, strconv.Itoa(int(week))+" week"+s(week))
	}

	if day > 0 {
		parts = append(parts, strconv.Itoa(int(day))+" day"+s(day))
	}

	if hour > 0 {
		parts = append(parts, strconv.Itoa(int(hour))+" hour"+s(hour))
	}

	if minute > 0 {
		parts = append(parts, strconv.Itoa(int(minute))+" minute"+s(minute))
	}

	if second > 0 {
		parts = append(parts, strconv.Itoa(int(second))+" second"+s(second))
	}

	if now.After(then) {
		text = " ago"
	} else {
		text = " after"
	}

	if len(parts) == 0 {
		return "just now"
	}

	if full {
		return strings.Join(parts, ", ") + text
	}
	return parts[0] + text
}

func BeginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func EndOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}
