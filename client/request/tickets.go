package request

import (
	"net/http"
	"time"
)

type claimTicketRequest struct {
	EventID  int64 `json:"event_id"`
	ZoneID   int64 `json:"zone_id"`
	Row      int32 `json:"row"` // 10 Rows each zone for now
	Quantity int   `json:"quantity"`
}

type claimTicketResponse struct {
	ClaimedTickets []int64   `json:"claimed_tickets"`
	ClaimExpiresAt time.Time `json:"claim_expires_at"`
}

func (client *Client) ClaimTicket(eventID, zoneID int64, row int32, quantity int, limit bool) (*claimTicketResponse, error) {
	url := "ticket"
	if limit {
		url = "limit/ticket"
	}

	res, err := MakeRequest[claimTicketResponse](client, http.MethodPost, url, claimTicketRequest{
		EventID:  eventID,
		ZoneID:   zoneID,
		Row:      row,
		Quantity: quantity,
	}, nil)

	return res, err
}
