package timekit

import (
	"regexp"
	"time"
)

type Layout = string

const YYYYMMDDhhmmss Layout = "2006-01-02 15:04:05"
const YYYYMMDDhhmmsspure Layout = "20060102150405"
const YYYYMMDD Layout = "2006-01-02"

var maptimeformate = map[*regexp.Regexp]string{
	regexp.MustCompile(`\d{4}-\d{2}-\d{2}\s{1}\d{2}\:\d{2}:\d{2}$`):  YYYYMMDDhhmmss,
	regexp.MustCompile(`\d{4}-\d{2}-\d{2}[T]{1}\d{2}\:\d{2}:\d{2}$`): YYYYMMDDhhmmss,
	regexp.MustCompile(`\d{4}/\d{2}/\d{2}\s{1}\d{2}\:\d{2}:\d{2}$`):  YYYYMMDDhhmmss,
	regexp.MustCompile(`\d{4}/\d{2}/\d{2}[T]{1}\d{2}\:\d{2}:\d{2}$`): YYYYMMDDhhmmss,
	regexp.MustCompile(`\d{4}/\d{2}/\d{2}$`):                         YYYYMMDD,
	regexp.MustCompile(`\d{4}-\d{2}-\d{2}$`):                         YYYYMMDD,
}

func Parse(input string, layouts ...Layout) (tm time.Time, err error) {
	if len(layouts) == 0 {
		layouts = append(layouts, YYYYMMDDhhmmss)
	}
	for _, layout := range layouts {
		tm, err = time.Parse(layout, input)
	}
	return
}
func Format(tm time.Time, layout Layout) string {
	return tm.Format(layout)
}

func DateTimeNow(layouts ...Layout) string {
	if len(layouts) == 0 {
		layouts = append(layouts, YYYYMMDDhhmmss)
	}
	r := ""
	for _, l := range layouts {
		r = Format(time.Now(), l)
	}
	return r
}
func DateNow() string {
	return Format(time.Now(), YYYYMMDD)
}
