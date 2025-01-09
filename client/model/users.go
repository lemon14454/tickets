package model

import "time"

type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	Host           bool      `json:"host"`
	CreatedAt      time.Time `json:"created_at"`
}
