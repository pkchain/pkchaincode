// Code generated by protoc-gen-go.
// source: transaction.proto
// DO NOT EDIT!

package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
import google_protobuf1 "github.com/golang/protobuf/ptypes/any"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Signature struct {
	P *ECPoint `protobuf:"bytes,1,opt,name=p" json:"p,omitempty"`
}

func (m *Signature) Reset()                    { *m = Signature{} }
func (m *Signature) String() string            { return proto.CompactTextString(m) }
func (*Signature) ProtoMessage()               {}
func (*Signature) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *Signature) GetP() *ECPoint {
	if m != nil {
		return m.P
	}
	return nil
}

type UserTxHeader struct {
	FundId string                     `protobuf:"bytes,1,opt,name=fundId" json:"fundId,omitempty"`
	Nounce []byte                     `protobuf:"bytes,2,opt,name=nounce,proto3" json:"nounce,omitempty"`
	Ts     *google_protobuf.Timestamp `protobuf:"bytes,3,opt,name=ts" json:"ts,omitempty"`
}

func (m *UserTxHeader) Reset()                    { *m = UserTxHeader{} }
func (m *UserTxHeader) String() string            { return proto.CompactTextString(m) }
func (*UserTxHeader) ProtoMessage()               {}
func (*UserTxHeader) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *UserTxHeader) GetTs() *google_protobuf.Timestamp {
	if m != nil {
		return m.Ts
	}
	return nil
}

// user can register a public key only if it has own some pais
type RegPublicKey struct {
	Pk *PublicKey `protobuf:"bytes,1,opt,name=pk" json:"pk,omitempty"`
}

func (m *RegPublicKey) Reset()                    { *m = RegPublicKey{} }
func (m *RegPublicKey) String() string            { return proto.CompactTextString(m) }
func (*RegPublicKey) ProtoMessage()               {}
func (*RegPublicKey) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *RegPublicKey) GetPk() *PublicKey {
	if m != nil {
		return m.Pk
	}
	return nil
}

type AuthChaincode struct {
	Code uint32 `protobuf:"varint,1,opt,name=code" json:"code,omitempty"`
}

func (m *AuthChaincode) Reset()                    { *m = AuthChaincode{} }
func (m *AuthChaincode) String() string            { return proto.CompactTextString(m) }
func (*AuthChaincode) ProtoMessage()               {}
func (*AuthChaincode) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

// fund tx can be invoked directly by user or other chaincode,
// if the fund tx. is invoked by user, invokeChaincode is 0
// and Funddata is used
// any other chaincode invoked fund tx must add their data fields
// in Fund, or use contract field for another parsing
type Fund struct {
	InvokeChaincode uint32 `protobuf:"varint,4,opt,name=invokeChaincode" json:"invokeChaincode,omitempty"`
	// Types that are valid to be assigned to D:
	//	*Fund_Userfund
	//	*Fund_Null
	//	*Fund_Contract
	D isFund_D `protobuf_oneof:"d"`
}

func (m *Fund) Reset()                    { *m = Fund{} }
func (m *Fund) String() string            { return proto.CompactTextString(m) }
func (*Fund) ProtoMessage()               {}
func (*Fund) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

type isFund_D interface {
	isFund_D()
}

type Fund_Userfund struct {
	Userfund *Funddata `protobuf:"bytes,1,opt,name=userfund,oneof"`
}
type Fund_Null struct {
	Null bool `protobuf:"varint,2,opt,name=null,oneof"`
}
type Fund_Contract struct {
	Contract *google_protobuf1.Any `protobuf:"bytes,3,opt,name=contract,oneof"`
}

func (*Fund_Userfund) isFund_D() {}
func (*Fund_Null) isFund_D()     {}
func (*Fund_Contract) isFund_D() {}

func (m *Fund) GetD() isFund_D {
	if m != nil {
		return m.D
	}
	return nil
}

func (m *Fund) GetUserfund() *Funddata {
	if x, ok := m.GetD().(*Fund_Userfund); ok {
		return x.Userfund
	}
	return nil
}

func (m *Fund) GetNull() bool {
	if x, ok := m.GetD().(*Fund_Null); ok {
		return x.Null
	}
	return false
}

