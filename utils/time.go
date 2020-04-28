package utils

import (
	"database/sql/driver"
	"github.com/araddon/dateparse"
	"strings"
	"time"
)

type Time time.Time

const DatabaseTimeFormat = "2006-01-02 15:04:05"

// 实现了 UnmarshalJSON 。 json解码时用此方法
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	// 常常遇到空字符串 "" 无法解析，
	if string(data) == "null" || len(data) < 4 {
		t = nil
		return nil
	}
	//  这里使用 dateparse 来解析时间格式
	tt, err := dateparse.ParseAny(strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	*t = Time(tt)
	return err
}

// 实现了 MarshalJSON 。 json编码时用此方法
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(time.RFC3339)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, time.RFC3339)
	b = append(b, '"')
	return b, nil
}

// implements Stringer interface
func (t Time) String() string {
	return time.Time(t).String()
}

// Value 和 Scan 方法用于数据库操作
// implements the sql.Valuer interface
func (t Time) Value() (driver.Value, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return nil, nil
	}
	return tt.Format(DatabaseTimeFormat), nil
}

// Scan implements the sql.Scanner interface
func (t *Time) Scan(src interface{}) error {
	if src == nil {
		t := Time(time.Time{})
		_ = t
		return nil
	}
	// 数据库存的总是 UTC 时间， 这里可以转换本地时间
	// localtime := src.(time.Time).In(time.Local)
	*t = Time(src.(time.Time))
	return nil
}
