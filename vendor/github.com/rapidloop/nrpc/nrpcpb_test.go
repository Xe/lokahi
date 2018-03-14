// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nrpc_test.proto

/*
Package nrpc_test is a generated protocol buffer package.

It is generated from these files:
	nrpc_test.proto

It has these top-level messages:
	DummyMessage
*/
package nrpc_test

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/rapidloop/nrpc"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type DummyMessage struct {
	Foobar string `protobuf:"bytes,1,opt,name=foobar" json:"foobar,omitempty"`
}

func (m *DummyMessage) Reset()                    { *m = DummyMessage{} }
func (m *DummyMessage) String() string            { return proto.CompactTextString(m) }
func (*DummyMessage) ProtoMessage()               {}
func (*DummyMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DummyMessage) GetFoobar() string {
	if m != nil {
		return m.Foobar
	}
	return ""
}

func init() {
	proto.RegisterType((*DummyMessage)(nil), "DummyMessage")
}

func init() { proto.RegisterFile("nrpc_test.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 86 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xcf, 0x2b, 0x2a, 0x48,
	0x8e, 0x2f, 0x49, 0x2d, 0x2e, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x97, 0xe2, 0x02, 0x09, 0x40,
	0xd8, 0x4a, 0x6a, 0x5c, 0x3c, 0x2e, 0xa5, 0xb9, 0xb9, 0x95, 0xbe, 0xa9, 0xc5, 0xc5, 0x89, 0xe9,
	0xa9, 0x42, 0x62, 0x5c, 0x6c, 0x69, 0xf9, 0xf9, 0x49, 0x89, 0x45, 0x12, 0x8c, 0x0a, 0x8c, 0x1a,
	0x9c, 0x41, 0x50, 0x5e, 0x12, 0x1b, 0x58, 0xb9, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xe8, 0x51,
	0x4a, 0xcb, 0x4d, 0x00, 0x00, 0x00,
}
