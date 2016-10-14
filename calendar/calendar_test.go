package calendar_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/henrytk/calendar-resource/calendar"
	"github.com/henrytk/calendar-resource/calendar/fakes"
	"github.com/henrytk/calendar-resource/models"
)

func NewTestCalendar(events []calendar.Event, timeKeeper calendar.TimeKeeper) *calendar.Calendar {
	return &calendar.Calendar{
		events,
		timeKeeper,
	}
}

var _ = Describe("Calendar", func() {

	var (
		fakeTimer   fakes.FakeTimer
		events1     []calendar.Event
		events2     []calendar.Event
		events3     []calendar.Event
		events4     []calendar.Event
		events5     []calendar.Event
		fakeVersion models.Version
	)

	// fakeTimer is a time fixture
	fakeTimer = fakes.NewFakeTimer(calendar.StringToTime("2016-10-02T14:05:00+01:00"))

	// fakeVersion is the version of a resource supplied via standard input
	fakeVersion = models.Version{}

	events1 = []calendar.Event{}

	// events2 contains an event that is currently happening
	events2 = []calendar.Event{
		calendar.Event{
			StartTime: calendar.StringToTime("2016-10-02T14:00:00+01:00"),
			EndTime:   calendar.StringToTime("2016-10-02T15:00:00+01:00"),
		},
	}

	// events3 contains an event that has already finished
	events3 = []calendar.Event{
		calendar.Event{
			StartTime: calendar.StringToTime("2016-10-02T13:00:00+01:00"),
			EndTime:   calendar.StringToTime("2016-10-02T14:00:00+01:00"),
		},
	}

	// events4 contains an event that hasn't started yet
	events4 = []calendar.Event{
		calendar.Event{
			StartTime: calendar.StringToTime("2016-10-02T15:00:00+01:00"),
			EndTime:   calendar.StringToTime("2016-10-02T16:00:00+01:00"),
		},
	}

	// events5 contains two current events
	events5 = []calendar.Event{
		calendar.Event{
			StartTime: calendar.StringToTime("2016-10-02T13:00:00+01:00"),
			EndTime:   calendar.StringToTime("2016-10-02T15:00:00+01:00"),
		},
		calendar.Event{
			StartTime: calendar.StringToTime("2016-10-02T14:00:00+01:00"),
			EndTime:   calendar.StringToTime("2016-10-02T16:00:00+01:00"),
		},
	}

	Context("when there are no events currently happening", func() {
		It("should return an empty array", func() {
			cal := NewTestCalendar(events1, fakeTimer)
			currentVersions := cal.CurrentVersions(fakeVersion)
			Expect(len(currentVersions)).To(Equal(0))
		})
	})

	Context("when there is only one event currently happening", func() {
		It("should return the event", func() {
			cal := NewTestCalendar(events2, fakeTimer)
			currentVersions := cal.CurrentVersions(fakeVersion)
			Expect(len(currentVersions)).To(Equal(1))
			Expect(currentVersions[0].StartTime).To(Equal(calendar.TimeToString(events2[0].StartTime)))
		})
	})

	Context("when the only event has already finished", func() {
		It("should return an empty array", func() {
			cal := NewTestCalendar(events3, fakeTimer)
			currentVersions := cal.CurrentVersions(fakeVersion)
			Expect(len(currentVersions)).To(Equal(0))
		})
	})

	Context("when the only event has not started yet", func() {
		It("should return an empty array", func() {
			cal := NewTestCalendar(events4, fakeTimer)
			currentVersions := cal.CurrentVersions(fakeVersion)
			Expect(len(currentVersions)).To(Equal(0))
		})
	})

	Context("when there are two events currently happening", func() {
		It("should return both events in chronological order", func() {
			cal := NewTestCalendar(events5, fakeTimer)
			currentVersions := cal.CurrentVersions(fakeVersion)
			Expect(len(currentVersions)).To(Equal(2))
			Expect(currentVersions[0].StartTime).To(Equal(calendar.TimeToString(events5[0].StartTime)))
			Expect(currentVersions[1].StartTime).To(Equal(calendar.TimeToString(events5[1].StartTime)))
		})

		Context("when a specific version is requested and is current", func() {
			It("should return only the current version", func() {
				fakeVersion2 := models.Version{
					StartTime: "2016-10-02T14:00:00+01:00",
					EndTime:   "2016-10-02T15:00:00+01:00",
				}
				cal := NewTestCalendar(events5, fakeTimer)
				currentVersions := cal.CurrentVersions(fakeVersion2)
				Expect(len(currentVersions)).To(Equal(1))
				Expect(currentVersions[0]).To(Equal(fakeVersion2))
			})
		})
	})

})