func (m *Fund) GetContract() *google_protobuf1.Any {
	if x, ok := m.GetD().(*Fund_Contract); ok {
		return x.Contract
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Fund) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Fund_OneofMarshaler, _Fund_OneofUnmarshaler, _Fund_OneofSizer, []interface{}{
		(*Fund_Userfund)(nil),
		(*Fund_Null)(nil),
		(*Fund_Contract)(nil),
	}
}

func _Fund_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Fund)
	// d
	switch x := m.D.(type) {
	case *Fund_Userfund:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Userfund); err != nil {
			return err
		}
	case *Fund_Null:
		t := uint64(0)
		if x.Null {
			t = 1
		}
		b.EncodeVarint(2<<3 | proto.WireVarint)
		b.EncodeVarint(t)
	case *Fund_Contract:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Contract); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Fund.D has unexpected type %T", x)
	}
	return nil
}

func _Fund_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Fund)
	switch tag {
	case 1: // d.userfund
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Funddata)
		err := b.DecodeMessage(msg)
		m.D = &Fund_Userfund{msg}
		return true, err
	case 2: // d.null
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.D = &Fund_Null{x != 0}
		return true, err
	case 3: // d.contract
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(google_protobuf1.Any)
		err := b.DecodeMessage(msg)
		m.D = &Fund_Contract{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Fund_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Fund)
	// d
	switch x := m.D.(type) {
	case *Fund_Userfund:
		s := proto.Size(x.Userfund)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Fund_Null:
		n += proto.SizeVarint(2<<3 | proto.WireVarint)
		n += 1
	case *Fund_Contract:
		s := proto.Size(x.Contract)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Funddata struct {
	Pai      uint32 `protobuf:"varint,1,opt,name=pai" json:"pai,omitempty"`
	ToUserId string `protobuf:"bytes,2,opt,name=toUserId" json:"toUserId,omitempty"`
}

func (m *Funddata) Reset()                    { *m = Funddata{} }
func (m *Funddata) String() string            { return proto.CompactTextString(m) }
func (*Funddata) ProtoMessage()               {}
func (*Funddata) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

// system transactions
type PreassignData struct {
	Userid string `protobuf:"bytes,1,opt,name=userid" json:"userid,omitempty"`
	Pais   int64  `protobuf:"varint,2,opt,name=pais" json:"pais,omitempty"`
}

func (m *PreassignData) Reset()                    { *m = PreassignData{} }
func (m *PreassignData) String() string            { return proto.CompactTextString(m) }
func (*PreassignData) ProtoMessage()               {}
func (*PreassignData) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{6} }

type InitChaincode struct {
	Mainsetting     *DeploySetting   `protobuf:"bytes,1,opt,name=mainsetting" json:"mainsetting,omitempty"`
	PreassignedUser []*PreassignData `protobuf:"bytes,2,rep,name=preassignedUser" json:"preassignedUser,omitempty"`
}

func (m *InitChaincode) Reset()                    { *m = InitChaincode{} }
func (m *InitChaincode) String() string            { return proto.CompactTextString(m) }
func (*InitChaincode) ProtoMessage()               {}
func (*InitChaincode) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{7} }

func (m *InitChaincode) GetMainsetting() *DeploySetting {
	if m != nil {
		return m.Mainsetting
	}
	return nil
}

func (m *InitChaincode) GetPreassignedUser() []*PreassignData {
	if m != nil {
		return m.PreassignedUser
	}
	return nil
}

type RegChaincode struct {
	ChaincodeName string `protobuf:"bytes,1,opt,name=chaincodeName" json:"chaincodeName,omitempty"`
	ChaincodeId   uint32 `protobuf:"varint,2,opt,name=chaincodeId" json:"chaincodeId,omitempty"`
}

func (m *RegChaincode) Reset()                    { *m = RegChaincode{} }
func (m *RegChaincode) String() string            { return proto.CompactTextString(m) }
func (*RegChaincode) ProtoMessage()               {}
func (*RegChaincode) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{8} }

type RecyclePai struct {
	DeadUserId string `protobuf:"bytes,1,opt,name=deadUserId" json:"deadUserId,omitempty"`
	ToUserId   string `protobuf:"bytes,2,opt,name=toUserId" json:"toUserId,omitempty"`
}

func (m *RecyclePai) Reset()                    { *m = RecyclePai{} }
func (m *RecyclePai) String() string            { return proto.CompactTextString(m) }
func (*RecyclePai) ProtoMessage()               {}
func (*RecyclePai) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{9} }

