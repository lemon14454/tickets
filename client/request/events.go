package request

import (
	"fmt"
	"net/http"
	"ticket/client/model"
	"time"
)

type getAllEventsResponse struct {
	ID      int64             `json:"id"`
	Name    string            `json:"name"`
	Status  model.EventStatus `json:"status"`
	StartAt time.Time         `json:"start_at"`
}

type createEventRequest struct {
	Name      string            `json:"name"`
	StartAt   string            `json:"start_at"`
	EventZone []model.EventZone `json:"event_zone"`
}

func (client *Client) GetAllEvents() ([]getAllEventsResponse, error) {
	events, err := MakeRequest[[]getAllEventsResponse](client, http.MethodGet, "event", nil, nil)
	return *events, err
}

func (client *Client) CreateEvent(name string, eventZone []model.EventZone) (*model.Event, error) {

	event, err := MakeRequest[model.Event](client, http.MethodPost, "event", createEventRequest{
		Name:      name,
		StartAt:   "2024-07-23T15:04:05Z", // I don't care, hard coded
		EventZone: eventZone,
	}, nil)
	return event, err
}

type listEventZoneRequest struct {
	EventID int64 `uri:"id" binding:"required,min=1"`
}

func (client *Client) GetEventZoneByID(eventID int64) (*[]model.EventZoneDetail, error) {
	url := fmt.Sprintf("event/%d", eventID)
	eventZones, err := MakeRequest[[]model.EventZoneDetail](client, http.MethodGet, url, nil, nil)
	return eventZones, err
}
