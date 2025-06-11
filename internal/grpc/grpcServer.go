package grpc

import (
	pb "collectionsservice/internal/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

func StartGRPCServer(ser pb.CollectionServiceServer) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(" Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCollectionServiceServer(grpcServer, ser)

	log.Println("gRPC Server started on port 50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Failed to start gRPC server:", err)
	}
}