type QueryUser struct {
	Userid string `protobuf:"bytes,1,opt,name=userid" json:"userid,omitempty"`
}

func (m *QueryUser) Reset()                    { *m = QueryUser{} }
func (m *QueryUser) String() string            { return proto.CompactTextString(m) }
func (*QueryUser) ProtoMessage()               {}
func (*QueryUser) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{10} }

func init() {
	proto.RegisterType((*Signature)(nil), "protos.Signature")
	proto.RegisterType((*UserTxHeader)(nil), "protos.UserTxHeader")
	proto.RegisterType((*RegPublicKey)(nil), "protos.RegPublicKey")
	proto.RegisterType((*AuthChaincode)(nil), "protos.AuthChaincode")
	proto.RegisterType((*Fund)(nil), "protos.Fund")
	proto.RegisterType((*Funddata)(nil), "protos.Funddata")
	proto.RegisterType((*PreassignData)(nil), "protos.PreassignData")
	proto.RegisterType((*InitChaincode)(nil), "protos.InitChaincode")
	proto.RegisterType((*RegChaincode)(nil), "protos.RegChaincode")
	proto.RegisterType((*RecyclePai)(nil), "protos.RecyclePai")
	proto.RegisterType((*QueryUser)(nil), "protos.QueryUser")
}

