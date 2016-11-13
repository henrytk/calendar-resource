package main

import (
	"encoding/json"
	"os"

	"github.com/henrytk/calendar-resource/client"
	"github.com/henrytk/calendar-resource/errors"
	"github.com/henrytk/calendar-resource/models"
	googleCalendarAPI "google.golang.org/api/calendar/v3"
)

func main() {
	var checkRequest models.CheckRequest
	inputRequest(&checkRequest)
	calendarClient := client.NewCalendarClient(checkRequest.Source, googleCalendarAPI.CalendarReadonlyScope)
	currentVersions := calendarClient.ListEvents(checkRequest.Version)
	outputResponse(currentVersions)
}

func inputRequest(request *models.CheckRequest) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		errors.Fatal("Reading request from standard input", err)
	}
}

func outputResponse(versions []models.Version) {
	if err := json.NewEncoder(os.Stdout).Encode(versions); err != nil {
		errors.Fatal("writing response to stdout", err)
	}
}
