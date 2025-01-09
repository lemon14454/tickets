package model

import "time"

type EventStatus string

const (
	EventStatusProcessing EventStatus = "processing"
	EventStatusCreated    EventStatus = "created"
	EventStatusAvailable  EventStatus = "available"
	EventStatusDone       EventStatus = "done"
	EventStatusFailure    EventStatus = "failure"
)

type Event struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	HostID    int64       `json:"host_id"`
	StartAt   time.Time   `json:"start_at"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Status    EventStatus `json:"status"`
}

type EventZone struct {
	Zone  string `json:"zone"`
	Rows  int32  `json:"rows"`
	Seats int32  `json:"seats"`
	Price int32  `json:"price"`
}

type EventZoneDetail struct {
	ID      int64  `json:"id"`
	Zone    string `json:"zone"`
	EventID int64  `json:"event_id"`
	Rows    int32  `json:"rows"`
	Seats   int32  `json:"seats"`
	Price   int32  `json:"price"`
}
