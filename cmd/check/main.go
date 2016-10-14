package main

import (
	"encoding/json"
	"os"

	"github.com/HenryTK/calendar-resource/calendar"
	"github.com/HenryTK/calendar-resource/client"
	"github.com/HenryTK/calendar-resource/errors"
	"github.com/HenryTK/calendar-resource/models"
)

func main() {
	var checkRequest models.CheckRequest
	inputRequest(&checkRequest)
	calendarClient := client.NewCalendarClient(checkRequest.Source)
	calendar := calendar.NewCalendar()
	calendar.Events = calendarClient.Events()
	currentVersions := calendar.CurrentVersions(checkRequest.Version)
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
