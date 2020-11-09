package model

import (
	"time"
)

type TimeStamp time.Time

const (
	jsonLayout = "2006-01-02 15:04:05"
)

func Now() TimeStamp {
	return TimeStamp(time.Now())
}

func (t *TimeStamp) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "" || string(data) == `""` {
		return nil
	}
	now, err := time.ParseInLocation(`"`+jsonLayout+`"`, string(data), time.Local)
	if err == nil {
		*t = TimeStamp(now)
	}
	return nil
}

func (t TimeStamp) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(jsonLayout)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, jsonLayout)
	b = append(b, '"')
	return b, nil
}