func init() { proto.RegisterFile("transaction.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 528 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x53, 0x51, 0x6f, 0xd3, 0x30,
	0x10, 0x5e, 0xd2, 0x6a, 0x6a, 0x6f, 0x8d, 0xda, 0x59, 0x63, 0x2a, 0x95, 0x80, 0x92, 0xf1, 0x50,
	0xed, 0x21, 0x13, 0xe3, 0x01, 0x24, 0x1e, 0xd0, 0xd8, 0x40, 0xad, 0x90, 0x50, 0xf1, 0x06, 0xef,
	0x6e, 0x72, 0xcb, 0x4c, 0x53, 0x3b, 0x8a, 0x1d, 0x44, 0x7e, 0x02, 0xbf, 0x86, 0xbf, 0x88, 0xec,
	0x38, 0x59, 0x57, 0x04, 0x4f, 0xbe, 0xef, 0xee, 0xbb, 0xf3, 0x77, 0x9f, 0x13, 0x38, 0xd4, 0x05,
	0x13, 0x8a, 0xc5, 0x9a, 0x4b, 0x11, 0xe5, 0x85, 0xd4, 0x92, 0xec, 0xdb, 0x43, 0x4d, 0x46, 0x39,
	0x16, 0x8a, 0x2b, 0x8d, 0x42, 0xd7, 0x95, 0xc9, 0xb3, 0x54, 0xca, 0x34, 0xc3, 0x33, 0x8b, 0x56,
	0xe5, 0xed, 0x99, 0xe6, 0x1b, 0x54, 0x9a, 0x6d, 0x72, 0x47, 0x78, 0xbc, 0x4b, 0x60, 0xa2, 0xaa,
	0x4b, 0xe1, 0x29, 0xf4, 0xaf, 0x79, 0x2a, 0x98, 0x2e, 0x0b, 0x24, 0x4f, 0xc0, 0xcb, 0xc7, 0xde,
	0xd4, 0x9b, 0x1d, 0x9c, 0x0f, 0xeb, 0xba, 0x8a, 0x3e, 0x5c, 0x2e, 0x25, 0x17, 0x9a, 0x7a, 0x79,
	0xf8, 0x1d, 0x06, 0x5f, 0x15, 0x16, 0x37, 0x3f, 0xe7, 0xc8, 0x12, 0x2c, 0xc8, 0x31, 0xec, 0xdf,
	0x96, 0x22, 0x59, 0x24, 0xb6, 0xa7, 0x4f, 0x1d, 0x32, 0x79, 0x21, 0x4b, 0x11, 0xe3, 0xd8, 0x9f,
	0x7a, 0xb3, 0x01, 0x75, 0x88, 0x9c, 0x82, 0xaf, 0xd5, 0xb8, 0x63, 0xe7, 0x4f, 0xa2, 0x5a, 0x53,
	0xd4, 0x68, 0x8a, 0x6e, 0x1a, 0xd1, 0xd4, 0xd7, 0x2a, 0x7c, 0x09, 0x03, 0x8a, 0xe9, 0xb2, 0x5c,
	0x65, 0x3c, 0xfe, 0x84, 0x15, 0x79, 0x0e, 0x7e, 0xbe, 0x76, 0xda, 0x0e, 0x1b, 0x6d, 0x6d, 0x99,
	0xfa, 0xf9, 0x3a, 0x3c, 0x81, 0xe0, 0xa2, 0xd4, 0x77, 0x97, 0x77, 0x8c, 0x8b, 0x58, 0x26, 0x48,
	0x08, 0x74, 0xcd, 0x69, 0xbb, 0x02, 0x6a, 0xe3, 0xf0, 0xb7, 0x07, 0xdd, 0x8f, 0xa5, 0x48, 0xc8,
	0x0c, 0x86, 0x5c, 0xfc, 0x90, 0x6b, 0x6c, 0xf9, 0xe3, 0xae, 0xe5, 0xed, 0xa6, 0x49, 0x04, 0xbd,
	0x52, 0x61, 0x61, 0x96, 0x73, 0x02, 0x46, 0x8d, 0x00, 0x33, 0x29, 0x61, 0x9a, 0xcd, 0xf7, 0x68,
	0xcb, 0x21, 0x47, 0xd0, 0x15, 0x65, 0x96, 0xd9, 0xe5, 0x7b, 0xf3, 0x3d, 0x6a, 0x11, 0x39, 0x87,
	0x5e, 0x2c, 0x85, 0x2e, 0x58, 0xac, 0x9d, 0x05, 0x47, 0x7f, 0x59, 0x70, 0x21, 0x2a, 0x33, 0xa9,
	0xe1, 0xbd, 0xef, 0x80, 0x97, 0x84, 0x6f, 0xa0, 0xd7, 0x5c, 0x43, 0x46, 0xd0, 0xc9, 0x19, 0x77,
	0x0b, 0x99, 0x90, 0x4c, 0xa0, 0xa7, 0xa5, 0x79, 0x95, 0x45, 0x62, 0x2f, 0xec, 0xd3, 0x16, 0x87,
	0x6f, 0x21, 0x58, 0x16, 0xc8, 0x94, 0xe2, 0xa9, 0xb8, 0x32, 0xed, 0xc7, 0xb0, 0x6f, 0x54, 0xf2,
	0xf6, 0xc1, 0x6a, 0x64, 0x8c, 0xca, 0x19, 0x57, 0x76, 0x40, 0x87, 0xda, 0x38, 0xfc, 0xe5, 0x41,
	0xb0, 0x10, 0x5c, 0xdf, 0xfb, 0xf0, 0x1a, 0x0e, 0x36, 0x8c, 0x0b, 0x85, 0x5a, 0x73, 0x91, 0x3a,
	0x2b, 0x1e, 0x35, 0x56, 0x5c, 0x61, 0x9e, 0xc9, 0xea, 0xba, 0x2e, 0xd2, 0x6d, 0x26, 0x79, 0x07,
	0xc3, 0xbc, 0xd1, 0x81, 0x89, 0x11, 0x37, 0xf6, 0xa7, 0x9d, 0xed, 0xe6, 0x07, 0x32, 0xe9, 0x2e,
	0x3b, 0xfc, 0x66, 0x3f, 0x86, 0x7b, 0x25, 0x2f, 0x20, 0x88, 0x1b, 0xf0, 0x99, 0x6d, 0xd0, 0xad,
	0xf3, 0x30, 0x49, 0xa6, 0x70, 0xd0, 0x26, 0x9c, 0x3b, 0x01, 0xdd, 0x4e, 0x85, 0x73, 0x00, 0x8a,
	0x71, 0x15, 0x67, 0xb8, 0x64, 0x9c, 0x3c, 0x05, 0x48, 0x90, 0x25, 0xce, 0xcc, 0x7a, 0xe4, 0x56,
	0xe6, 0xbf, 0x56, 0x9f, 0x40, 0xff, 0x4b, 0x89, 0x45, 0x65, 0xe0, 0xbf, 0x6c, 0x5e, 0xd5, 0x7f,
	0xf0, 0xab, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xb6, 0x28, 0xc8, 0xca, 0xdd, 0x03, 0x00, 0x00,
}
