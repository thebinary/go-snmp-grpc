package main

import (
	"log"
	"net"

	pb "github.com/thebinary/go-snmp-grpc/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	lis, err := net.Listen("tcp", Server.ListenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Listening on ", Server.ListenAddr)

	s := grpc.NewServer()
	pb.RegisterCommandServer(s, &Server)
	grpc_health_v1.RegisterHealthServer(s, &HealthCheckServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
