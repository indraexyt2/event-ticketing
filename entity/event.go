package entity

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type EventStatus string

const (
	ActiveEvent    EventStatus = "active"
	OngoingEvent   EventStatus = "ongoing"
	CompletedEvent EventStatus = "completed"
)

type Event struct {
	BaseEntity
	Name        string      `gorm:"unique" json:"name"`
	Description string      `json:"description"`
	StartDate   time.Time   `json:"start_date"`
	EndDate     time.Time   `json:"end_date"`
	Capacity    int         `json:"capacity"`
	Price       float64     `json:"price"`
	Status      EventStatus `json:"status" gorm:"type:ENUM('active', 'ongoing', 'completed');default:'active'"`
	Location    string      `json:"location"`
	Tickets     []Ticket    `json:"-" gorm:"foreignKey:EventID"`
	CreatedBy   uuid.UUID   `gorm:"type:char(36)" json:"created_by"`
	TicketsSold int         `json:"tickets_sold" gorm:"-"`
	User        User        `json:"-" gorm:"foreignKey:CreatedBy;references:ID"`
}

func (e *Event) CanBeModified() bool {
	return e.Status != CompletedEvent && e.Status != OngoingEvent
}

func (e *Event) HasAvailableTickets(count int) bool {
	return e.Capacity >= e.TicketsSold+count
}
