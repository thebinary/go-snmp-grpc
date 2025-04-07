//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=Msnmp.proto=github.com/thebinary/go-snmp-grpc/protobuf --go-grpc_opt=Msnmp.proto=github.com/thebinary/go-snmp-grpc/protobuf snmp.proto
package protobuf
