// Code generated by protoc-gen-go.
// source: persistent.proto
// DO NOT EDIT!

/*
Package protos is a generated protocol buffer package.

It is generated from these files:
	persistent.proto
	transaction.proto

It has these top-level messages:
	FuncRecord
	UserData
	PublicKey
	ECPoint
	DeploySetting
	NounceData
	RegChainCodeTable
	Signature
	UserTxHeader
	RegPublicKey
	AuthChaincode
	Fund
	Funddata
	PreassignData
	InitChaincode
	RegChaincode
	RecyclePai
	QueryUser
*/
package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type FuncRecord struct {
	Nouncekey []byte `protobuf:"bytes,1,opt,name=nouncekey,proto3" json:"nouncekey,omitempty"`
	IsSend    bool   `protobuf:"varint,2,opt,name=isSend" json:"isSend,omitempty"`
}

func (m *FuncRecord) Reset()                    { *m = FuncRecord{} }
func (m *FuncRecord) String() string            { return proto.CompactTextString(m) }
func (*FuncRecord) ProtoMessage()               {}
func (*FuncRecord) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type UserData struct {
	Pais          int64                      `protobuf:"varint,1,opt,name=pais" json:"pais,omitempty"`
	Pk            *PublicKey                 `protobuf:"bytes,2,opt,name=pk" json:"pk,omitempty"`
	LastActive    *google_protobuf.Timestamp `protobuf:"bytes,3,opt,name=lastActive" json:"lastActive,omitempty"`
	ManagedRegion string                     `protobuf:"bytes,5,opt,name=managedRegion" json:"managedRegion,omitempty"`
	LastFund      *FuncRecord                `protobuf:"bytes,6,opt,name=lastFund" json:"lastFund,omitempty"`
	Authcodes     []uint32                   `protobuf:"varint,10,rep,packed,name=authcodes" json:"authcodes,omitempty"`
}

func (m *UserData) Reset()                    { *m = UserData{} }
func (m *UserData) String() string            { return proto.CompactTextString(m) }
func (*UserData) ProtoMessage()               {}
func (*UserData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *UserData) GetPk() *PublicKey {
	if m != nil {
		return m.Pk
	}
	return nil
}

func (m *UserData) GetLastActive() *google_protobuf.Timestamp {
	if m != nil {
		return m.LastActive
	}
	return nil
}

func (m *UserData) GetLastFund() *FuncRecord {
	if m != nil {
		return m.LastFund
	}
	return nil
}

type PublicKey struct {
	Curvetype int32    `protobuf:"varint,1,opt,name=curvetype" json:"curvetype,omitempty"`
	P         *ECPoint `protobuf:"bytes,2,opt,name=p" json:"p,omitempty"`
}

func (m *PublicKey) Reset()                    { *m = PublicKey{} }
func (m *PublicKey) String() string            { return proto.CompactTextString(m) }
func (*PublicKey) ProtoMessage()               {}
func (*PublicKey) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *PublicKey) GetP() *ECPoint {
	if m != nil {
		return m.P
	}
	return nil
}

type ECPoint struct {
	X []byte `protobuf:"bytes,1,opt,name=x,proto3" json:"x,omitempty"`
	Y []byte `protobuf:"bytes,2,opt,name=y,proto3" json:"y,omitempty"`
}

func (m *ECPoint) Reset()                    { *m = ECPoint{} }
func (m *ECPoint) String() string            { return proto.CompactTextString(m) }
func (*ECPoint) ProtoMessage()               {}
func (*ECPoint) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

// the table of registared chaincode
type DeploySetting struct {
	DebugMode      bool  `protobuf:"varint,1,opt,name=debugMode" json:"debugMode,omitempty"`
	NetworkCode    int32 `protobuf:"varint,2,opt,name=networkCode" json:"networkCode,omitempty"`
	TotalPais      int64 `protobuf:"varint,10,opt,name=totalPais" json:"totalPais,omitempty"`
	UnassignedPais int64 `protobuf:"varint,11,opt,name=unassignedPais" json:"unassignedPais,omitempty"`
}

func (m *DeploySetting) Reset()                    { *m = DeploySetting{} }
func (m *DeploySetting) String() string            { return proto.CompactTextString(m) }
func (*DeploySetting) ProtoMessage()               {}
func (*DeploySetting) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type NounceData struct {
	Txid       string                     `protobuf:"bytes,1,opt,name=txid" json:"txid,omitempty"`
	NounceTime *google_protobuf.Timestamp `protobuf:"bytes,2,opt,name=nounceTime" json:"nounceTime,omitempty"`
	FromLast   *FuncRecord                `protobuf:"bytes,3,opt,name=fromLast" json:"fromLast,omitempty"`
	ToLast     *FuncRecord                `protobuf:"bytes,4,opt,name=toLast" json:"toLast,omitempty"`
}

