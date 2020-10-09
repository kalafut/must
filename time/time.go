package time

import (
	"time"

	"github.com/kalafut/must"
)

type Time time.Time

func Parse(layout, value string) Time {
	ret, err := time.Parse(layout, value)
	must.PanicErr(err)

	return (Time)(ret)

}

func (t Time) MarshalText() []byte {
	ret, err := (time.Time)(t).MarshalText()
	must.PanicErr(err)

	return ret
}
