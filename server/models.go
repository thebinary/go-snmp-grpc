package server

//go:generate protoc --doc_out=../protobuf --doc_opt=markdown,../protobuf/README.md ../protobuf/snmp.proto
import (
	"github.com/soniah/gosnmp"
	"github.com/thebinary/go-snmp-grpc/protobuf"
)

// Logger should implement logging functionality
type Logger interface{}

// CommandServer is the GRPC Service for snmp.Command
type CommandServer struct {
	ListenAddr       string
	DefaultVersion   gosnmp.SnmpVersion
	DefaultCommunity string
	Logger           Logger
	SNMPLogger       gosnmp.Logger

	protobuf.UnimplementedCommandServer
}
