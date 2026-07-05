package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"booking-service/internal/config"
	"booking-service/internal/handlers"
	"booking-service/internal/middleware"
	"booking-service/internal/repository"
	"booking-service/internal/service"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DBConnectionString())
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)

	userRepo := repository.NewUserRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	slotRepo := repository.NewSlotRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	roomService := service.NewRoomService(roomRepo)
	scheduleService := service.NewScheduleService(scheduleRepo, slotRepo, roomRepo)
	slotService := service.NewSlotService(slotRepo, scheduleRepo, roomRepo, cfg.SlotsLookupDays)
	bookingService := service.NewBookingService(bookingRepo, slotRepo, userRepo, db)

	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)
	roomHandler := handlers.NewRoomHandler(roomService)
	scheduleHandler := handlers.NewScheduleHandler(scheduleService)
	slotHandler := handlers.NewSlotHandler(slotService, roomService)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	infoHandler := handlers.NewInfoHandler()

	router := mux.NewRouter()
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/dummyLogin", authHandler.DummyLogin).Methods("POST")
	router.HandleFunc("/_info", infoHandler.Info).Methods("GET")

	api := router.NewRoute().Subrouter()
	api.Use(middleware.JWTAuth(cfg.JWTSecret))

	api.HandleFunc("/rooms/list", roomHandler.ListRooms).Methods("GET")
	api.HandleFunc("/rooms/create", roomHandler.CreateRoom).Methods("POST")
	api.HandleFunc("/rooms/{roomId}/schedule/create", scheduleHandler.CreateSchedule).Methods("POST")
	api.HandleFunc("/rooms/{roomId}/slots/list", slotHandler.ListAvailableSlots).Methods("GET")
	api.HandleFunc("/bookings/create", bookingHandler.CreateBooking).Methods("POST")
	api.HandleFunc("/bookings/list", bookingHandler.ListAllBookings).Methods("GET")
	api.HandleFunc("/bookings/my", bookingHandler.ListMyBookings).Methods("GET")
	api.HandleFunc("/bookings/{bookingId}/cancel", bookingHandler.CancelBooking).Methods("POST")

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
