// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
)

type Querier interface {
	CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error)
	CreateEventZone(ctx context.Context, arg CreateEventZoneParams) (EventZone, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetEventByID(ctx context.Context, id int64) (GetEventByIDRow, error)
	GetEventZones(ctx context.Context, eventID int64) ([]EventZone, error)
	GetRowTickets(ctx context.Context, arg GetRowTicketsParams) ([]Ticket, error)
	GetUser(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
}

var _ Querier = (*Queries)(nil)
