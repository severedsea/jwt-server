package timex

import "time"

// UnixMillis returns the unix time in millies for the provided timestamp
func UnixMillis(ts time.Time) int64 {
	return ts.UnixNano() / int64(time.Millisecond)
}
