package models

import (
	"encoding/json"
)

type Source struct {
	Provider    string          `json:"provider"`
	CalendarId  string          `json:"calendar_id"`
	EventName   string          `json:"event_name"`
	Credentials json.RawMessage `json:"credentials"`
}

type Version struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}