func (m *NounceData) Reset()                    { *m = NounceData{} }
func (m *NounceData) String() string            { return proto.CompactTextString(m) }
func (*NounceData) ProtoMessage()               {}
func (*NounceData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *NounceData) GetNounceTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.NounceTime
	}
	return nil
}

func (m *NounceData) GetFromLast() *FuncRecord {
	if m != nil {
		return m.FromLast
	}
	return nil
}

func (m *NounceData) GetToLast() *FuncRecord {
	if m != nil {
		return m.ToLast
	}
	return nil
}

type RegChainCodeTable struct {
	T map[uint32]string `protobuf:"bytes,1,rep,name=t" json:"t,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *RegChainCodeTable) Reset()                    { *m = RegChainCodeTable{} }
func (m *RegChainCodeTable) String() string            { return proto.CompactTextString(m) }
func (*RegChainCodeTable) ProtoMessage()               {}
func (*RegChainCodeTable) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *RegChainCodeTable) GetT() map[uint32]string {
	if m != nil {
		return m.T
	}
	return nil
}

func init() {
	proto.RegisterType((*FuncRecord)(nil), "protos.FuncRecord")
	proto.RegisterType((*UserData)(nil), "protos.UserData")
	proto.RegisterType((*PublicKey)(nil), "protos.PublicKey")
	proto.RegisterType((*ECPoint)(nil), "protos.ECPoint")
	proto.RegisterType((*DeploySetting)(nil), "protos.DeploySetting")
	proto.RegisterType((*NounceData)(nil), "protos.NounceData")
	proto.RegisterType((*RegChainCodeTable)(nil), "protos.RegChainCodeTable")
}

func init() { proto.RegisterFile("persistent.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 528 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x93, 0xdd, 0x6e, 0xd3, 0x30,
	0x14, 0xc7, 0xe5, 0x76, 0x2d, 0xed, 0x69, 0x0b, 0x9b, 0x85, 0x50, 0x54, 0x81, 0x08, 0x11, 0xa0,
	0x8a, 0x8b, 0x4c, 0x1a, 0x5c, 0xa0, 0xdd, 0x41, 0xb7, 0x09, 0x89, 0x0f, 0x55, 0x5e, 0x79, 0x00,
	0x27, 0xf1, 0x32, 0xab, 0xa9, 0x1d, 0xc5, 0x27, 0xa5, 0x79, 0x10, 0xde, 0x85, 0xd7, 0xe2, 0x0d,
	0x90, 0x9d, 0x64, 0x19, 0x5f, 0xbb, 0xaa, 0xcf, 0xbf, 0xbf, 0x93, 0xf3, 0xf1, 0xb7, 0xe1, 0x30,
	0x17, 0x85, 0x91, 0x06, 0x85, 0xc2, 0x30, 0x2f, 0x34, 0x6a, 0x3a, 0x74, 0x3f, 0x66, 0xfe, 0x34,
	0xd5, 0x3a, 0xcd, 0xc4, 0xb1, 0x0b, 0xa3, 0xf2, 0xea, 0x18, 0xe5, 0x56, 0x18, 0xe4, 0xdb, 0xbc,
	0x06, 0x83, 0xf7, 0x00, 0x17, 0xa5, 0x8a, 0x99, 0x88, 0x75, 0x91, 0xd0, 0xc7, 0x30, 0x56, 0xba,
	0x54, 0xb1, 0xd8, 0x88, 0xca, 0x23, 0x3e, 0x59, 0x4c, 0x59, 0x27, 0xd0, 0x47, 0x30, 0x94, 0xe6,
	0x52, 0xa8, 0xc4, 0xeb, 0xf9, 0x64, 0x31, 0x62, 0x4d, 0x14, 0xfc, 0x24, 0x30, 0xfa, 0x6a, 0x44,
	0x71, 0xc6, 0x91, 0x53, 0x0a, 0x07, 0x39, 0x97, 0xc6, 0x65, 0xf7, 0x99, 0x3b, 0xd3, 0x67, 0xd0,
	0xcb, 0x37, 0x2e, 0x69, 0x72, 0x72, 0x54, 0x17, 0x36, 0xe1, 0xaa, 0x8c, 0x32, 0x19, 0x7f, 0x14,
	0x15, 0xeb, 0xe5, 0x1b, 0x7a, 0x0a, 0x90, 0x71, 0x83, 0xef, 0x62, 0x94, 0x3b, 0xe1, 0xf5, 0x1d,
	0x3a, 0x0f, 0xeb, 0xee, 0xc3, 0xb6, 0xfb, 0x70, 0xdd, 0x76, 0xcf, 0x6e, 0xd1, 0xf4, 0x39, 0xcc,
	0xb6, 0x5c, 0xf1, 0x54, 0x24, 0x4c, 0xa4, 0x52, 0x2b, 0x6f, 0xe0, 0x93, 0xc5, 0x98, 0xfd, 0x2e,
	0xd2, 0x10, 0x46, 0x36, 0xe7, 0xa2, 0x54, 0x89, 0x37, 0x74, 0xdf, 0xa7, 0x6d, 0x2b, 0xdd, 0x06,
	0xd8, 0x0d, 0x63, 0x77, 0xc1, 0x4b, 0xbc, 0x8e, 0x75, 0x22, 0x8c, 0x07, 0x7e, 0x7f, 0x31, 0x63,
	0x9d, 0x10, 0x7c, 0x80, 0xf1, 0xcd, 0x00, 0x16, 0x8d, 0xcb, 0x62, 0x27, 0xb0, 0xca, 0x85, 0x1b,
	0x7c, 0xc0, 0x3a, 0x81, 0x3e, 0x01, 0x92, 0x37, 0xc3, 0x3f, 0x68, 0x2b, 0x9e, 0x2f, 0x57, 0x5a,
	0x2a, 0x64, 0x24, 0x0f, 0x5e, 0xc0, 0xbd, 0x26, 0xa2, 0x53, 0x20, 0xfb, 0x66, 0xed, 0x64, 0x6f,
	0xa3, 0xca, 0xe5, 0x4d, 0x19, 0xa9, 0x82, 0xef, 0x04, 0x66, 0x67, 0x22, 0xcf, 0x74, 0x75, 0x29,
	0x10, 0xa5, 0x4a, 0x6d, 0xd5, 0x44, 0x44, 0x65, 0xfa, 0x59, 0x27, 0x75, 0xd5, 0x11, 0xeb, 0x04,
	0xea, 0xc3, 0x44, 0x09, 0xfc, 0xa6, 0x8b, 0xcd, 0xd2, 0xfe, 0xdf, 0x73, 0x5d, 0xdd, 0x96, 0x6c,
	0x3e, 0x6a, 0xe4, 0xd9, 0xca, 0xda, 0x05, 0xce, 0xae, 0x4e, 0xa0, 0x2f, 0xe1, 0x7e, 0xa9, 0xb8,
	0x31, 0x32, 0x55, 0x22, 0x71, 0xc8, 0xc4, 0x21, 0x7f, 0xa8, 0xc1, 0x0f, 0x02, 0xf0, 0xc5, 0x5d,
	0x91, 0xd6, 0x7e, 0xdc, 0xcb, 0xc4, 0xf5, 0x33, 0x66, 0xee, 0x6c, 0xbd, 0xad, 0x2f, 0x91, 0xb5,
	0xaf, 0xd9, 0xc4, 0x9d, 0xde, 0x76, 0xb4, 0x75, 0xed, 0xaa, 0xd0, 0xdb, 0x4f, 0xdc, 0x60, 0x73,
	0x2b, 0xfe, 0xe9, 0x5a, 0xcb, 0xd0, 0x57, 0x30, 0x44, 0xed, 0xe8, 0x83, 0xff, 0xd2, 0x0d, 0x11,
	0x54, 0x70, 0xc4, 0x44, 0xba, 0xbc, 0xe6, 0x52, 0xd9, 0x85, 0xac, 0x79, 0x94, 0xd9, 0x82, 0x04,
	0x3d, 0xe2, 0xf7, 0x17, 0x93, 0x13, 0xbf, 0xcd, 0xfd, 0x8b, 0x0a, 0xd7, 0xe7, 0x0a, 0x8b, 0x8a,
	0x11, 0x9c, 0xbf, 0x81, 0x61, 0x1d, 0xd0, 0x43, 0xe8, 0xb7, 0xcf, 0x66, 0xc6, 0xec, 0x91, 0x3e,
	0x84, 0xc1, 0x8e, 0x67, 0x65, 0x3d, 0xf3, 0x98, 0xd5, 0xc1, 0x69, 0xef, 0x2d, 0x89, 0xea, 0xf7,
	0xf9, 0xfa, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0x25, 0xc1, 0x83, 0x66, 0xba, 0x03, 0x00, 0x00,
}
