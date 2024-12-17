package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go-expert-list-orders/internal/pb"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Não foi possível conectar: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Faz a chamada do método ListOrders
	response, err := client.ListOrders(ctx, &pb.ListOrdersRequest{})
	if err != nil {
		log.Fatalf("Erro ao chamar ListOrders: %v", err)
	}

	log.Println("Orders encontradas:")
	for _, order := range response.Orders {
		log.Printf("ID: %s, CustomerID: %s, Status: %s, TotalAmount: %.2f",
			order.Id, order.CustomerId, order.Status, order.TotalAmount)
	}
}
