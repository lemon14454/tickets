package model

type Ticket struct {
	EventID int64 `json:"event_id"`
	ZoneID  int64 `json:"zone_id"`
	Row     int32 `json:"row"`
	Seat    int32 `json:"seat"`
}
