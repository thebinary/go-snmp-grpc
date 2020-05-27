package main

import (
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	pb "github.com/thebinary/go-snmp-grpc/protobuf"
)

func main() {
	lis, err := net.Listen("tcp", commandSrv.ListenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Listening on ", commandSrv.ListenAddr)

	s := grpc.NewServer()
	pb.RegisterCommandServer(s, &commandSrv)
	grpc_health_v1.RegisterHealthServer(s, &HealthCheckServer{})

	if config.MetricsEnabled {
		log.Printf("Metrics listening on address=%s, path=%s", config.MetricsAddr, config.MetricsPath)
		grpc_prometheus.Register(s)
		http.Handle(config.MetricsPath, promhttp.Handler())

		go func() {
			log.Fatal(http.ListenAndServe(config.MetricsAddr, nil))
		}()
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
