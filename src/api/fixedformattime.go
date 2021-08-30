package api

import (
	"errors"
	"time"
)

const fixedFormatTimeLayout = time.RFC3339

type FixedFormatTime time.Time

func (t *FixedFormatTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	var err error
	var v time.Time
	v, err = time.Parse(`"`+fixedFormatTimeLayout+`"`, string(data))
	*t = FixedFormatTime(v)
	return err
}

func (t FixedFormatTime) MarshalJSON() ([]byte, error) {
	v := time.Time(t)
	if y := v.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("FixedFormatTime.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(fixedFormatTimeLayout)+2)
	b = append(b, '"')
	b = v.AppendFormat(b, fixedFormatTimeLayout)
	b = append(b, '"')
	return b, nil
}
