package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"go-expert-list-orders/internal/graph"
	dbConn "go-expert-list-orders/internal/infra/database/postgres"       // Alias para dbConn
	repository "go-expert-list-orders/internal/infra/repository/postgres" // Alias para repository
	"go-expert-list-orders/internal/usecase"
)

func main() {
	db, err := dbConn.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	orderRepo := repository.NewPostgresOrderRepository(db)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepo)

	resolver := graph.NewResolver(listOrdersUseCase)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))

	mux := http.NewServeMux()
	mux.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("GraphQL server running on :8081")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("GraphQL server error: %v", err)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}
