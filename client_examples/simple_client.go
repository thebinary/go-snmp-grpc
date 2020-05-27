package main

import (
	"context"
	"flag"
	"io"
	"log"

	pb "github.com/thebinary/go-snmp-grpc/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

func main() {
	var serverAddr string
	var target string
	var community string
	var targetPort string
	flag.StringVar(&serverAddr, "addr", "127.0.0.1:8161", "grpc server address to connect to")
	flag.StringVar(&target, "target", "", "snmp target")
	flag.StringVar(&community, "community", "public", "snmp community")
	flag.StringVar(&targetPort, "port", "161", "snmp target port")
	flag.Parse()

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Panicf("grpc dial error: %v", err)
	}
	defer conn.Close()

	hclient := grpc_health_v1.NewHealthClient(conn)
	d, err := hclient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "GSNMP"})
	log.Println(d, err)

	client := pb.NewCommandClient(conn)

	data := metadata.New(map[string]string{
		"snmp-target":    target,
		"snmp-community": community,
		"snmp-port":      targetPort,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), data)

	sPacket, err := client.Get(ctx, &pb.OidList{
		Oids: []string{
			// OctetString: SNMPv2-MIB::sysDescr.0
			".1.3.6.1.2.1.1.1.0",
			// ObjectIdentifier: SNMPv2-MIB::sysObjectID.0
			".1.3.6.1.2.1.1.2.0",
			// TimeTicks: DISMAN-EVENT-MIB::sysUpTimeInstance
			".1.3.6.1.2.1.1.3.0",
			// Integer: SNMPv2-MIB::sysServices.0
			".1.3.6.1.2.1.1.7.0",
			// Gauge32: IF-MIB::ifSpeed.1
			".1.3.6.1.2.1.2.2.1.5.1",
			// Counter32: IP-MIB::ipSystemStatsInReceives.ipv4
			".1.3.6.1.2.1.4.31.1.1.3.1",
			// Counter64: IP-MIB::ipSystemStatsHCInReceives.ipv4
			".1.3.6.1.2.1.4.31.1.1.4.1",
			//NoSuchInstance
			".1.3.6.1.2.1.1.7.10",
		},
	})

	if err != nil {
		log.Printf("[ERR] %v", err)
	}

	log.Println(sPacket.GetVariable())

	wClient, err := client.StreamWalk(ctx, &pb.Oid{Oid: ".1.3.6.1.2.1.31.1.1.1.1"})
	if err != nil {
		log.Println(err)
	} else {
		for {
			p, err := wClient.Recv()
			log.Println(p.GetType(), p.GetStr(), err)
			if err == io.EOF {
				break
			} else if err != nil {
				log.Printf("err: %v", err)
				break
			}
		}
	}
}
