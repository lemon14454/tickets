package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	db "ticket/backend/db/sqlc"
	"ticket/backend/token"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type createOrderRequest struct {
	ClaimedTickets []int64 `json:"claimed_tickets" binding:"required"`
}

type orderDetail struct {
	Zone  string `json:"zone"`
	Seat  int32  `json:"seat"`
	Row   int32  `json:"row"`
	Price int32  `json:"price"`
}

type createOrderResponse struct {
	Order  db.Order      `json:"order"`
	Detail []orderDetail `json:"ticekts"`
}

func (server *Server) createOrder(ctx *gin.Context) {
	var req createOrderRequest
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

	// Check if userID match the ticket's claim
	for _, ticketID := range req.ClaimedTickets {
		key := fmt.Sprintf("ticket:%d", ticketID)
		userID, err := server.cache.Get(ctx, key).Result()
		if errors.Is(err, redis.Nil) {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		id, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if id != payload.UserID {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
	}

	// Create Order
	var order db.Order
	details := make([]orderDetail, 0)
	err = server.store.ExecTx(ctx, func(q *db.Queries) error {
		tickets, err := q.GetTicketsForUpdate(ctx, req.ClaimedTickets)
		if err != nil {
			return err
		}

		if len(tickets) != len(req.ClaimedTickets) {
			return fmt.Errorf("Request includes sold tickets: %v", req.ClaimedTickets)
		}

		var totalPrice int32
		purchasedTicketsID := make([]int64, 0)
		for _, ticket := range tickets {
			totalPrice += ticket.Price
			details = append(details, orderDetail{
				Zone:  ticket.Zone,
				Seat:  ticket.Seat,
				Row:   ticket.Row,
				Price: ticket.Price,
			})
			purchasedTicketsID = append(purchasedTicketsID, ticket.ID)
		}

		order, err = q.CreateOrder(ctx, db.CreateOrderParams{
			UserID:     &payload.UserID,
			EventID:    tickets[0].EventID,
			TotalPrice: totalPrice,
		})

		if err != nil {
			return err
		}

		err = q.UpdateTicketsUser(ctx, db.UpdateTicketsUserParams{
			UserID:  &payload.UserID,
			OrderID: &order.ID,
			ID:      req.ClaimedTickets,
		})

		return err
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createOrderResponse{
		Order:  order,
		Detail: details,
	})
}
