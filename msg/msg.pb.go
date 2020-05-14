// Code generated by protoc-gen-go. DO NOT EDIT.
// source: msg/msg.proto

package msg

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

//normal message
type Message struct {
	Type                 int32    `protobuf:"varint,1,opt,name=type,proto3" json:"type,omitempty"`
	Sender               string   `protobuf:"bytes,2,opt,name=sender,proto3" json:"sender,omitempty"`
	Destination          string   `protobuf:"bytes,3,opt,name=destination,proto3" json:"destination,omitempty"`
	Content              []byte   `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0f0a1b324c95b77, []int{0}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *Message) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

func (m *Message) GetDestination() string {
	if m != nil {
		return m.Destination
	}
	return ""
}

func (m *Message) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

//subscribe request
type SubscribeReq struct {
	Topics               []string `protobuf:"bytes,1,rep,name=topics,proto3" json:"topics,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubscribeReq) Reset()         { *m = SubscribeReq{} }
func (m *SubscribeReq) String() string { return proto.CompactTextString(m) }
func (*SubscribeReq) ProtoMessage()    {}
func (*SubscribeReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0f0a1b324c95b77, []int{1}
}

func (m *SubscribeReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubscribeReq.Unmarshal(m, b)
}
func (m *SubscribeReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubscribeReq.Marshal(b, m, deterministic)
}
func (m *SubscribeReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubscribeReq.Merge(m, src)
}
func (m *SubscribeReq) XXX_Size() int {
	return xxx_messageInfo_SubscribeReq.Size(m)
}
func (m *SubscribeReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SubscribeReq.DiscardUnknown(m)
}

var xxx_messageInfo_SubscribeReq proto.InternalMessageInfo

func (m *SubscribeReq) GetTopics() []string {
	if m != nil {
		return m.Topics
	}
	return nil
}

//Auth request
type AuthReq struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Credential           []byte   `protobuf:"bytes,2,opt,name=credential,proto3" json:"credential,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthReq) Reset()         { *m = AuthReq{} }
func (m *AuthReq) String() string { return proto.CompactTextString(m) }
func (*AuthReq) ProtoMessage()    {}
func (*AuthReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0f0a1b324c95b77, []int{2}
}

func (m *AuthReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthReq.Unmarshal(m, b)
}
func (m *AuthReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthReq.Marshal(b, m, deterministic)
}
func (m *AuthReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthReq.Merge(m, src)
}
func (m *AuthReq) XXX_Size() int {
	return xxx_messageInfo_AuthReq.Size(m)
}
func (m *AuthReq) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthReq.DiscardUnknown(m)
}

var xxx_messageInfo_AuthReq proto.InternalMessageInfo

func (m *AuthReq) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *AuthReq) GetCredential() []byte {
	if m != nil {
		return m.Credential
	}
	return nil
}

//ack of auth request
type Ack struct {
	Ack                  bool     `protobuf:"varint,1,opt,name=ack,proto3" json:"ack,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ack) Reset()         { *m = Ack{} }
func (m *Ack) String() string { return proto.CompactTextString(m) }
func (*Ack) ProtoMessage()    {}
func (*Ack) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0f0a1b324c95b77, []int{3}
}

func (m *Ack) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ack.Unmarshal(m, b)
}
func (m *Ack) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ack.Marshal(b, m, deterministic)
}
func (m *Ack) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ack.Merge(m, src)
}
func (m *Ack) XXX_Size() int {
	return xxx_messageInfo_Ack.Size(m)
}
func (m *Ack) XXX_DiscardUnknown() {
	xxx_messageInfo_Ack.DiscardUnknown(m)
}

var xxx_messageInfo_Ack proto.InternalMessageInfo

func (m *Ack) GetAck() bool {
	if m != nil {
		return m.Ack
	}
	return false
}

func init() {
	proto.RegisterType((*Message)(nil), "msg.Message")
	proto.RegisterType((*SubscribeReq)(nil), "msg.SubscribeReq")
	proto.RegisterType((*AuthReq)(nil), "msg.AuthReq")
	proto.RegisterType((*Ack)(nil), "msg.Ack")
}

func init() { proto.RegisterFile("msg/msg.proto", fileDescriptor_d0f0a1b324c95b77) }

var fileDescriptor_d0f0a1b324c95b77 = []byte{
	// 222 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8f, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x85, 0xe5, 0xb8, 0x34, 0xf4, 0x08, 0x08, 0x79, 0x00, 0x4f, 0xc8, 0xca, 0x80, 0x32, 0xc1,
	0xc0, 0xc4, 0x58, 0x26, 0x16, 0x16, 0xf3, 0x0b, 0x12, 0xfb, 0x14, 0xac, 0x12, 0x3b, 0xcd, 0x5d,
	0x07, 0xfe, 0x3d, 0xb2, 0xa1, 0x52, 0xb6, 0xf7, 0xbd, 0xbb, 0xd3, 0x7b, 0x07, 0xd7, 0x13, 0x8d,
	0xcf, 0x13, 0x8d, 0x4f, 0xf3, 0x92, 0x38, 0x29, 0x39, 0xd1, 0xd8, 0x1e, 0xa1, 0xfe, 0x40, 0xa2,
	0x7e, 0x44, 0xa5, 0x60, 0xc3, 0x3f, 0x33, 0x6a, 0x61, 0x44, 0x77, 0x61, 0x8b, 0x56, 0x77, 0xb0,
	0x25, 0x8c, 0x1e, 0x17, 0x5d, 0x19, 0xd1, 0xed, 0xec, 0x3f, 0x29, 0x03, 0x57, 0x1e, 0x89, 0x43,
	0xec, 0x39, 0xa4, 0xa8, 0x65, 0x19, 0xae, 0x2d, 0xa5, 0xa1, 0x76, 0x29, 0x32, 0x46, 0xd6, 0x1b,
	0x23, 0xba, 0xc6, 0x9e, 0xb1, 0x7d, 0x84, 0xe6, 0xf3, 0x34, 0x90, 0x5b, 0xc2, 0x80, 0x16, 0x8f,
	0x39, 0x83, 0xd3, 0x1c, 0x1c, 0x69, 0x61, 0x64, 0xce, 0xf8, 0xa3, 0xf6, 0x15, 0xea, 0xfd, 0x89,
	0xbf, 0xf2, 0xca, 0x0d, 0x54, 0xc1, 0x97, 0x62, 0x3b, 0x5b, 0x05, 0xaf, 0x1e, 0x00, 0xdc, 0x82,
	0x1e, 0x23, 0x87, 0xfe, 0xbb, 0x54, 0x6b, 0xec, 0xca, 0x69, 0xef, 0x41, 0xee, 0xdd, 0x41, 0xdd,
	0x82, 0xec, 0xdd, 0xa1, 0xdc, 0x5d, 0xda, 0x2c, 0xdf, 0xaa, 0x77, 0x31, 0x6c, 0xcb, 0xfb, 0x2f,
	0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x03, 0x63, 0x0f, 0x98, 0x0f, 0x01, 0x00, 0x00,
}
