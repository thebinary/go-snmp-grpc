package main

//go:generate protoc -I ./protobuf --go_out=plugins=grpc:./protobuf ./protobuf/healthcheck.proto

import (
	"context"

	pb "github.com/thebinary/go-snmp-grpc/protobuf"
)

type HealthCheckServer struct{}

func (s *HealthCheckServer) Check(ctx context.Context, request *pb.HealthCheckRequest) (response *pb.HealthCheckResponse, err error) {
	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthCheckServer) Watch(request *pb.HealthCheckRequest, srv pb.Health_WatchServer) (err error) {
	err = srv.Send(&pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	})
	return err
}
