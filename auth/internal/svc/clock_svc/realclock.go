package clock_svc

import "time"

type ClockInterface interface {
	Now() time.Time
}

type RealClockStruct struct{}

func (RealClockStruct) Now() time.Time {
	return time.Now()
}
