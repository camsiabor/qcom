package qtime

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func ShiftUTCPointer(year, month, day int, delta, trunc time.Duration) *time.Time {

	var t = time.Now().UTC()
	if year != 0 || month != 0 || day != 0 {
		t = t.AddDate(year, month, day)
	}

	if delta != 0 {
		t = t.Add(delta)
	}

	if trunc != 0 {
		t = t.Truncate(trunc)
	}

	return &t
}

func NowUTCPointer() *time.Time {
	var t = time.Now().UTC()
	return &t
}

func NowUTCRFC33399() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func NowUTCRFC33399Shift(year, month, day int) string {
	return time.Now().UTC().AddDate(year, month, day).Format(time.RFC3339)
}

func Date2Int(t *time.Time) int {
	if t == nil {
		var n = time.Now()
		t = &n
	}
	var y, m, d = t.Date()
	return int(y*10000) + int(m*100) + d
}

func Time2Int64(t *time.Time) int64 {
	if t == nil {
		var n = time.Now()
		t = &n
	}
	var y, m, d = t.Date()
	var hour, min, sec = t.Clock()
	var date = int64(y*10000000000) + int64(m*100000000) + int64(d*1000000)
	var clock int64 = int64(hour*10000) + int64(min*100) + int64(sec)
	return date + clock
}

func ParseTime(s string) (*time.Time, error) {
	var err error
	var t *time.Time
	if strings.Contains(s, "-") {
		if strings.Contains(s, ":") {
			if len(s) > 8 {
				*t, err = time.Parse("2006-01-02 15:04:05", s)
			} else {
				*t, err = time.Parse("15:04:05", s)
			}
		} else {
			*t, err = time.Parse("2006-01-02", s)
		}
	} else if strings.Contains(s, "/") {
		if strings.Contains(s, ":") {
			if len(s) > 8 {
				*t, err = time.Parse("2006/01/02 15:04:05", s)
			} else {
				*t, err = time.Parse("15:04:05", s)
			}
		} else {
			*t, err = time.Parse("2006/01/02", s)
		}
	} else {
		if len(s) > 8 {
			*t, err = time.Parse("20060102 150405", s)
		} else {
			*t, err = time.Parse("20060102", s)
		}
	}
	if err != nil {
		return nil, err
	}
	return t, err
}

func TimeInterval(t1 *time.Time, t2 *time.Time, duration time.Duration) float64 {
	if t2 == nil {
		var n = time.Now()
		t2 = &n
	}
	var interval = float64(t1.UnixNano() - t2.UnixNano())
	var d = float64(duration)
	return math.Ceil(interval / d)
}

func YYYY_MM_dd_HH_mm_ss(t *time.Time) string {
	if t == nil {
		var s = time.Now()
		t = &s
	}
	return t.Format("2006-01-02 15:04:05")
}

func YYYY_MM_dd(t *time.Time) string {
	if t == nil {
		var s = time.Now()
		t = &s
	}
	return t.Format("2006-01-02")
}

func GetTimeFormatIntervalArray(from *time.Time, to *time.Time, layout string, lastday bool, exclude ...time.Weekday) ([]string, error) {

	if to == nil {
		var s = time.Now()
		to = &s
	}

	var i = 0
	var time_array []string = make([]string, 256)
	time_array[i] = from.Format(layout)
	i = i + 1

	var fromnum = Date2Int(from)
	var tonum = Date2Int(to)
	if fromnum > tonum {
		return nil, fmt.Errorf("from time > to time")
	}

	var excludearray []int
	var doexclude = exclude != nil && len(exclude) > 0
	if doexclude {
		excludearray = make([]int, 7)
		for i := 0; i < len(exclude); i++ {
			var one = exclude[i]
			excludearray[one] = 1
		}
	}

	if fromnum < tonum {
		var middle = from.AddDate(0, 0, 0)
		for {
			middle = middle.AddDate(0, 0, 1)
			var middlenum = Date2Int(&middle)
			if middlenum > tonum {
				break
			}

			if doexclude {
				if excludearray[middle.Weekday()] != 0 {
					continue
				}
			}

			if lastday {
				var day = middle.Day()
				if day >= 28 {
					var next = middle.AddDate(0, 0, 1)
					if next.Month() == middle.Month() {
						continue
					}
				} else {
					continue
				}
			}

			time_array[i] = middle.Format(layout)
			i = i + 1

		}
	}

	return time_array[:i], nil
}

func TruncateHMS(t *time.Time) *time.Time {
	var delta = time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute + time.Duration(t.Second())*time.Second
	*t = t.Add(-delta)
	return t
}
