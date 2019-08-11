package main

//go:generate protoc -I ./protobuf --go_out=plugins=grpc:./protobuf ./protobuf/healthcheck.proto

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

type HealthCheckServer struct{}

func (s *HealthCheckServer) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (response *grpc_health_v1.HealthCheckResponse, err error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthCheckServer) Watch(request *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) (err error) {
	err = srv.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
	return err
}
