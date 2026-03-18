// @title Event Service API
// @version 1.0
// @description REST API для управления мероприятиями, участниками и площадками.
// @host localhost:8080
// @BasePath /api
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"event-service/internal/config"
	"event-service/internal/handler"
	"event-service/internal/repository"
	"event-service/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	eventRepo := repository.NewEventRepository(db)
	participantRepo := repository.NewParticipantRepository(db)
	venueRepo := repository.NewVenueRepository(db)

	eventSvc := service.NewEventService(eventRepo)
	participantSvc := service.NewParticipantService(participantRepo, eventRepo)
	venueSvc := service.NewVenueService(venueRepo)

	r := handler.NewRouter(eventSvc, participantSvc, venueSvc)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("сервер запускается на %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("ошибка сервера: %v", err)
	}
}

func connectDB(cfg *config.Config) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", cfg.DSN())
		if err != nil {
			log.Printf("попытка %d: не удалось открыть БД: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		if err = db.Ping(); err != nil {
			log.Printf("попытка %d: не удалось выполнить ping БД: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		log.Println("подключение к базе данных установлено")
		return db, nil
	}
	return nil, fmt.Errorf("не удалось подключиться к БД после 10 попыток: %w", err)
}
