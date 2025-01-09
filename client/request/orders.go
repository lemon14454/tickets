package request

import (
	"net/http"
	"ticket/client/model"
)

type createOrderRequest struct {
	ClaimedTickets []int64 `json:"claimed_tickets"`
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

func (client *Client) CreateOrder(ticketsID []int64) (*orderResponse, error) {
	res, err := MakeRequest[orderResponse](client, http.MethodPost, "order", createOrderRequest{
		ClaimedTickets: ticketsID,
	}, nil)

	return res, err
}
