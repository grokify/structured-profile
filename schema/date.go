package schema

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Date represents a date in YYYYMM or YYYYMMDD format.
// It supports both 6-digit (YYYYMM) and 8-digit (YYYYMMDD) formats.
type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day,omitempty"`
}

// NewDate creates a Date from year and month.
func NewDate(year, month int) Date {
	return Date{Year: year, Month: month}
}

// NewDateFull creates a Date from year, month, and day.
func NewDateFull(year, month, day int) Date {
	return Date{Year: year, Month: month, Day: day}
}

// NewDateFromTime creates a Date from a time.Time.
func NewDateFromTime(t time.Time) Date {
	return Date{
		Year:  t.Year(),
		Month: int(t.Month()),
		Day:   t.Day(),
	}
}

// NewDateFromDT6 creates a Date from a DT6 integer (YYYYMM format).
func NewDateFromDT6(dt6 int) Date {
	year := dt6 / 100
	month := dt6 % 100
	return Date{Year: year, Month: month}
}

// NewDateFromDT8 creates a Date from a DT8 integer (YYYYMMDD format).
func NewDateFromDT8(dt8 int) Date {
	year := dt8 / 10000
	month := (dt8 / 100) % 100
	day := dt8 % 100
	return Date{Year: year, Month: month, Day: day}
}

// DT6 returns the date as a 6-digit integer (YYYYMM).
func (d Date) DT6() int {
	return d.Year*100 + d.Month
}

// DT8 returns the date as an 8-digit integer (YYYYMMDD).
func (d Date) DT8() int {
	day := d.Day
	if day == 0 {
		day = 1
	}
	return d.Year*10000 + d.Month*100 + day
}

// Time returns the Date as a time.Time.
func (d Date) Time() time.Time {
	day := d.Day
	if day == 0 {
		day = 1
	}
	return time.Date(d.Year, time.Month(d.Month), day, 0, 0, 0, 0, time.UTC)
}

// IsZero returns true if the date is unset.
func (d Date) IsZero() bool {
	return d.Year == 0 && d.Month == 0 && d.Day == 0
}

// String returns the date as a string in "YYYY-MM" or "YYYY-MM-DD" format.
func (d Date) String() string {
	if d.Day > 0 {
		return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
	}
	return fmt.Sprintf("%04d-%02d", d.Year, d.Month)
}

// DisplayString returns a human-readable date string.
// If the date is zero, it returns "Present".
func (d Date) DisplayString() string {
	if d.IsZero() {
		return "Present"
	}
	months := []string{
		"", "Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
	}
	if d.Month >= 1 && d.Month <= 12 {
		return fmt.Sprintf("%s %d", months[d.Month], d.Year)
	}
	return fmt.Sprintf("%d", d.Year)
}

// Before returns true if d is before other.
func (d Date) Before(other Date) bool {
	return d.DT8() < other.DT8()
}

// After returns true if d is after other.
func (d Date) After(other Date) bool {
	return d.DT8() > other.DT8()
}

// Equal returns true if d equals other.
func (d Date) Equal(other Date) bool {
	return d.Year == other.Year && d.Month == other.Month && d.Day == other.Day
}

// MarshalJSON implements json.Marshaler.
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		// Try unmarshaling as integer (DT6 or DT8 format)
		var n int
		if err := json.Unmarshal(data, &n); err != nil {
			return fmt.Errorf("date must be string or integer: %w", err)
		}
		if n > 999999 {
			*d = NewDateFromDT8(n)
		} else {
			*d = NewDateFromDT6(n)
		}
		return nil
	}

	// Parse string format
	if len(s) == 0 {
		*d = Date{}
		return nil
	}

	// Try YYYY-MM-DD format
	if len(s) == 10 {
		t, err := time.Parse("2006-01-02", s)
		if err == nil {
			*d = NewDateFromTime(t)
			return nil
		}
	}

	// Try YYYY-MM format
	if len(s) == 7 {
		// Handle zero date
		if s == "0000-00" {
			*d = Date{}
			return nil
		}
		t, err := time.Parse("2006-01", s)
		if err == nil {
			*d = Date{Year: t.Year(), Month: int(t.Month())}
			return nil
		}
	}

	// Try YYYYMM format
	if len(s) == 6 {
		n, err := strconv.Atoi(s)
		if err == nil {
			*d = NewDateFromDT6(n)
			return nil
		}
	}

	// Try YYYYMMDD format
	if len(s) == 8 {
		n, err := strconv.Atoi(s)
		if err == nil {
			*d = NewDateFromDT8(n)
			return nil
		}
	}

	return fmt.Errorf("invalid date format: %s", s)
}

// DateRange represents a range of dates with start and optional end.
type DateRange struct {
	Start Date  `json:"start"`
	End   *Date `json:"end,omitempty"`
}

// IsCurrent returns true if the date range has no end date (ongoing).
func (r DateRange) IsCurrent() bool {
	return r.End == nil || r.End.IsZero()
}

// Duration returns the approximate duration in months.
func (r DateRange) Duration() int {
	end := r.End
	if end == nil || end.IsZero() {
		now := NewDateFromTime(time.Now())
		end = &now
	}
	return (end.Year-r.Start.Year)*12 + (end.Month - r.Start.Month)
}

// String returns the date range as a string.
func (r DateRange) String() string {
	if r.IsCurrent() {
		return fmt.Sprintf("%s - Present", r.Start.DisplayString())
	}
	return fmt.Sprintf("%s - %s", r.Start.DisplayString(), r.End.DisplayString())
}
