package model

type eventZone struct {
	Zone  string `json:"zone" binding:"required"`
	Rows  int32  `json:"rows" binding:"required"`
	Seats int32  `json:"seats" binding:"required"`
	Price int32  `json:"price" binding:"required"`
}

type Message struct {
	EventID int64 `json:"event_id" binding:"required,min=1"`
}
