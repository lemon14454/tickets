package request

import (
	"net/http"
	"ticket/client/model"
	"time"
)

type GetAllEventsResponse struct {
	ID      int64             `json:"id"`
	Name    string            `json:"name"`
	Status  model.EventStatus `json:"status"`
	StartAt time.Time         `json:"start_at"`
}

func (client *Client) GetAllEvents() ([]GetAllEventsResponse, error) {
	events, err := MakeRequest[[]GetAllEventsResponse](client, http.MethodGet, "event", nil, nil)
	return *events, err
}
