package timex

import (
	"math"
	"time"
)

// SecondsSince returns the seconds since start time relative to end time
func SecondsSince(start, end time.Time) int {
	s := end.In(SGT).Sub(start).Seconds()

	return int(math.Round(s))
}

// AfterNow returns whether input is after now)
func AfterNow(now, input time.Time) bool {
	return input.After(now.In(SGT))
}
