package entity

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type TicketStatus string

const (
	AvailableTicket TicketStatus = "available"
	PurchasedTicket TicketStatus = "purchased"
	CancelledTicket TicketStatus = "cancelled"
)

type Ticket struct {
	BaseEntity
	EventID      uuid.UUID    `json:"event_id" binding:"required"`
	UserID       uuid.UUID    `json:"user_id"`
	PurchaseDate time.Time    `json:"purchase_date"`
	Status       TicketStatus `json:"status" gorm:"type:ENUM('available', 'purchased', 'cancelled');default:'available'"`
	BookingCode  string       `json:"booking_code" gorm:"unique"`
	Price        float64      `json:"price"`
	Event        Event        `json:"event" gorm:"foreignKey:EventID"`
	User         User         `json:"-" gorm:"foreignKey:UserID"`
}

func (t *Ticket) CanBeCancelled() bool {
	return t.Status == PurchasedTicket && time.Now().Before(t.Event.StartDate)
}
