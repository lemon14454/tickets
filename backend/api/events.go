package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	db "ticket/backend/db/sqlc"
	"ticket/backend/token"

	"github.com/gin-gonic/gin"
)

var queueName = "tickets"

type createEventRequest struct {
	Name      string      `json:"name" binding:"required"`
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
			HostID: user.ID,
			Name:   req.Name,
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

	// TODO: Considering about MQ Publish Failed
	msg, err := serialize(Message{
		EventID: event.ID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.mq.Publish(queueName, msg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, event)
}
