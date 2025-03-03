package request

import (
	"net/http"
	"ticket/client/model"
)

type createOrderRequest struct {
	ClaimedTickets []int64 `json:"claimed_tickets"`
	EventID        int64   `json:"event_id"`
}

type orderDetail struct {
	Zone  string `json:"zone"`
	Seat  int32  `json:"seat"`
	Row   int32  `json:"row"`
	Price int32  `json:"price"`
}

type orderResponse struct {
	Order  model.Order   `json:"order"`
	Detail []orderDetail `json:"ticekts"`
}

func (client *Client) CreateOrder(ticketsID []int64, eventID int64, limit bool) (*orderResponse, error) {
	url := "order"
	if limit {
		url = "limit/order"
	}
	res, err := MakeRequest[orderResponse](client, http.MethodPost, url, createOrderRequest{
		ClaimedTickets: ticketsID,
		EventID:        eventID,
	}, nil)

	return res, err
}
