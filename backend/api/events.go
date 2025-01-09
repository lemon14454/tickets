package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	db "ticket/backend/db/sqlc"
	"ticket/backend/token"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	TicketQueue      = "ticket_queue"
	TicketExchange   = "ticket_exchange"
	TicketRoutingKey = "ticket_routing_key"
)

type createEventRequest struct {
	Name      string      `json:"name" binding:"required"`
	StartAt   string      `json:"start_at" binding:"required"`
	EventZone []eventZone `json:"event_zone" binding:"dive"`
}

type eventZone struct {
	Zone  string `json:"zone" binding:"required"`
	Rows  int32  `json:"rows" binding:"required"`
	Seats int32  `json:"seats" binding:"required"`
	Price int32  `json:"price" binding:"required"`
}

type Message struct {
	EventID int64 `json:"event_id" binding:"required,min=1"`
}

func serialize(msg Message) ([]byte, error) {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	err := encoder.Encode(msg)
	return b.Bytes(), err
}

func (server *Server) createEvent(ctx *gin.Context) {
	var req createEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	startTime, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.store.GetUserByID(ctx, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !user.Host {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if payload.UserID != user.ID {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var event db.Event
	server.store.ExecTx(ctx, func(q *db.Queries) error {
		event, err = server.store.CreateEvent(ctx, db.CreateEventParams{
			HostID:  user.ID,
			StartAt: startTime,
			Name:    req.Name,
		})

		if err != nil {
			return err
		}

		for _, ez := range req.EventZone {
			_, err := server.store.CreateEventZone(ctx, db.CreateEventZoneParams{
				Zone:    ez.Zone,
				Rows:    ez.Rows,
				Seats:   ez.Seats,
				EventID: event.ID,
				Price:   ez.Price,
			})

			if err != nil {
				return err
			}
		}

		return nil
	})

	msg, err := serialize(Message{
		EventID: event.ID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.mq.Publish(TicketExchange, TicketRoutingKey, msg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, event)
}

func (server *Server) listEvent(ctx *gin.Context) {
	events, err := server.store.GetAllEvent(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, events)
}

type listEventZoneRequest struct {
	EventID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) listEventZone(ctx *gin.Context) {
	var req listEventZoneRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	zones, err := server.store.GetEventZone(ctx, req.EventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, zones)
}
