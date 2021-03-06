// Code generated by protoc-gen-go. DO NOT EDIT.
// source: judy.proto

package judy

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

type Heartbeat struct {
	ClientId             string   `protobuf:"bytes,1,opt,name=ClientId,proto3" json:"ClientId,omitempty"`
	MessageId            string   `protobuf:"bytes,2,opt,name=MessageId,proto3" json:"MessageId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Heartbeat) Reset()         { *m = Heartbeat{} }
func (m *Heartbeat) String() string { return proto.CompactTextString(m) }
func (*Heartbeat) ProtoMessage()    {}
func (*Heartbeat) Descriptor() ([]byte, []int) {
	return fileDescriptor_96a8885cde79d710, []int{0}
}

func (m *Heartbeat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Heartbeat.Unmarshal(m, b)
}
func (m *Heartbeat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Heartbeat.Marshal(b, m, deterministic)
}
func (m *Heartbeat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Heartbeat.Merge(m, src)
}
func (m *Heartbeat) XXX_Size() int {
	return xxx_messageInfo_Heartbeat.Size(m)
}
func (m *Heartbeat) XXX_DiscardUnknown() {
	xxx_messageInfo_Heartbeat.DiscardUnknown(m)
}

var xxx_messageInfo_Heartbeat proto.InternalMessageInfo

func (m *Heartbeat) GetClientId() string {
	if m != nil {
		return m.ClientId
	}
	return ""
}

func (m *Heartbeat) GetMessageId() string {
	if m != nil {
		return m.MessageId
	}
	return ""
}

type Task struct {
	Id                   string   `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Payload              []byte   `protobuf:"bytes,2,opt,name=Payload,proto3" json:"Payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Task) Reset()         { *m = Task{} }
func (m *Task) String() string { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()    {}
func (*Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_96a8885cde79d710, []int{1}
}

func (m *Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Task.Unmarshal(m, b)
}
func (m *Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Task.Marshal(b, m, deterministic)
}
func (m *Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Task.Merge(m, src)
}
func (m *Task) XXX_Size() int {
	return xxx_messageInfo_Task.Size(m)
}
func (m *Task) XXX_DiscardUnknown() {
	xxx_messageInfo_Task.DiscardUnknown(m)
}

var xxx_messageInfo_Task proto.InternalMessageInfo

func (m *Task) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Task) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type HeartbeatResponse struct {
	ClientId             string   `protobuf:"bytes,1,opt,name=ClientId,proto3" json:"ClientId,omitempty"`
	MessageId            string   `protobuf:"bytes,2,opt,name=MessageId,proto3" json:"MessageId,omitempty"`
	Tasks                []*Task  `protobuf:"bytes,3,rep,name=tasks,proto3" json:"tasks,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartbeatResponse) Reset()         { *m = HeartbeatResponse{} }
func (m *HeartbeatResponse) String() string { return proto.CompactTextString(m) }
func (*HeartbeatResponse) ProtoMessage()    {}
func (*HeartbeatResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_96a8885cde79d710, []int{2}
}

func (m *HeartbeatResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartbeatResponse.Unmarshal(m, b)
}
func (m *HeartbeatResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartbeatResponse.Marshal(b, m, deterministic)
}
func (m *HeartbeatResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartbeatResponse.Merge(m, src)
}
func (m *HeartbeatResponse) XXX_Size() int {
	return xxx_messageInfo_HeartbeatResponse.Size(m)
}
func (m *HeartbeatResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartbeatResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HeartbeatResponse proto.InternalMessageInfo

func (m *HeartbeatResponse) GetClientId() string {
	if m != nil {
		return m.ClientId
	}
	return ""
}

func (m *HeartbeatResponse) GetMessageId() string {
	if m != nil {
		return m.MessageId
	}
	return ""
}

func (m *HeartbeatResponse) GetTasks() []*Task {
	if m != nil {
		return m.Tasks
	}
	return nil
}

type CommandRequest struct {
	Command              string   `protobuf:"bytes,1,opt,name=Command,proto3" json:"Command,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandRequest) Reset()         { *m = CommandRequest{} }
func (m *CommandRequest) String() string { return proto.CompactTextString(m) }
func (*CommandRequest) ProtoMessage()    {}
func (*CommandRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_96a8885cde79d710, []int{3}
}

func (m *CommandRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandRequest.Unmarshal(m, b)
}
func (m *CommandRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandRequest.Marshal(b, m, deterministic)
}
func (m *CommandRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandRequest.Merge(m, src)
}
func (m *CommandRequest) XXX_Size() int {
	return xxx_messageInfo_CommandRequest.Size(m)
}
func (m *CommandRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CommandRequest proto.InternalMessageInfo

func (m *CommandRequest) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

type CommandResponse struct {
	Response             string   `protobuf:"bytes,1,opt,name=Response,proto3" json:"Response,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandResponse) Reset()         { *m = CommandResponse{} }
func (m *CommandResponse) String() string { return proto.CompactTextString(m) }
func (*CommandResponse) ProtoMessage()    {}
func (*CommandResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_96a8885cde79d710, []int{4}
}

func (m *CommandResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandResponse.Unmarshal(m, b)
}
func (m *CommandResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandResponse.Marshal(b, m, deterministic)
}
func (m *CommandResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandResponse.Merge(m, src)
}
func (m *CommandResponse) XXX_Size() int {
	return xxx_messageInfo_CommandResponse.Size(m)
}
func (m *CommandResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CommandResponse proto.InternalMessageInfo

func (m *CommandResponse) GetResponse() string {
	if m != nil {
		return m.Response
	}
	return ""
}

func init() {
	proto.RegisterType((*Heartbeat)(nil), "judy.Heartbeat")
	proto.RegisterType((*Task)(nil), "judy.Task")
	proto.RegisterType((*HeartbeatResponse)(nil), "judy.HeartbeatResponse")
	proto.RegisterType((*CommandRequest)(nil), "judy.CommandRequest")
	proto.RegisterType((*CommandResponse)(nil), "judy.CommandResponse")
}

func init() { proto.RegisterFile("judy.proto", fileDescriptor_96a8885cde79d710) }

var fileDescriptor_96a8885cde79d710 = []byte{
	// 212 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x90, 0xcb, 0x4a, 0xc5, 0x30,
	0x10, 0x86, 0x39, 0x3d, 0xc7, 0x4b, 0x47, 0xa9, 0x98, 0x55, 0x10, 0x17, 0x25, 0xab, 0x22, 0x58,
	0x44, 0x1f, 0xa1, 0x08, 0x76, 0x21, 0x48, 0xf0, 0x05, 0xa6, 0x64, 0x10, 0xed, 0x25, 0xb5, 0x93,
	0x2e, 0xfa, 0xf6, 0xd2, 0x4b, 0xd2, 0xbd, 0xbb, 0xf9, 0x66, 0x32, 0xfc, 0x5f, 0x06, 0xe0, 0x67,
	0x34, 0x53, 0xde, 0x0f, 0xd6, 0x59, 0x71, 0x9a, 0x6b, 0xf5, 0x0a, 0xf1, 0x1b, 0xe1, 0xe0, 0x2a,
	0x42, 0x27, 0xee, 0xe0, 0xb2, 0x68, 0xbe, 0xa9, 0x73, 0xa5, 0x91, 0x87, 0xf4, 0x90, 0xc5, 0x3a,
	0xb0, 0xb8, 0x87, 0xf8, 0x9d, 0x98, 0xf1, 0x8b, 0x4a, 0x23, 0xa3, 0x65, 0xb8, 0x37, 0xd4, 0x13,
	0x9c, 0x3e, 0x91, 0x6b, 0x91, 0x40, 0x14, 0x76, 0xa3, 0xd2, 0x08, 0x09, 0x17, 0x1f, 0x38, 0x35,
	0x16, 0xd7, 0x9d, 0x6b, 0xed, 0x51, 0x59, 0xb8, 0x0d, 0xc1, 0x9a, 0xb8, 0xb7, 0x1d, 0xd3, 0xff,
	0x05, 0x44, 0x0a, 0x67, 0x0e, 0xb9, 0x66, 0x79, 0x4c, 0x8f, 0xd9, 0xd5, 0x33, 0xe4, 0xcb, 0x4f,
	0x67, 0x27, 0xbd, 0x0e, 0xd4, 0x03, 0x24, 0x85, 0x6d, 0x5b, 0xec, 0x8c, 0xa6, 0xdf, 0x91, 0xd8,
	0xcd, 0x72, 0x5b, 0x67, 0x0b, 0xf3, 0xa8, 0x1e, 0xe1, 0x26, 0xbc, 0xdd, 0xd5, 0x7c, 0xed, 0xd5,
	0x3c, 0x57, 0xe7, 0xcb, 0x45, 0x5f, 0xfe, 0x02, 0x00, 0x00, 0xff, 0xff, 0xbe, 0x36, 0xc8, 0xc2,
	0x5f, 0x01, 0x00, 0x00,
}
