package types

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

type Date time.Time

func (p *Date) UnmarshalJSON(data []byte) error {

	if len(data) < 10 {
		*p = Date(time.Time{})
		return nil
	}
	local, err := time.ParseInLocation(`"2006-01-02"`, string(data), time.Local)

	*p = Date(local)

	return err
}

func (date *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = Date(nullTime.Time)
	return
}

func (date Date) Value() (driver.Value, error) {
	y, m, d := time.Time(date).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Time(date).Location()), nil
}

// GormDataType gorm common data type
func (date Date) GormDataType() string {
	return "date"
}

func (date Date) GobEncode() ([]byte, error) {
	return time.Time(date).GobEncode()
}

func (date *Date) GobDecode(b []byte) error {
	return (*time.Time)(date).GobDecode(b)
}

func (c Date) MarshalJSON() ([]byte, error) {
	data := make([]byte, 0)
	data = append(data, '"')
	data = time.Time(c).AppendFormat(data, "2006-01-02")
	data = append(data, '"')
	return data, nil
}

func (c Date) IsZero() bool {
	return time.Time(c).IsZero()
}
func (c Date) Time() time.Time {
	return time.Time(c)
}
func (c Date) String() string {
	return time.Time(c).Format("2006-01-02")
}
func (c Date) Unix() int64 {
	return time.Time(c).Unix()
}
func DateFromTime(t time.Time) Date {
	return Date(t)
}

func DateNow() Date {
	return Date(time.Now())
}
