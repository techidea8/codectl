package types

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

type DateTime time.Time

func (p *DateTime) UnmarshalJSON(data []byte) error {

	if len(data) < 10 {
		*p = DateTime(time.Time{})
		return nil
	}

	local, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	*p = DateTime(local)

	return err
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *DateTime) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*j = DateTime(nullTime.Time)
	return
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (date DateTime) Value() (driver.Value, error) {
	d := time.Time(date)
	return time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second(), 0, time.Time(date).Location()), nil
}

// GormDataType gorm common data type
func (date DateTime) GormDataType() string {
	return "datetime"
}
func (c DateTime) MarshalJSON() ([]byte, error) {
	data := make([]byte, 0)
	data = append(data, '"')
	data = time.Time(c).AppendFormat(data, "2006-01-02 15:04:05")
	data = append(data, '"')
	return data, nil
}

func (c DateTime) IsZero() bool {
	return time.Time(c).IsZero()
}

func (c DateTime) Unix() int64 {
	return time.Time(c).Unix()
}

func (c DateTime) String() string {
	return time.Time(c).Format("2006-01-02 15:04:05")
}
func (c DateTime) Now() DateTime {
	return DateTime(time.Now())
}
func (c DateTime) Time() time.Time {
	return time.Time(c)
}
func DateTimeFromTime(t time.Time) DateTime {
	return DateTime(t)
}
func DateTimeNow() DateTime {
	return DateTime(time.Now())
}

func (c DateTime) FormatDay() string {
	return time.Time(c).Format("2006-01-02")
}
func (c DateTime) FormatMonth() string {

	return time.Time(c).Format("2006-01")
}
