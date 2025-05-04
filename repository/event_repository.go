package repository

import (
	"errors"
	"event-ticketing/entity"
	"event-ticketing/utils"
	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *entity.Event) error
	FindByID(id string) (*entity.Event, error)
	FindByIDForUpdate(id string) (*entity.Event, error)
	FindAll(params utils.PaginationParams, keyword, status, startDate, endDate string) ([]entity.Event, int64, error)
	FindByName(name string) (*entity.Event, error)
	Update(event *entity.Event) error
	Delete(id string) error
	CountTicketsSold(eventID string) (int, error)
	WithTx(tx *gorm.DB) EventRepository
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db}
}

func (r *eventRepository) Create(event *entity.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) FindByID(id string) (*entity.Event, error) {
	var event entity.Event
	err := r.db.Where("id = ?", id).First(&event).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, err
	}

	soldTickets, err := r.CountTicketsSold(id)
	if err != nil {
		return nil, err
	}
	event.TicketsSold = soldTickets

	return &event, nil
}

func (r *eventRepository) FindByIDForUpdate(id string) (*entity.Event, error) {
	var event entity.Event
	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&event).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, err
	}

	soldTickets, err := r.CountTicketsSold(id)
	if err != nil {
		return nil, err
	}
	event.TicketsSold = soldTickets

	return &event, nil
}

func (r *eventRepository) WithTx(tx *gorm.DB) EventRepository {
	return &eventRepository{db: tx}
}

func (r *eventRepository) FindAll(params utils.PaginationParams, keyword, status, startDate, endDate string) ([]entity.Event, int64, error) {
	var events []entity.Event
	var count int64

	query := r.db.Model(&entity.Event{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ? OR location LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if startDate != "" {
		query = query.Where("start_date >= ?", startDate)
	}

	if endDate != "" {
		query = query.Where("end_date <= ?", endDate)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(params.GetOffset()).Limit(params.GetLimit()).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	for i, event := range events {
		soldTickets, err := r.CountTicketsSold(event.ID.String())
		if err != nil {
			return nil, 0, err
		}
		events[i].TicketsSold = soldTickets
	}

	return events, count, nil
}

func (r *eventRepository) FindByName(name string) (*entity.Event, error) {
	var event entity.Event
	err := r.db.Where("name = ?", name).First(&event).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) Update(event *entity.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(id string) error {
	var ticketCount int64
	err := r.db.Model(&entity.Ticket{}).Where("event_id = ? AND status = ?", id, entity.PurchasedTicket).Count(&ticketCount).Error
	if err != nil {
		return err
	}

	if ticketCount > 0 {
		return errors.New("cannot delete event with sold tickets")
	}

	return r.db.Where("id = ?", id).Delete(&entity.Event{}).Error
}

func (r *eventRepository) CountTicketsSold(eventID string) (int, error) {
	var count int64
	err := r.db.Model(&entity.Ticket{}).Where("event_id = ? AND status = ?", eventID, entity.PurchasedTicket).Count(&count).Error
	return int(count), err
}
