// routes/route.go
package routes

import (
	"event-ticketing/config"
	"event-ticketing/controller"
	"event-ticketing/middleware"
	"event-ticketing/repository"
	"event-ticketing/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func SetupRoutes(db *gorm.DB, config config.Config) *gin.Engine {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	ticketRepo := repository.NewTicketRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, config)
	eventService := service.NewEventService(eventRepo, ticketRepo)
	ticketService := service.NewTicketService(db, ticketRepo, eventRepo)
	reportService := service.NewReportService(eventRepo, userRepo, ticketRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	eventController := controller.NewEventController(eventService)
	ticketController := controller.NewTicketController(ticketService)
	reportController := controller.NewReportController(reportService)

	// Create router
	router := gin.Default()

	// Auth routes
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.GET("/profile", middleware.AuthMiddleware(userRepo, config), authController.GetProfile)
	}

	// Event routes
	eventRoutes := router.Group("/api/events")
	{
		eventRoutes.GET("", eventController.GetAllEvents)
		eventRoutes.GET("/:id", eventController.GetEventByID)

		// Protected routes
		eventRoutes.POST("", middleware.AuthMiddleware(userRepo, config), middleware.AdminMiddleware(), eventController.CreateEvent)
		eventRoutes.PUT("/:id", middleware.AuthMiddleware(userRepo, config), middleware.AdminMiddleware(), eventController.UpdateEvent)
		eventRoutes.DELETE("/:id", middleware.AuthMiddleware(userRepo, config), middleware.AdminMiddleware(), eventController.DeleteEvent)

		// Admin route for event tickets
		eventRoutes.GET("/:id/tickets", middleware.AuthMiddleware(userRepo, config), middleware.AdminMiddleware(), ticketController.GetEventTickets)
	}

	// Ticket routes
	ticketRoutes := router.Group("/api/tickets")
	{
		ticketRoutes.Use(middleware.AuthMiddleware(userRepo, config))

		ticketRoutes.POST("", ticketController.BuyTicket)
		ticketRoutes.GET("/my-tickets", ticketController.GetUserTickets)
		ticketRoutes.GET("/:id", middleware.AdminMiddleware(), ticketController.GetTicketByID)
		ticketRoutes.PUT("/:id/cancel", middleware.AdminMiddleware(), ticketController.CancelTicket)
	}

	// Report routes (admin only)
	reportRoutes := router.Group("/api/reports")
	{
		reportRoutes.Use(middleware.AuthMiddleware(userRepo, config))
		reportRoutes.Use(middleware.AdminMiddleware())

		reportRoutes.GET("/summary", reportController.GetSummaryReport)
		reportRoutes.GET("/events/:id", reportController.GetEventReport)
	}

	if config.Environment != "production" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Add health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Event Ticketing API is running",
		})
	})

	return router
}
