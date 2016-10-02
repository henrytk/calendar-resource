package fakes

import (
	"time"
)

type FakeTimer struct {
	FakeNow time.Time
}

func (fk FakeTimer) Now() time.Time {
	return fk.FakeNow
}

// IsHappeningNow takes an event start and end time and returns true if the
// event is currently happening
func (ft FakeTimer) IsHappeningNow(eventStartTime, eventEndTime time.Time) bool {
	if ft.Now().After(eventStartTime) && ft.Now().Before(eventEndTime) {
		return true
	}
	return false
}

func NewFakeTimer(t time.Time) FakeTimer {
	return FakeTimer{FakeNow: t}
}
