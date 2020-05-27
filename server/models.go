package server

//go:generate protoc -I ../protobuf --go_out=plugins=grpc:../protobuf ../protobuf/snmp.proto
//go:generate protoc --doc_out=../protobuf --doc_opt=markdown,../protobuf/README.md ../protobuf/snmp.proto
import "github.com/soniah/gosnmp"

// Logger should implement logging functionality
type Logger interface{}

// CommandServer is the GRPC Service for snmp.Command
type CommandServer struct {
	ListenAddr       string
	DefaultVersion   gosnmp.SnmpVersion
	DefaultCommunity string
	Logger           Logger
	SNMPLogger       gosnmp.Logger
}
