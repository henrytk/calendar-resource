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
		errors.Fatal("command line input", fmt.Errorf("Must pass target directory as input"))
	}

	targetDirectory := os.Args[1]
	if err := os.MkdirAll(targetDirectory, 0755); err != nil {
		errors.Fatal("creating target directory", err)
	}

	var inRequest models.InRequest
	inputRequest(&inRequest)
	calendarClient := client.NewCalendarClient(inRequest.Source)
	inResponse, file, err := calendarClient.GetEvent(&inRequest, targetDirectory)
	if err != nil {
		errors.Fatal("getting event details for input file", err)
	}
	defer file.Close()

	outputResponse(&inResponse)
}

func inputRequest(request *models.InRequest) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		errors.Fatal("reading request from standard input", err)
	}
}

func outputResponse(response *models.InResponse) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		errors.Fatal("writing response to standard output", err)
	}
}
