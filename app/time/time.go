package time

import (
	"context"
	"database/sql/driver"
	"fmt"
	"time"
)

var (
	CSTZone = time.FixedZone("CST", 8*3600)       // 东八
	_ = time.Now().In(CSTZone).Format(StandardLayout)
)

const StandardLayout = "2006-01-02 15:04:05"
const StandardYmdLayout = "2006-01-02"

// StandardTime alias time.Time
type Time DateTime

// StandardTime alias time.Time
type DateTime struct {
	time.Time
}

// Value insert timestamp into mysql need this function.
func (it DateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if it.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return it.Time, nil
}

// Scan valueof time.Time
func (it *DateTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*it = DateTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// MarshalJSON implements the json.Marshaler interface.
func (it DateTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(StandardLayout)+2)
	b = append(b, '"')
	b = Time(it).In(CSTZone).AppendFormat(b, StandardLayout)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (it *DateTime) UnmarshalJSON(data []byte) error {
	t, err := time.ParseInLocation(`"`+StandardLayout+`"`, string(data), CSTZone)
	if err == nil {
		*it = DateTime{Time: t,}
	}
	return err
}


// StandardTime duration
type Duration time.Duration

// UnmarshalText unmarshal text to duration.
func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := time.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp)
	}
	return err
}

// Shrink will decrease the duration by comparing with context's timeout duration
// and return new timeout\context\CancelFunc.
func (d Duration) Shrink(c context.Context) (Duration, context.Context, context.CancelFunc) {
	if deadline, ok := c.Deadline(); ok {
		if ctimeout := time.Until(deadline); ctimeout < time.Duration(d) {
			// deliver small timeout
			return Duration(ctimeout), c, func() {}
		}
	}
	ctx, cancel := context.WithTimeout(c, time.Duration(d))
	return d, ctx, cancel
}
