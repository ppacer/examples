package main

import (
	"fmt"
	"time"

	"github.com/ppacer/core/dag/schedule"
)

type MyCustomSched struct {
	start     time.Time
	dailyCron *schedule.Cron
}

func NewCustomSched(hour, minute int, start time.Time) *MyCustomSched {
	cron := schedule.NewCron().AtHour(hour).AtMinute(minute)
	return &MyCustomSched{
		start:     start,
		dailyCron: cron,
	}
}

func (mcs *MyCustomSched) Start() time.Time { return mcs.start }

func (mcs *MyCustomSched) String() string {
	return fmt.Sprintf("MyCustomSched: %s", mcs.dailyCron.String())
}

func (mcs *MyCustomSched) Next(currentTime time.Time, _ *time.Time) time.Time {
	cronNext := mcs.dailyCron.Next(currentTime, nil)
	if cronNext.Month() == time.June || cronNext.Month() == time.August {
		if cronNext.Weekday() == time.Friday {
			return mcs.dailyCron.Next(cronNext, nil)
		}
	}
	if cronNext.Month() == time.December && cronNext.Day() == 24 {
		return time.Date(
			cronNext.Year(), cronNext.Month(), cronNext.Day(), 8, 0, 0, 0,
			cronNext.Location(),
		)
	}
	return cronNext
}

func main() {
	start := time.Date(2024, 4, 1, 8, 0, 0, 0, time.UTC)
	mySched := NewCustomSched(10, 15, start)

	rand := time.Date(2024, time.April, 2, 7, 0, 0, 0, time.UTC)
	summerFriday := time.Date(2024, time.August, 9, 9, 30, 0, 0, time.UTC)
	beforeXMas := time.Date(2024, time.December, 23, 12, 0, 0, 0, time.UTC)

	fmt.Printf("Next sched for %v: %v\n", rand, mySched.Next(rand, nil))
	fmt.Printf("Next sched for Friday in August (%v): %v\n", summerFriday,
		mySched.Next(summerFriday, nil))
	fmt.Printf("Next sched for one day before XMas (%v): %v\n", beforeXMas,
		mySched.Next(beforeXMas, nil))
}
