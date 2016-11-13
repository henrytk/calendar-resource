package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/henrytk/calendar-resource/client"
	"github.com/henrytk/calendar-resource/errors"
	"github.com/henrytk/calendar-resource/models"
	googleCalendarAPI "google.golang.org/api/calendar/v3"
)

func main() {
	if len(os.Args) < 2 {
		errors.Fatal("command line input", fmt.Errorf("Must pass path to build sources"))
	}

	var outRequest models.OutRequest
	inputRequest(&outRequest)

	calendarClient := client.NewCalendarClient(outRequest.Source, googleCalendarAPI.CalendarScope)
	outResponse := calendarClient.AddEvent(&outRequest, os.Args[1])
	outputResponse(&outResponse)
}

func inputRequest(request *models.OutRequest) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		errors.Fatal("reading request from standard input", err)
	}
}

func outputResponse(response *models.OutResponse) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		errors.Fatal("writing response to stdout", err)
	}
}
