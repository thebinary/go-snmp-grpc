package server

import (
	"context"
	"io"
	"log"

	"github.com/soniah/gosnmp"
	pb "github.com/thebinary/go-snmp-grpc/protobuf"
	"google.golang.org/grpc/metadata"
)

// Walk implements snmpwalk command as a oneshot data
func (s *CommandServer) Walk(ctx context.Context, oid *pb.Oid) (snmpPacket *pb.SnmpPDUs, err error) {
	md, _ := metadata.FromIncomingContext(ctx)
	snmp, err := s.snmpConnectionFromMetadata(md)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = snmp.Connect()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer snmp.Conn.Close()

	vars, err := snmp.WalkAll(oid.Oid)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	pbVars := make([]*pb.SnmpPDU, len(vars))
	for i, v := range vars {
		pbVars[i] = ToPbSnmpPDU(v)
	}
	return &pb.SnmpPDUs{
		Pdus: pbVars,
	}, nil
}

// StreamWalk implements snmpwalk command in a stream
func (s *CommandServer) StreamWalk(oid *pb.Oid, srv pb.Command_StreamWalkServer) error {
	ctx := srv.Context()
	md, _ := metadata.FromIncomingContext(ctx)
	snmp, err := s.snmpConnectionFromMetadata(md)
	if err != nil {
		log.Println(err)
		return err
	}

	err = snmp.Connect()
	if err != nil {
		log.Println(err)
		return err
	}
	defer snmp.Conn.Close()

	err = snmp.Walk(oid.Oid, func(pdu gosnmp.SnmpPDU) error {
		err := srv.Send(ToPbSnmpPDU(pdu))
		if err != nil {
			return err
		}
		return nil
	})

	if err == io.EOF {
		log.Printf("Send error: %v", err)
	} else if err != nil {
		log.Printf("%v", err)
	}
	return err
}
