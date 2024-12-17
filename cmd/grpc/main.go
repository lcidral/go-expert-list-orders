package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	dbConn "go-expert-list-orders/internal/infra/database/postgres"
	repository "go-expert-list-orders/internal/infra/repository/postgres"
	grpcHandler "go-expert-list-orders/internal/interfaces/grpc" // Alias para evitar conflito
	"go-expert-list-orders/internal/pb"                          // Protobuf gerado
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

	grpcServer := grpc.NewServer()
	orderService := grpcHandler.NewOrderService(listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("gRPC server running on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	<-stop

	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped gracefully")
}
