package main

//go:generate protoc -I ./protobuf --go_out=plugins=grpc:./protobuf ./protobuf/snmp.proto
//go:generate protoc --doc_out=./protobuf --doc_opt=markdown,./protobuf/README.md ./protobuf/snmp.proto
import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/soniah/gosnmp"
	pb "github.com/thebinary/go-snmp-grpc/protobuf"
	"google.golang.org/grpc/metadata"
)

func (s *CommandServer) snmpConnectionFromMetadata(md metadata.MD) (snmp *gosnmp.GoSNMP, err error) {
	mdTarget := md.Get("snmp-target")
	if mdTarget == nil {
		return nil, fmt.Errorf("snmp-target not given")
	}
	target := mdTarget[0]

	mdPort := md.Get("snmp-port")
	port := gosnmp.Default.Port
	if mdPort != nil {
		iport, err := strconv.Atoi(mdPort[0])
		if err != nil {
			return nil, fmt.Errorf("invalid snmp-port: %s", mdPort[0])
		}
		port = uint16(iport)
	}

	community := s.DefaultCommunity
	mdCommunity := md.Get("snmp-community")
	if mdCommunity != nil {
		community = mdCommunity[0]
	}

	version := s.DefaultVersion
	mdVersion := md.Get("snmp-version")
	if mdVersion != nil {
		switch mdVersion[0] {
		case "1":
			version = gosnmp.Version1
		case "2", "2c":
			version = gosnmp.Version2c
		case "3":
			return nil, fmt.Errorf("snmp-version 3 not supported yet")
		default:
			return nil, fmt.Errorf("invalid snmp-version; %s", mdVersion[0])
		}
	}

	return &gosnmp.GoSNMP{
		Target:    target,
		Port:      port,
		Community: community,
		Version:   version,
		Timeout:   gosnmp.Default.Timeout,
		Retries:   gosnmp.Default.Retries,
		MaxOids:   gosnmp.Default.MaxOids,
	}, nil
}

func ToPbSnmpPDU(pdu gosnmp.SnmpPDU) (pbSnmpPdu *pb.SnmpPDU) {
	pbSnmpPdu = &pb.SnmpPDU{
		Name: pdu.Name,
		Type: pb.Asn1BER(pdu.Type),
	}

	switch pdu.Type {
	case gosnmp.OctetString:
		pbSnmpPdu.Value = &pb.SnmpPDU_Str{
			Str: string(pdu.Value.([]byte)),
		}
	case gosnmp.TimeTicks:
		pbSnmpPdu.Value = &pb.SnmpPDU_UI64{
			UI64: uint64(pdu.Value.(uint)),
		}
	case gosnmp.Integer:
		pbSnmpPdu.Value = &pb.SnmpPDU_I32{
			I32: int32(pdu.Value.(int)),
		}
	}

	return pbSnmpPdu
}

func ToGoSnmpPDU(pbSnmpPdu *pb.SnmpPDU) (pdu gosnmp.SnmpPDU) {
	pdu = gosnmp.SnmpPDU{
		Name: pbSnmpPdu.Name,
		Type: gosnmp.Asn1BER(pbSnmpPdu.Type),
	}

	switch pbSnmpPdu.Type {
	case pb.Asn1BER_OctetString:
		pdu.Value = pbSnmpPdu.GetStr()
	case pb.Asn1BER_TimeTicks:
		pdu.Value = pbSnmpPdu.GetUI64()
	case pb.Asn1BER_Integer:
		pdu.Value = pbSnmpPdu.GetI32()
	}

	return pdu
}

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
