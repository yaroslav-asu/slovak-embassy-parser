package datetime

import (
	"database/sql"
	"database/sql/driver"
	"go.uber.org/zap"
	"time"
)

type Date time.Time

var (
	MonthAndYear  = "01.2006"
	DateOnly      = "02.01.2006"
	DateTime      = "15:04 02.01.2006"
	FormDateTime  = "02.01.2006 15:04"
	SiteMonthYear = "01/2006"
)

// Scan implementation for Gorm
func (d *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*d = Date(nullTime.Time)
	return
}

// Value implementation for Gorm
func (d *Date) Value() (driver.Value, error) {
	year, month, date := time.Time(*d).Date()
	return time.Date(year, month, date, 0, 0, 0, 0, time.UTC), nil
}

func (d *Date) GormDataType() string {
	return "time"
}

func (d *Date) GobEncode() ([]byte, error) {
	return time.Time(*d).GobEncode()
}

func (d *Date) GobDecode(b []byte) error {
	return (*time.Time)(d).GobDecode(b)
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return time.Time(*d).MarshalJSON()
}

func (d *Date) UnmarshalJSON(b []byte) error {
	return (*time.Time)(d).UnmarshalJSON(b)
}

func Now() Date {
	return Date(time.Now())
}

func NewDate(year int, month int, day int, hour int, minute int) Date {
	return Date(time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC))
}

func NewDateYMD(year int, month int, day int) Date {
	return NewDate(year, month, day, 0, 0)
}

func NewDateYM(year int, month int) Date {
	return NewDate(year, month, 1, 0, 0)
}

func NewBlankDate() Date {
	return NewDate(1, 1, 1, 0, 0)
}

func (d *Date) ChangeMinutes(minutes int) {
	d.Set(minutes, d.Hour(), d.Day(), d.Month(), d.Year())
}

func (d *Date) SetHour(hour int) {
	d.Set(d.Minute(), hour, d.Day(), d.Month(), d.Year())
}

func (d *Date) SetDay(day int) {
	d.Set(d.Minute(), d.Hour(), day, d.Month(), d.Year())
}

func (d *Date) SetMonth(month int) Date {
	return d.Set(d.Minute(), d.Hour(), d.Day(), month, d.Year())
}

func (d *Date) MoveMonth(delta int) Date {
	return d.SetMonth(d.Month() + delta)
}

func (d *Date) SetYear(year int) {
	a := time.Time(*d)
	d.Set(a.Minute(), d.Hour(), d.Day(), d.Month(), year)
}

func (d *Date) Set(minutes int, hour int, day int, month int, year int) Date {
	*d = Date(time.Date(year, time.Month(month), day, hour, minutes, 0, 0, time.UTC))
	return *d
}

func (d *Date) Minute() int {
	return time.Time(*d).Minute()
}

func (d *Date) Hour() int {
	return time.Time(*d).Hour()
}

func (d *Date) Day() int {
	return time.Time(*d).Day()
}

func (d *Date) Month() int {
	return int(time.Time(*d).Month())
}

func (d *Date) Year() int {
	return time.Time(*d).Year()
}

func (d *Date) Time() time.Time {
	return time.Time(*d)
}

func (d *Date) Format(format string) string {
	return d.Time().Format(format)
}

func ParseDateFromString(formatString, dateString string) (Date, error) {
	parsedTime, err := time.Parse(formatString, dateString)
	if err != nil {
		zap.L().Error("Failed to parse date from: " + dateString + " by: " + formatString)
		return Date{}, err
	}
	return Date(parsedTime), nil
}
