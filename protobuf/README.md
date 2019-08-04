# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [protobuf/snmp.proto](#protobuf/snmp.proto)
    - [Oid](#protobuf.Oid)
    - [OidList](#protobuf.OidList)
    - [SnmpPDU](#protobuf.SnmpPDU)
    - [SnmpPDUs](#protobuf.SnmpPDUs)
    - [SnmpPacket](#protobuf.SnmpPacket)
  
    - [Asn1BER](#protobuf.Asn1BER)
  
  
    - [Command](#protobuf.Command)
  

- [Scalar Value Types](#scalar-value-types)



<a name="protobuf/snmp.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## protobuf/snmp.proto
SNMP Protobuf.

A protobuf interface to SNMP functions.


<a name="protobuf.Oid"></a>

### Oid
Represents SNMP OID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| oid | [string](#string) |  | SNMP OID |






<a name="protobuf.OidList"></a>

### OidList
Represents list of SNMP Oids


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| oids | [string](#string) | repeated | List of SNMP OIDs |






<a name="protobuf.SnmpPDU"></a>

### SnmpPDU
Represents a single SNMP PDU
consisting of oid, type and a value.

The value is in any one of the following fields,
determined by the type of value it stores.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Name | [string](#string) |  | OID |
| Type | [Asn1BER](#protobuf.Asn1BER) |  | PDU Type (Asn1BER encoding type) |
| I32 | [int32](#int32) |  | Stores 32-bit integer |
| I64 | [int64](#int64) |  | Stores 64-bit signed integer |
| UI32 | [uint32](#uint32) |  | Stores 32-bit unsigned integer |
| UI64 | [uint64](#uint64) |  | Stores 64-bit unsigned integer |
| Str | [string](#string) |  | Stores string |






<a name="protobuf.SnmpPDUs"></a>

### SnmpPDUs
Represents multiple SNMP PDU


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pdus | [SnmpPDU](#protobuf.SnmpPDU) | repeated |  |






<a name="protobuf.SnmpPacket"></a>

### SnmpPacket



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Error | [uint32](#uint32) |  |  |
| Variable | [SnmpPDU](#protobuf.SnmpPDU) | repeated |  |





 


<a name="protobuf.Asn1BER"></a>

### Asn1BER
Asn1Ber Enum Type definitions

| Name | Number | Description |
| ---- | ------ | ----------- |
| EndOfContents | 0 |  |
| UnknownType | 0 |  |
| Boolean | 1 |  |
| Integer | 2 |  |
| BitString | 3 |  |
| OctetString | 4 |  |
| Null | 5 |  |
| ObjectIdentifier | 6 |  |
| ObjectDescription | 7 |  |
| IPAddress | 64 |  |
| Counter32 | 65 |  |
| Gauge32 | 66 |  |
| TimeTicks | 67 |  |
| Opaque | 68 |  |
| NsapAddress | 69 |  |
| Counter64 | 70 |  |
| Uinteger32 | 71 |  |
| OpaqueFloat | 120 |  |
| OpaqueDouble | 121 |  |
| NoSuchObject | 128 |  |
| NoSuchInstance | 129 |  |
| EndOfMibView | 130 |  |


 

 


<a name="protobuf.Command"></a>

### Command
The SNMP command service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Get | [OidList](#protobuf.OidList) | [SnmpPacket](#protobuf.SnmpPacket) | SNMP Get given the list of OIDs |
| Set | [SnmpPDUs](#protobuf.SnmpPDUs) | [SnmpPacket](#protobuf.SnmpPacket) | SNMP Set |
| Walk | [Oid](#protobuf.Oid) | [SnmpPDUs](#protobuf.SnmpPDUs) | SNMP Walk and return all variables in one shot |
| StreamWalk | [Oid](#protobuf.Oid) | [SnmpPDU](#protobuf.SnmpPDU) stream | Stream each SNMP PDU while running SNMP WALK from the given OID. |

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

