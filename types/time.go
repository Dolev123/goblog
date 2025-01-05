package types

import (
    "time"
)

type Time struct {
    time.Time
}

// based on time.DateTime
const TimeFormat = `"2006-01-02 15:04:05"`

func (t *Time) UnmarshalJSON(b []byte) error {
    date, err := time.Parse(TimeFormat, string(b))
    if err != nil {
	return err
    }
    t.Time = date
    return nil
}
