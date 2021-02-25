// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: page.proto

/*
Package cache is a generated protocol buffer package.

It is generated from these files:
	page.proto

It has these top-level messages:
	ResponseCache
	HeaderValue
*/
package degrade

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type ResponseCache struct {
	Status int32                   `protobuf:"varint,1,opt,name=Status,proto3" json:"Status,omitempty"`
	Header map[string]*HeaderValue `protobuf:"bytes,2,rep,name=Header" json:"Header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
	Data   []byte                  `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (m *ResponseCache) Reset()                    { *m = ResponseCache{} }
func (m *ResponseCache) String() string            { return proto.CompactTextString(m) }
func (*ResponseCache) ProtoMessage()               {}
func (*ResponseCache) Descriptor() ([]byte, []int) { return fileDescriptorPage, []int{0} }

func (m *ResponseCache) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *ResponseCache) GetHeader() map[string]*HeaderValue {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *ResponseCache) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type HeaderValue struct {
	Value []string `protobuf:"bytes,1,rep,name=Value" json:"Value,omitempty"`
}

func (m *HeaderValue) Reset()                    { *m = HeaderValue{} }
func (m *HeaderValue) String() string            { return proto.CompactTextString(m) }
func (*HeaderValue) ProtoMessage()               {}
func (*HeaderValue) Descriptor() ([]byte, []int) { return fileDescriptorPage, []int{1} }

func (m *HeaderValue) GetValue() []string {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*ResponseCache)(nil), "cache.responseCache")
	proto.RegisterType((*HeaderValue)(nil), "cache.headerValue")
}

func init() { proto.RegisterFile("page.proto", fileDescriptorPage) }

var fileDescriptorPage = []byte{
	// 231 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8f, 0x41, 0x4b, 0xc4, 0x30,
	0x10, 0x85, 0x49, 0x6b, 0x0b, 0x3b, 0x55, 0x90, 0x41, 0x24, 0xec, 0x29, 0xac, 0x97, 0x5c, 0xcc,
	0xc2, 0x7a, 0x59, 0xbc, 0xaa, 0xe0, 0xc5, 0x4b, 0x04, 0xef, 0x69, 0x1d, 0x5b, 0x51, 0x37, 0xa5,
	0x4d, 0x84, 0xfd, 0x7f, 0xfe, 0x30, 0xe9, 0xa4, 0x87, 0xee, 0xed, 0x3d, 0xde, 0xf7, 0xe6, 0x31,
	0x00, 0xbd, 0x6b, 0xc9, 0xf4, 0x83, 0x0f, 0x1e, 0x8b, 0xc6, 0x35, 0x1d, 0xad, 0x6f, 0xdb, 0xcf,
	0xd0, 0xc5, 0xda, 0x34, 0xfe, 0x67, 0xdb, 0xfa, 0xd6, 0x6f, 0x39, 0xad, 0xe3, 0x07, 0x3b, 0x36,
	0xac, 0x52, 0x6b, 0xf3, 0x27, 0xe0, 0x62, 0xa0, 0xb1, 0xf7, 0x87, 0x91, 0x1e, 0xa6, 0x03, 0x78,
	0x0d, 0xe5, 0x6b, 0x70, 0x21, 0x8e, 0x52, 0x28, 0xa1, 0x0b, 0x3b, 0x3b, 0xdc, 0x43, 0xf9, 0x4c,
	0xee, 0x9d, 0x06, 0x99, 0xa9, 0x5c, 0x57, 0x3b, 0x65, 0x78, 0xd0, 0x9c, 0xb4, 0x4d, 0x42, 0x9e,
	0x0e, 0x61, 0x38, 0xda, 0x99, 0x47, 0x84, 0xb3, 0x47, 0x17, 0x9c, 0xcc, 0x95, 0xd0, 0xe7, 0x96,
	0xf5, 0xfa, 0x05, 0xaa, 0x05, 0x8a, 0x97, 0x90, 0x7f, 0xd1, 0x91, 0x17, 0x57, 0x76, 0x92, 0xa8,
	0xa1, 0xf8, 0x75, 0xdf, 0x91, 0x64, 0xa6, 0x84, 0xae, 0x76, 0x38, 0xaf, 0x75, 0x5c, 0x7a, 0x9b,
	0x12, 0x9b, 0x80, 0xfb, 0x6c, 0x2f, 0x36, 0x37, 0x50, 0x2d, 0x12, 0xbc, 0x82, 0x82, 0x85, 0x14,
	0x2a, 0xd7, 0x2b, 0x9b, 0x4c, 0x5d, 0xf2, 0xcb, 0x77, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff, 0x6f,
	0x05, 0xcf, 0xb0, 0x36, 0x01, 0x00, 0x00,
}