package model

import "time"

type Order struct {
	ID         int64     `json:"id"`
	UserID     *int64    `json:"user_id"`
	EventID    int64     `json:"event_id"`
	CreatedAt  time.Time `json:"created_at"`
	TotalPrice int32     `json:"total_price"`
}
