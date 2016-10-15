package client

import (
	"fmt"
	"os"

	"github.com/henrytk/calendar-resource/calendar"
	"github.com/henrytk/calendar-resource/errors"
	"github.com/henrytk/calendar-resource/models"
)

// CalendarClient is an interface that must be satisfied in order to
// implement other calendar providers.
type CalendarClient interface {

	// ListEvents uses the calendar provider's API to return a list of events
	ListEvents() []calendar.Event

	// GetEvent takes the `in` request data and a directory path under which
	// a file will be created. It uses the calendar provider's API to get
	// the event details necessary to provide a response on standard output
	// and populate a file. The file will then be placed in the Concourse
	// task's file system.
	GetEvent(*models.InRequest, string) (models.InResponse, *os.File, error)
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
