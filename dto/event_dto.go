package dto

import (
	"event-ticketing/entity"
	"github.com/gofrs/uuid/v5"
	"time"
)

type CreateEventReqDto struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Capacity    int       `json:"capacity" binding:"required,min=1"`
	Price       float64   `json:"price" binding:"required,min=0"`
	Status      string    `json:"status" binding:"required,oneof=active ongoing completed"`
	Location    string    `json:"location" binding:"required"`
	CreatedBy   uuid.UUID `json:"-"`
}

func (e *CreateEventReqDto) ToEntity() *entity.Event {
	return &entity.Event{
		Name:        e.Name,
		Description: e.Description,
		StartDate:   e.StartDate,
		EndDate:     e.EndDate,
		Capacity:    e.Capacity,
		Price:       e.Price,
		Status:      entity.EventStatus(e.Status),
		Location:    e.Location,
		CreatedBy:   e.CreatedBy,
	}
}

type UpdateEventReqDto struct {
	ID          uuid.UUID `json:"-"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Capacity    int       `json:"capacity" binding:"required,min=1"`
	Price       float64   `json:"price" binding:"required,min=0"`
	Status      string    `json:"status" binding:"required,oneof=active ongoing completed"`
	Location    string    `json:"location" binding:"required"`
}

func (e *UpdateEventReqDto) ToEntity() *entity.Event {
	return &entity.Event{
		BaseEntity:  entity.BaseEntity{ID: e.ID},
		Name:        e.Name,
		Description: e.Description,
		StartDate:   e.StartDate,
		EndDate:     e.EndDate,
		Capacity:    e.Capacity,
		Price:       e.Price,
		Status:      entity.EventStatus(e.Status),
		Location:    e.Location,
	}
}
