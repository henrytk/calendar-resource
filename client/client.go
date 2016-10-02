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

	// Events takes a calendar and must populate calendar.Events using
	// the calendar provider's API.
	Events(*calendar.Calendar) error
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
