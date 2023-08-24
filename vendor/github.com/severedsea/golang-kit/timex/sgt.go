package timex

import (
	"log"
	"time"
)

// SGT is the timezone for Singapore. Used e.g. `createdAt.In(utils.SGT)`
var SGT *time.Location

const (
	startTimestamp  = "T00:00:00+08:00"
	endTimestamp    = "T23:59:59+08:00"
	midDayTimestamp = "T11:59:59+08:00"
)

func init() {
	var err error
	SGT, err = time.LoadLocation("Singapore")
	if err != nil {
		log.Fatalf(err.Error())
	}
}

// NowSGT returns time.Now() in SGT
func NowSGT() time.Time {
	return time.Now().In(SGT)
}

// DateSGT returns the Time in SGT using time.Date()
func DateSGT(year int, month time.Month, day, hour, min int) time.Time {
	return time.Date(year, month, day, hour, min, 0, 0, SGT)
}

// FormatSGT formats the time provided into the provided layout in SGT
func FormatSGT(t time.Time, layout string) string {
	return t.In(SGT).Format(layout)
}

// BeginningOfDay sets 00:00:00 SGT as time
func BeginningOfDaySGT(now time.Time) time.Time {
	now = now.In(SGT)
	t, _ := time.Parse(time.RFC3339, now.Format("2006-01-02")+startTimestamp)

	return t
}

// EndOfDay sets 23:59:59 SGT as time
func EndOfDaySGT(now time.Time) time.Time {
	now = now.In(SGT)
	t, _ := time.Parse(time.RFC3339, now.Format("2006-01-02")+endTimestamp)

	return t
}

// MidOfDay sets 11:59:59 SGT as time
func MidOfDaySGT(now time.Time) time.Time {
	now = now.In(SGT)
	t, _ := time.Parse(time.RFC3339, now.Format("2006-01-02")+midDayTimestamp)

	return t
}
