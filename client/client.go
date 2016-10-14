package client

import (
	"fmt"

	"github.com/HenryTK/calendar-resource/calendar"
	"github.com/HenryTK/calendar-resource/errors"
	"github.com/HenryTK/calendar-resource/models"
)

// CalendarClient is an interface that must be satisfied in order to
// implement other calendar providers.
type CalendarClient interface {

	// Events uses the calendar provider's API to return a list of events
	Events() []calendar.Event
}

func NewCalendarClient(source models.Source) CalendarClient {
	var client CalendarClient
	switch source.Provider {
	case "google":
		client = NewGoogleCalendarClient(source)
	default:
		errors.Fatal("Provider error: ", fmt.Errorf("Provider '%v' is not supported", source.Provider))
	}
	return client
}
