package calendar

import (
	"time"

	"github.com/HenryTK/calendar-resource/errors"
	"github.com/HenryTK/calendar-resource/models"
)

type Event struct {
	StartTime time.Time
	EndTime   time.Time
}

type TimeKeeper interface {
	Now() time.Time
	IsHappeningNow(time.Time, time.Time) bool
}

type Timer struct{}

func (t Timer) Now() time.Time {
	return time.Now()
}

// IsHappeningNow takes an event start and end time and returns true if the
// event is currently happening
func (t Timer) IsHappeningNow(eventStartTime, eventEndTime time.Time) bool {
	if t.Now().After(eventStartTime) && t.Now().Before(eventEndTime) {
		return true
	}
	return false
}

type Calendar struct {
	Events []Event
	TimeKeeper
}

func NewCalendar() *Calendar {
	return &Calendar{
		[]Event{},
		Timer{},
	}
}

// CurrentVersions takes Calendar.Events (which has already been populated
// by the Calendar.Client) and sorts out which events are currently happening.
// It satisfies the `check` part of the Concourse resource interface. It will
// take the version passed via standard input to the check script and return a
// list containing just that version if it is still current, otherwise it will
// return a list of those events which are currently happening.
func (c *Calendar) CurrentVersions(version models.Version) []models.Version {
	if version.StartTime != "" && version.EndTime != "" {
		if c.IsHappeningNow(StringToTime(version.StartTime), StringToTime(version.EndTime)) {
			currentVersion := models.Version{StartTime: version.StartTime, EndTime: version.EndTime}
			return []models.Version{currentVersion}
		}
	}
	versions := []models.Version{}
	for _, event := range c.Events {
		if c.IsHappeningNow(event.StartTime, event.EndTime) {
			versions = append(versions, models.Version{
				StartTime: TimeToString(event.StartTime),
				EndTime:   TimeToString(event.EndTime),
			})
		}
	}
	return versions
}

// TimeToString takes a time and converts
// it to the RFC 3339 format. For example:
// 2016-10-02T14:00:00+01:00
func TimeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}

// StringToTime takes a string in the RFC
// 3339 format and converts it to a Golang
// Time. Example formatted string:
// 2016-10-02T14:00:00+01:00
func StringToTime(str string) time.Time {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		errors.Fatal("Failed parsing time: ", err)
	}
	return t

}
