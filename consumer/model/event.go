package model

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type EventStatus string

const (
	EventStatusProcessing EventStatus = "processing"
	EventStatusCreated    EventStatus = "created"
	EventStatusAvailable  EventStatus = "available"
	EventStatusDone       EventStatus = "done"
	EventStatusFailure    EventStatus = "failure"
)

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

func CreateEvent(db *gorm.DB, eventID int64) error {
	var evt Event
	result := db.First(&evt, eventID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("Unable to find Event: %d \n", eventID)
		// normally it won't happend
	}
	if result.Error != nil {
		return result.Error
	}

	var eventZone []EventZone
	result = db.Where("event_id = ?", eventID).Find(&eventZone)
	if result.Error != nil {
		log.Printf("Unable to find EventZone: %v \n", result.Error)
		return result.Error
	}

	tickets := make([]Ticket, 0)
	for _, zone := range eventZone {
		for row := range zone.Rows {
			for seat := range zone.Seats {
				tickets = append(tickets, Ticket{
					EventID: evt.ID,
					ZoneID:  zone.ID,
					Row:     row + 1,
					Seat:    seat + 1,
				})
			}
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Create(&tickets).Error; err != nil {
			return err
		}

		if err := db.Model(&evt).Update("status", EventStatusCreated).Update("updated_at", time.Now()).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func MarkEventFailure(db *gorm.DB, eventID int64) error {

	err := db.Transaction(func(tx *gorm.DB) error {
		var evt Event
		result := db.First(&evt, eventID)
		if err := result.Error; err != nil {
			return err
		}

		if err := db.Model(&evt).Update("status", EventStatusFailure).Update("updated_at", time.Now()).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
