package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/henrytk/calendar-resource/calendar"
	"github.com/henrytk/calendar-resource/errors"
	"github.com/henrytk/calendar-resource/models"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	googleCalendarAPI "google.golang.org/api/calendar/v3"
)

type GoogleCalendarClient struct {
	Source     models.Source
	HTTPClient *http.Client
}

func NewGoogleCalendarClient(source models.Source) CalendarClient {
	ctx := context.Background()
	config, err := google.JWTConfigFromJSON(source.Credentials, googleCalendarAPI.CalendarReadonlyScope)
	if err != nil {
		errors.Fatal("JWTConfigFromJSON error: ", err)
	}
	client := config.Client(ctx)
	return &GoogleCalendarClient{
		Source:     source,
		HTTPClient: client,
	}
}

func (gcc *GoogleCalendarClient) getService() *googleCalendarAPI.Service {
	service, err := googleCalendarAPI.New(gcc.HTTPClient)
	if err != nil {
		errors.Fatal("Google calendar API error: ", err)
	}
	return service
}

func (gcc *GoogleCalendarClient) ListEvents() []calendar.Event {
	var calendarEvents []calendar.Event
	service := gcc.getService()

	t := time.Now().Format(time.RFC3339)
	events, err := service.Events.List(gcc.Source.CalendarId).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).OrderBy("startTime").Do()
	if err != nil {
		errors.Fatal("getting events using calendar client", err)
	}
	if len(events.Items) > 0 {
		for _, item := range events.Items {
			if item.Summary == gcc.Source.EventName {
				// If the DateTime is an empty string the Event is an all-day Event.
				// So only Date is available.
				var startTime time.Time
				if item.Start.DateTime != "" {
					startTime = gcc.parseTime(item.Start.DateTime)
				} else {
					startTime = gcc.parseDate(item.Start.Date, events.TimeZone)
				}
				var endTime time.Time
				if item.End.DateTime != "" {
					endTime = gcc.parseTime(item.End.DateTime)
				} else {
					endTime = gcc.parseDate(item.End.Date, events.TimeZone)
				}
				calendarEvents = append(calendarEvents, calendar.Event{
					StartTime: startTime,
					EndTime:   endTime,
				})
			}
		}
	}
	return calendarEvents
}

func (gcc *GoogleCalendarClient) GetEvent(inRequest *models.InRequest, targetDirectory string) (models.InResponse, *os.File, error) {
	startTime := inRequest.Version.StartTime
	endTime := inRequest.Version.EndTime
	if startTime == "" || endTime == "" {
		errors.Fatal("fetching resource version", fmt.Errorf("resource version start time or end time not specified"))
	}
	service := gcc.getService()
	events, err := service.Events.List(gcc.Source.CalendarId).ShowDeleted(false).
		SingleEvents(true).TimeMin(endTime).OrderBy("startTime").Do()
	if err != nil {
		errors.Fatal("getting events using calendar client", err)
	}
	if len(events.Items) == 0 {
		errors.Fatal("fetching resource version", fmt.Errorf("Event start time: %v, end time: %v", startTime, endTime))
	}
	event := events.Items[0]
	if event.Summary != gcc.Source.EventName {
		errors.Fatal("fetching resource version", fmt.Errorf("Event not found"))
	}
	keyValuePairs := []models.KeyValuePair{
		models.KeyValuePair{
			Name:  "created",
			Value: event.Created,
		},
		models.KeyValuePair{
			Name:  "description",
			Value: event.Description,
		},
		models.KeyValuePair{
			Name:  "hangoutLink",
			Value: event.HangoutLink,
		},
		models.KeyValuePair{
			Name:  "htmlLink",
			Value: event.HtmlLink,
		},
		models.KeyValuePair{
			Name:  "iCalUid",
			Value: event.ICalUID,
		},
		models.KeyValuePair{
			Name:  "Id",
			Value: event.Id,
		},
		models.KeyValuePair{
			Name:  "summary",
			Value: event.Summary,
		},
	}
	inResponse := models.InResponse{
		Version:  inRequest.Version,
		MetaData: keyValuePairs,
	}

	file, err := os.Create(filepath.Join(targetDirectory, "input"))
	defer file.Close()
	if err != nil {
		errors.Fatal("creating input file", err)
	}
	err = json.NewEncoder(file).Encode(inResponse)
	if err != nil {
		errors.Fatal("reading from input request", err)
	}
	return inResponse, file, nil
}

// AddEventParams holds data passed in via `params` from the
// Conncourse task. StartTime and EndTime must be an RFC3339
// formatted time string. For example, "2016-10-15T08:00:00+01:00"
type AddEventParams struct {
	Description string `json:"description,omitempty"`
	EndTime     string `json:"end_time"`
	StartTime   string `json:"start_time"`
	Summary     string `json:"summary,omitempty"`
	TimeZone    string `json:"time_zone,omitempty"`
}

func (gcc *GoogleCalendarClient) AddEvent(outRequest *models.OutRequest, buildSourcePath string) {
	var addEventParams AddEventParams
	if err := json.Unmarshal(outRequest.Params, &addEventParams); err != nil {
		errors.Fatal("decoding event params", err)
	}
	if addEventParams.StartTime == "" || addEventParams.EndTime == "" {
		errors.Fatal("adding event", fmt.Errorf("You must supply an event start and end time"))
	}
	service := gcc.getService()
	event := googleCalendarAPI.Event{
		// These values are not currently supported by the calendar resource, but
		// are not optional, so we pass empty literals.
		Attachments: []*googleCalendarAPI.EventAttachment{},
		Attendees:   []*googleCalendarAPI.EventAttendee{},
		Reminders:   &googleCalendarAPI.EventReminders{},

		// These values can be set by the calendar resource. Only Start and End
		// are mandatory.
		Description: addEventParams.Description,
		End:         &googleCalendarAPI.EventDateTime{TimeZone: addEventParams.TimeZone, DateTime: addEventParams.EndTime},
		Start:       &googleCalendarAPI.EventDateTime{TimeZone: addEventParams.TimeZone, DateTime: addEventParams.StartTime},
		Summary:     addEventParams.Summary,
	}
	_, err := service.Events.Insert(outRequest.Source.CalendarId, &event).Do()
	if err != nil {
		errors.Fatal("adding event", err)
	}
}

func (gcc *GoogleCalendarClient) parseTime(timeString string) time.Time {
	t, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		errors.Fatal("parsing time: ", err)
	}
	return t
}

func (gcc *GoogleCalendarClient) parseDate(date, location string) time.Time {
	loc, err := time.LoadLocation(location)
	if err != nil {
		errors.Fatal("parsing time: ", err)
	}
	t, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		errors.Fatal("parsing date: ", err)
	}
	return t
}
