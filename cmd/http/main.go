package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-expert-list-orders/internal/infra/database/postgres"
	repository "go-expert-list-orders/internal/infra/repository/postgres"
	"go-expert-list-orders/internal/interfaces/handler"
	"go-expert-list-orders/internal/usecase"
)

func main() {
	db, err := postgres.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	orderRepo := repository.NewPostgresOrderRepository(db)

	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepo)
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepo)

	orderHandler := handler.NewOrderHandler(listOrdersUseCase, createOrderUseCase)

	mux := http.NewServeMux()
	mux.Handle("/order", orderHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("REST server running on :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}
