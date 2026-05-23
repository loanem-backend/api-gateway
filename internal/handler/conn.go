package handler

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitConnections() *grpc.ClientConn {
	authConn, err := grpc.NewClient(serviceAddr(auth), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("failed connecting to services: %w", err))
	}

	return authConn
}

func CloseConnections(conns ...*grpc.ClientConn) {
	for _, c := range conns {
		if err := c.Close(); err != nil {
			log.Println("failed closing connection:", err)
		}
	}

	log.Println("gRPC connections closed.")
}
