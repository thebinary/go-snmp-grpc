package server

import (
	"context"
	"fmt"
	"log"

	pb "github.com/thebinary/go-snmp-grpc/protobuf"
	"google.golang.org/grpc/metadata"
)

// Get implements snmpget command
func (s *CommandServer) Get(ctx context.Context, oids *pb.OidList) (snmpPacket *pb.SnmpPacket, err error) {
	md, _ := metadata.FromIncomingContext(ctx)
	snmp, err := s.snmpConnectionFromMetadata(md)
	if err != nil {
		log.Printf("[ERR] %v", err)
		return nil, err
	}

	err = snmp.Connect()
	if err != nil {
		return nil, fmt.Errorf("ERR_SNMP_CONN: %v", err)
	}
	defer snmp.Conn.Close()

	packet, err := snmp.Get(oids.Oids)
	if err != nil {
		return nil, fmt.Errorf("ERR_SNMP_GET: %v", err)
	}

	vars := packet.Variables
	pbVars := make([]*pb.SnmpPDU, len(vars))
	for i, v := range vars {
		pbVars[i] = ToPbSnmpPDU(v)
	}

	return &pb.SnmpPacket{
		Error:    uint32(packet.ErrorIndex),
		Variable: pbVars,
	}, nil
}
