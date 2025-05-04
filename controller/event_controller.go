package controller

import (
	"event-ticketing/dto"
	"event-ticketing/service"
	"event-ticketing/utils"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"net/http"
)

type EventController interface {
	CreateEvent(c *gin.Context)
	GetEventByID(c *gin.Context)
	GetAllEvents(c *gin.Context)
	UpdateEvent(c *gin.Context)
	DeleteEvent(c *gin.Context)
}

type eventController struct {
	eventService service.EventService
}

func NewEventController(eventService service.EventService) EventController {
	return &eventController{
		eventService: eventService,
	}
}

// CreateEvent godoc
// @Summary Create a new event
// @Description Create a new event with the provided information
// @Tags events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param event body dto.CreateEventReqDto true "Event creation info"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /events [post]
func (ctrl *eventController) CreateEvent(c *gin.Context) {
	var log = utils.Log

	var eventRequest dto.CreateEventReqDto
	if err := c.ShouldBindJSON(&eventRequest); err != nil {
		log.Errorf("Failed to bind JSON: %v", err)
		utils.BadRequestResponse(c, "Invalid request body", err.Error())
		return
	}

	if err := utils.ValidateStruct(eventRequest); err != nil {
		log.Errorf("Validation error: %v", err)
		utils.BadRequestResponse(c, "Validation error", err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.Error("User not authenticated")
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	idUuid, err := uuid.FromString(userID.(string))
	if err != nil {
		log.Errorf("Failed to generate UUID: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to generate UUID", err.Error())
		return
	}

	eventRequest.CreatedBy = idUuid

	event := eventRequest.ToEntity()

	if err := ctrl.eventService.CreateEvent(event); err != nil {
		log.Errorf("Event creation failed: %v", err)
		utils.ConflictResponse(c, "Event creation failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Event created successfully", event)
}

// GetEventByID godoc
// @Summary Get an event by ID
// @Description Get an event by its ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /events/{id} [get]
func (ctrl *eventController) GetEventByID(c *gin.Context) {
	var log = utils.Log

	id := c.Param("id")
	event, err := ctrl.eventService.GetEventByID(id)
	if err != nil {
		log.Errorf("Event not found: %v", err)
		utils.NotFoundResponse(c, "Event not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Event retrieved successfully", event)
}

// GetAllEvents godoc
// @Summary Get all events with pagination and filtering
// @Description Get all events with pagination and optional filtering
// @Tags events
// @Accept json
// @Produce json
// @Param page query string false "Page number (default: 1)"
// @Param limit query string false "Results per page (default: 10)"
// @Param keyword query string false "Search keyword"
// @Param status query string false "Event status (active/ongoing/completed)"
// @Param start_date query string false "Start date filter (YYYY-MM-DD)"
// @Param end_date query string false "End date filter (YYYY-MM-DD)"
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /events [get]
func (ctrl *eventController) GetAllEvents(c *gin.Context) {
	var log = utils.Log

	params := utils.NewPaginationParams(c.DefaultQuery("page", "1"), c.DefaultQuery("limit", "10"))

	keyword := c.Query("keyword")
	status := c.Query("status")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	events, totalItems, err := ctrl.eventService.GetAllEvents(params, keyword, status, startDate, endDate)
	if err != nil {
		log.Errorf("Failed to retrieve events: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to retrieve events", err.Error())
		return
	}

	utils.PaginatedResponse(c, http.StatusOK, "Events retrieved successfully", events, totalItems, params.Page, params.Limit)
}

// UpdateEvent godoc
// @Summary Update an event
// @Description Update an event with the provided information
// @Tags events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Event ID"
// @Param event body dto.UpdateEventReqDto true "Updated event info"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /events/{id} [put]
func (ctrl *eventController) UpdateEvent(c *gin.Context) {
	var log = utils.Log

	id := c.Param("id")
	if id == "" {
		log.Error("Event ID not found in request")
		utils.BadRequestResponse(c, "Event not found", "Event ID not found in request")
		return
	}

	var updatedEvent dto.UpdateEventReqDto
	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		log.Errorf("Failed to bind JSON: %v", err)
		utils.BadRequestResponse(c, "Invalid request body", err.Error())
		return
	}

	IDUuid, err := uuid.FromString(id)
	if err != nil {
		log.Errorf("Invalid event ID: %v", err)
		utils.BadRequestResponse(c, "Invalid event ID", err.Error())
		return
	}

	updatedEvent.ID = IDUuid

	event := updatedEvent.ToEntity()
	event, err = ctrl.eventService.UpdateEvent(event)
	if err != nil {
		log.Errorf("Failed to update event: %v", err)
		utils.BadRequestResponse(c, "Failed to update event", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Event updated successfully", event)
}

// DeleteEvent godoc
// @Summary Delete an event
// @Description Delete an event by its ID
// @Tags events
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Event ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /events/{id} [delete]
func (ctrl *eventController) DeleteEvent(c *gin.Context) {
	var log = utils.Log

	id := c.Param("id")

	if err := ctrl.eventService.DeleteEvent(id); err != nil {
		log.Errorf("Failed to delete event: %v", err)
		utils.BadRequestResponse(c, "Failed to delete event", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Event deleted successfully", nil)
}
