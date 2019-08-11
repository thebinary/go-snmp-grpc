package main

import (
	"log"
	"net"

	pb "github.com/thebinary/go-snmp-grpc/protobuf"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", Server.ListenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Listening on ", Server.ListenAddr)

	s := grpc.NewServer()
	pb.RegisterCommandServer(s, &Server)
	pb.RegisterHealthServer(s, &HealthCheckServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
