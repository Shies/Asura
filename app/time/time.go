package time

import (
	"time"
	"context"
)

var (
	CSTZone = time.FixedZone("CST", 8*3600)       // 东八
	_ = time.Now().In(CSTZone).Format(StandardLayout)
)

// ISO8601Layout time parse layout
const ISO8601Layout = "2006-01-02T15:04:05.000-0700"
const StandardLayout = "2006-01-02 15:04:05"
const StandardYmdLayout = "2006-01-02"

// ISO8601Time alias time.Time
type ISO8601Time time.Time

// StandardTime alias time.Time
type DateTime time.Time

// MarshalJSON implements the json.Marshaler interface.
func (it ISO8601Time) MarshalJSON() ([]byte, error) {
	return time.Time(it).MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (it *ISO8601Time) UnmarshalJSON(data []byte) (err error) {
	t, err := time.Parse(`"`+ISO8601Layout+`"`, string(data))
	*it = ISO8601Time(t)
	return
}

func (it ISO8601Time) String() string {
	return time.Time(it).Format(ISO8601Layout)
}

// MarshalJSON implements the json.Marshaler interface.
func (it DateTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(StandardLayout)+2)
	b = append(b, '"')
	b = time.Time(it).In(CSTZone).AppendFormat(b, StandardLayout)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (it *DateTime) UnmarshalJSON(data []byte) (err error) {
	t, err := time.ParseInLocation(`"`+StandardLayout+`"`, string(data), CSTZone)
	*it = DateTime(t)
	return
}

func (it DateTime) String() string {
	return time.Time(it).In(CSTZone).Format(StandardLayout)
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