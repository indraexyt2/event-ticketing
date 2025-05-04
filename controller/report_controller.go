package controller

import (
	"event-ticketing/service"
	"event-ticketing/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportController interface {
	GetSummaryReport(c *gin.Context)
	GetEventReport(c *gin.Context)
}

type reportController struct {
	reportService service.ReportService
}

func NewReportController(reportService service.ReportService) ReportController {
	return &reportController{
		reportService: reportService,
	}
}

// GetSummaryReport godoc
// @Summary Get system summary report
// @Description Get a summary report of the entire ticketing system (admin only)
// @Tags reports
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /reports/summary [get]
func (ctrl *reportController) GetSummaryReport(c *gin.Context) {
	var log = utils.Log

	report, err := ctrl.reportService.GenerateSummaryReport()
	if err != nil {
		log.Errorf("Failed to generate report: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to generate report", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Report generated successfully", report)
}

// GetEventReport godoc
// @Summary Get report for a specific event
// @Description Get a detailed report for a specific event (admin only)
// @Tags reports
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Event ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /reports/events/{id} [get]
func (ctrl *reportController) GetEventReport(c *gin.Context) {
	var log = utils.Log

	eventID := c.Param("id")
	report, err := ctrl.reportService.GenerateEventReport(eventID)
	if err != nil {
		log.Errorf("Failed to generate report: %v", err)
		utils.NotFoundResponse(c, "Event not found or failed to generate report")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Report generated successfully", report)
}
