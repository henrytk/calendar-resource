package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/henrytk/calendar-resource/client"
	"github.com/henrytk/calendar-resource/errors"
	"github.com/henrytk/calendar-resource/models"
)

func main() {
	if len(os.Args) < 2 {
		errors.Fatal("command line input", fmt.Errorf("Must pass path to build sources"))
	}

	var outRequest models.OutRequest
	inputRequest(&outRequest)

	calendarClient := client.NewCalendarClient(outRequest.Source)
	calendarClient.AddEvent(&outRequest, os.Args[1])
	outputResponse(&outRequest)
}

func inputRequest(request *models.OutRequest) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		errors.Fatal("reading request from standard input", err)
	}
}

// outputResponse emits the requested calendar event to standard output.
// The Concourse interface demands the generated version of the resource
// is emitted, but calendar event versions are identified by a calendar ID,
// which is a value you cannot do much with.
func outputResponse(request *models.OutRequest) {
	os.Stdout.Write(request.Params)
}
