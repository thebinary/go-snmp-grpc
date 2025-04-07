package server

import (
	"fmt"
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
		Logger:    s.SNMPLogger,
	}, nil
}

// ToPbSnmpPDU helper comverts SnmpPDU to protobuf SnmpPDU
func ToPbSnmpPDU(pdu gosnmp.SnmpPDU) (pbSnmpPdu *pb.SnmpPDU) {
	pbSnmpPdu = &pb.SnmpPDU{
		Name: pdu.Name,
		Type: pb.Asn1BER(pdu.Type),
	}

	switch pdu.Type {
	case gosnmp.OctetString:
		pbSnmpPdu.Value = &pb.SnmpPDU_Bytes{
			Bytes: pdu.Value.([]byte),
		}
	case gosnmp.Gauge32:
		fallthrough
	case gosnmp.Counter32:
		fallthrough
	case gosnmp.TimeTicks:
		var val uint32
		var ok bool
		if val, ok = pdu.Value.(uint32); !ok {
			val = uint32(pdu.Value.(uint))
		}
		pbSnmpPdu.Value = &pb.SnmpPDU_UI32{
			UI32: val,
		}
	case gosnmp.Counter64:
		var val uint64
		var ok bool
		if val, ok = pdu.Value.(uint64); !ok {
			val = uint64(pdu.Value.(uint))
		}
		pbSnmpPdu.Value = &pb.SnmpPDU_UI64{
			UI64: val,
		}
	case gosnmp.Integer:
		pbSnmpPdu.Value = &pb.SnmpPDU_I32{
			I32: int32(pdu.Value.(int)),
		}
	case gosnmp.ObjectIdentifier:
		pbSnmpPdu.Value = &pb.SnmpPDU_Str{
			Str: pdu.Value.(string),
		}
	}

	return pbSnmpPdu
}

// ToGoSnmpPDU helper converts protobuf snmppdu to go snmpPDU
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
		pdu.Value = int(pbSnmpPdu.GetI32())
	}

	return pdu
}
