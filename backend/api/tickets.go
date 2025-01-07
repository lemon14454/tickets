package api

import (
	"net/http"
	"strconv"
	"ticket/backend/cache"
	db "ticket/backend/db/sqlc"
	"ticket/backend/token"
	"time"

	"github.com/gin-gonic/gin"
)

type claimTicketRequest struct {
	EventID  int64 `json:"event_id" binding:"required,min=1"`
	ZoneID   int64 `json:"zone_id" binding:"required,min=1"`
	Row      int32 `json:"row" binding:"required,min=1,max=10"` // 10 Rows each zone for now
	Quantity int   `json:"quantity" binding:"required,min=1,max=4"`
}

type claimTicketResponse struct {
	ClaimedTickets []int64   `json:"claimed_tickets"`
	ClaimExpiresAt time.Time `json:"claim_expires_at"`
}

var notEnoughTicket = gin.H{"error": "The remaining quantity of this item does not meet the required amount."}

func (server *Server) claimTicket(ctx *gin.Context) {
	var req claimTicketRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	_, err := server.store.GetUserByID(ctx, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	tickets, err := server.store.GetRowTickets(ctx, db.GetRowTicketsParams{
		EventID: req.EventID,
		ZoneID:  req.ZoneID,
		Row:     req.Row,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// This Filter could be done by SQL
	availableTickets := make([]interface{}, 0)
	for _, ticket := range tickets {
		if !ticket.UserID.Valid && !ticket.OrderID.Valid {
			availableTickets = append(availableTickets, ticket.ID)
		}
	}

	if len(availableTickets) >= req.Quantity {
		keys := []string{strconv.Itoa(req.Quantity), strconv.FormatInt(payload.UserID, 10)}
		// The claim will force expired if the client close the window
		result, err := cache.ClaimCacheTicket.Run(ctx, server.cache, keys, availableTickets...).Int64Slice()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		} else if len(result) == 0 {
			ctx.JSON(http.StatusOK, notEnoughTicket)
			return
		}

		ctx.JSON(http.StatusOK, claimTicketResponse{
			ClaimExpiresAt: time.Now().Add(5 * time.Minute),
			ClaimedTickets: result,
		})

	} else {
		ctx.JSON(http.StatusOK, notEnoughTicket)
		return
	}

}
