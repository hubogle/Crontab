package json

import (
	"fmt"
	"time"
)

type Time time.Time

// MarshalJson 直接转换 time.Time 为 json 字符串
func (j Time) MarshalJson() ([]byte, error) {
	var tmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))
	return []byte(tmp), nil
}
