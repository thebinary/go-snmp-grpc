package main

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

// HealthCheckServer is implementation of google grpc_health_v1 interface
type HealthCheckServer struct{}

// Check implements the Check method of the interface
func (s *HealthCheckServer) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (response *grpc_health_v1.HealthCheckResponse, err error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch implements the Watch method of the interface
func (s *HealthCheckServer) Watch(request *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) (err error) {
	err = srv.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
	return err
}
