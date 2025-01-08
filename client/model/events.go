package model

type EventStatus string

const (
	EventStatusProcessing EventStatus = "processing"
	EventStatusCreated    EventStatus = "created"
	EventStatusAvailable  EventStatus = "available"
	EventStatusDone       EventStatus = "done"
	EventStatusFailure    EventStatus = "failure"
)
