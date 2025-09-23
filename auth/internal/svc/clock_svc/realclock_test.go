package clock_svc

import "testing"

func TestNow(t *testing.T) {
	realClock := RealClockStruct{}

	currentTime := realClock.Now()

	if currentTime.IsZero() {
		t.Error("expected current time to be non-zero")
	}
}
