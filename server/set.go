package server

import (
	"context"
	"fmt"
	"log"

	"github.com/soniah/gosnmp"
	pb "github.com/thebinary/go-snmp-grpc/protobuf"
	"google.golang.org/grpc/metadata"
)

// Set implements snmpset command
func (s *CommandServer) Set(ctx context.Context, pdus *pb.SnmpPDUs) (snmpPacket *pb.SnmpPacket, err error) {
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

	setPdus := make([]gosnmp.SnmpPDU, len(pdus.Pdus))
	for i, pdu := range pdus.Pdus {
		setPdus[i] = ToGoSnmpPDU(pdu)
	}
	packet, err := snmp.Set(setPdus)
	if err != nil {
		return nil, fmt.Errorf("ERR_SNMP_SET: %v", err)
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
