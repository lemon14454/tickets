// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type EventStatus string

const (
	EventStatusProcessing EventStatus = "processing"
	EventStatusCreated    EventStatus = "created"
	EventStatusAvailable  EventStatus = "available"
	EventStatusDone       EventStatus = "done"
	EventStatusFailure    EventStatus = "failure"
)

func (e *EventStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = EventStatus(s)
	case string:
		*e = EventStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for EventStatus: %T", src)
	}
	return nil
}

type NullEventStatus struct {
	EventStatus EventStatus `json:"event_status"`
	Valid       bool        `json:"valid"` // Valid is true if EventStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullEventStatus) Scan(value interface{}) error {
	if value == nil {
		ns.EventStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.EventStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullEventStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.EventStatus), nil
}

type Event struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	HostID    int64       `json:"host_id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Status    EventStatus `json:"status"`
}

type EventZone struct {
	ID      int64  `json:"id"`
	Zone    string `json:"zone"`
	EventID int64  `json:"event_id"`
	Rows    int32  `json:"rows"`
	Seats   int32  `json:"seats"`
	Price   int32  `json:"price"`
}

type Order struct {
	ID         int64       `json:"id"`
	UserID     pgtype.Int8 `json:"user_id"`
	CreatedAt  time.Time   `json:"created_at"`
	TotalPrice int32       `json:"total_price"`
}

type Ticket struct {
	ID        int64       `json:"id"`
	UserID    pgtype.Int8 `json:"user_id"`
	OrderID   pgtype.Int8 `json:"order_id"`
	EventID   int64       `json:"event_id"`
	ZoneID    int64       `json:"zone_id"`
	Row       int32       `json:"row"`
	Seat      int32       `json:"seat"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	Host           bool      `json:"host"`
	CreatedAt      time.Time `json:"created_at"`
}
