package timeutils

import (
	"fmt"
	"github.com/sonnt85/gosutils/endec"
	"time"
)

type SleepStep struct {
	sleepDuration  time.Duration
	step, min, max time.Duration
}

type TickerStep struct {
	sleepDuration  time.Duration
	step, min, max time.Duration
	ticker         *time.Ticker
}

func NewSleepStep(step, min, max time.Duration) *SleepStep {
	return &SleepStep{
		max:           max,
		step:          step,
		min:           min,
		sleepDuration: min,
	}
}

func (ss *SleepStep) Sleep() {
	ss.sleepDuration = time.Duration(ss.sleepDuration.Nanoseconds() + endec.RandRangeInt64(0, int64(ss.step.Nanoseconds())))
	if ss.sleepDuration >= ss.max {
		ss.sleepDuration = ss.min
	}
	time.Sleep(ss.sleepDuration)
}

func NewTickerStep(step, min, max time.Duration) *TickerStep {
	return &TickerStep{
		max:           max,
		min:           min,
		step:          step,
		sleepDuration: min,
		ticker:        time.NewTicker(min),
	}
}

func (ss *TickerStep) C() <-chan time.Time {
	return ss.ticker.C
}

func (ss *TickerStep) Update() {
	ss.sleepDuration = time.Duration(ss.sleepDuration.Nanoseconds() + endec.RandRangeInt64(0, int64(ss.step.Nanoseconds())))
	if ss.sleepDuration >= ss.max {
		ss.sleepDuration = ss.min
	}
	ss.ticker.Reset(ss.sleepDuration)
}

// sleep random from mind to maxd duration
func SleepRandRange(mind, maxd time.Duration) {
	time.Sleep(time.Duration(endec.RandRangeInt64(int64(mind.Nanoseconds()), int64(maxd.Nanoseconds()))) * time.Nanosecond)
}

// sleep random from 0 to max duration
func SleepRandMax(d time.Duration) {
	SleepRandRange(0, d)
}

const day = time.Minute * 60 * 24

func StringDuration(d time.Duration) string {
	// if fm, err := durafmt.DefaultUnitsCoder.Decode("year:y,week:w,day:d,hour:h,minute:m,second:s,millisecond:ms,microsecond:us"); err == nil {
	// 	return durafmt.Parse(d).Format(fm)
	// }
	signed := ""
	if d < 0 {
		signed = "-"
		d *= -1
	}
	if d < day {
		return d.String()
	}

	n := d / day
	d -= n * day

	if d == 0 {
		return fmt.Sprintf("%s%dd", signed, n)
	}

	return fmt.Sprintf("%s%dd%s", signed, n, d)
}
