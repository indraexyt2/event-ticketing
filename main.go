package main

import (
	"errors"
	"event-ticketing/config"
	_ "event-ticketing/docs"
	"event-ticketing/routes"
	"event-ticketing/utils"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// @title           EventTicketing API
// @version         1.0
// @description     REST API for EventTicketing
// @contact.name    Indrawansyah
// @contact.email   indra@dev.com
// @host      localhost:8080
// @BasePath  /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @type apiKey
func main() {
	// Load konfigurasi
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or failed to load")
	}

	cfg := config.LoadConfig()

	// Inisialisasi logger
	utils.InitLogger(cfg.Environment)

	// Inisialisasi database
	db, err := config.SetupDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Setup routes
	r := routes.SetupRoutes(db, cfg)
	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	// Channel untuk menangkap signal interupsi
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Jalankan server di goroutine terpisah
	go func() {
		log.Printf("Server running on port %s", cfg.AppPort)

		if cfg.Environment != "production" {
			url := "http://localhost:" + cfg.AppPort + "/swagger/index.html"
			err := exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()

			if err != nil {
				log.Printf("Failed to open browser: %v", err)
			}
		}

		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	// Tunggu signal untuk shutdown
	<-quit
	log.Println("Shutting down server...")

	// Tutup koneksi database
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting database instance: %v", err)
		return
	}
	sqlDB.Close()

	log.Println("Server exited properly")
}
