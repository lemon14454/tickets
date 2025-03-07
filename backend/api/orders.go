package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	db "ticket/backend/db/sqlc"
	"ticket/backend/token"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type createOrderRequest struct {
	ClaimedTickets []int64 `json:"claimed_tickets" binding:"required"`
	EventID        int64   `json:"event_id" binding:"required,min=1"`
}

type orderDetail struct {
	Zone  string `json:"zone"`
	Seat  int32  `json:"seat"`
	Row   int32  `json:"row"`
	Price int32  `json:"price"`
}

type orderResponse struct {
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

	ticketsToDeleteFromCache := make([]string, len(req.ClaimedTickets))
	// Check if userID match the ticket's claim
	for _, ticketID := range req.ClaimedTickets {
		key := fmt.Sprintf("ticket:%d-%d", req.EventID, ticketID)
		ticketsToDeleteFromCache = append(ticketsToDeleteFromCache, key)
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
			ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("Ticket belongs to userID: %d, not %d \n", id, payload.UserID)))
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

	for _, soldTicket := range ticketsToDeleteFromCache {
		server.cache.Set(ctx, soldTicket, "SOLD", 60*time.Second)
	}

	ctx.JSON(http.StatusOK, orderResponse{
		Order:  order,
		Detail: details,
	})
}

type listOrderResponse struct {
	Orders []db.GetUserOrdersRow `json:"order"`
}

func (server *Server) listOrder(ctx *gin.Context) {

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	orders, err := server.store.GetUserOrders(ctx, &payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

type getOrderDetail struct {
	OrderID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) orderDetail(ctx *gin.Context) {

	var req getOrderDetail
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	detail, err := server.store.GetOrderDetail(ctx, &req.OrderID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, detail)
}
