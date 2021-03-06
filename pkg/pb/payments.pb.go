// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: payments.proto

package pb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// The request message containing the details needed to pay a storage node.
type PaymentRequest struct {
	// ID of the storage node to be paid
	NodeId               string   `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PaymentRequest) Reset()         { *m = PaymentRequest{} }
func (m *PaymentRequest) String() string { return proto.CompactTextString(m) }
func (*PaymentRequest) ProtoMessage()    {}
func (*PaymentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_payments_f4730bcfc665723c, []int{0}
}
func (m *PaymentRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PaymentRequest.Unmarshal(m, b)
}
func (m *PaymentRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PaymentRequest.Marshal(b, m, deterministic)
}
func (dst *PaymentRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PaymentRequest.Merge(dst, src)
}
func (m *PaymentRequest) XXX_Size() int {
	return xxx_messageInfo_PaymentRequest.Size(m)
}
func (m *PaymentRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PaymentRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PaymentRequest proto.InternalMessageInfo

func (m *PaymentRequest) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

// The response message for payments.
type PaymentResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PaymentResponse) Reset()         { *m = PaymentResponse{} }
func (m *PaymentResponse) String() string { return proto.CompactTextString(m) }
func (*PaymentResponse) ProtoMessage()    {}
func (*PaymentResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_payments_f4730bcfc665723c, []int{1}
}
func (m *PaymentResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PaymentResponse.Unmarshal(m, b)
}
func (m *PaymentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PaymentResponse.Marshal(b, m, deterministic)
}
func (dst *PaymentResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PaymentResponse.Merge(dst, src)
}
func (m *PaymentResponse) XXX_Size() int {
	return xxx_messageInfo_PaymentResponse.Size(m)
}
func (m *PaymentResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PaymentResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PaymentResponse proto.InternalMessageInfo

// The request message containing the details needed to calculate outstanding balance for a storage node.
type CalculateRequest struct {
	// ID of the storage node to be calculated
	NodeId               string   `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CalculateRequest) Reset()         { *m = CalculateRequest{} }
func (m *CalculateRequest) String() string { return proto.CompactTextString(m) }
func (*CalculateRequest) ProtoMessage()    {}
func (*CalculateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_payments_f4730bcfc665723c, []int{2}
}
func (m *CalculateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CalculateRequest.Unmarshal(m, b)
}
func (m *CalculateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CalculateRequest.Marshal(b, m, deterministic)
}
func (dst *CalculateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CalculateRequest.Merge(dst, src)
}
func (m *CalculateRequest) XXX_Size() int {
	return xxx_messageInfo_CalculateRequest.Size(m)
}
func (m *CalculateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CalculateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CalculateRequest proto.InternalMessageInfo

func (m *CalculateRequest) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

// The response message for payment calculations.
type CalculateResponse struct {
	// ID of the storage node calculation made for
	NodeId string `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	// total balance in Storj of outstanding credit
	Total                int64    `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CalculateResponse) Reset()         { *m = CalculateResponse{} }
func (m *CalculateResponse) String() string { return proto.CompactTextString(m) }
func (*CalculateResponse) ProtoMessage()    {}
func (*CalculateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_payments_f4730bcfc665723c, []int{3}
}
func (m *CalculateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CalculateResponse.Unmarshal(m, b)
}
func (m *CalculateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CalculateResponse.Marshal(b, m, deterministic)
}
func (dst *CalculateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CalculateResponse.Merge(dst, src)
}
func (m *CalculateResponse) XXX_Size() int {
	return xxx_messageInfo_CalculateResponse.Size(m)
}
func (m *CalculateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CalculateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CalculateResponse proto.InternalMessageInfo

func (m *CalculateResponse) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *CalculateResponse) GetTotal() int64 {
	if m != nil {
		return m.Total
	}
	return 0
}

// The request message for adjusting the cost of storage/bandwidth for a satelitte.
type AdjustPricesRequest struct {
	// price per gigabyte of bandwidth calculated in Storj
	Bandwidth int64 `protobuf:"varint,1,opt,name=bandwidth,proto3" json:"bandwidth,omitempty"`
	// price for GB/H of storage calculated in Storj
	Storage              int64    `protobuf:"varint,2,opt,name=storage,proto3" json:"storage,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AdjustPricesRequest) Reset()         { *m = AdjustPricesRequest{} }
func (m *AdjustPricesRequest) String() string { return proto.CompactTextString(m) }
func (*AdjustPricesRequest) ProtoMessage()    {}
func (*AdjustPricesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_payments_f4730bcfc665723c, []int{4}
}
func (m *AdjustPricesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AdjustPricesRequest.Unmarshal(m, b)
}
func (m *AdjustPricesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AdjustPricesRequest.Marshal(b, m, deterministic)
}
func (dst *AdjustPricesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdjustPricesRequest.Merge(dst, src)
}
func (m *AdjustPricesRequest) XXX_Size() int {
	return xxx_messageInfo_AdjustPricesRequest.Size(m)
}
func (m *AdjustPricesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AdjustPricesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AdjustPricesRequest proto.InternalMessageInfo

func (m *AdjustPricesRequest) GetBandwidth() int64 {
	if m != nil {
		return m.Bandwidth
	}
	return 0
}

func (m *AdjustPricesRequest) GetStorage() int64 {
	if m != nil {
		return m.Storage
	}
	return 0
}

// The response message from adjusting cost basis on satelittes.
type AdjustPricesResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AdjustPricesResponse) Reset()         { *m = AdjustPricesResponse{} }
func (m *AdjustPricesResponse) String() string { return proto.CompactTextString(m) }
func (*AdjustPricesResponse) ProtoMessage()    {}
func (*AdjustPricesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_payments_f4730bcfc665723c, []int{5}
}
func (m *AdjustPricesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AdjustPricesResponse.Unmarshal(m, b)
}
func (m *AdjustPricesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AdjustPricesResponse.Marshal(b, m, deterministic)
}
func (dst *AdjustPricesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdjustPricesResponse.Merge(dst, src)
}
func (m *AdjustPricesResponse) XXX_Size() int {
	return xxx_messageInfo_AdjustPricesResponse.Size(m)
}
func (m *AdjustPricesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AdjustPricesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AdjustPricesResponse proto.InternalMessageInfo

// The request message for querying the data needed to generate a payments CSV
type GenerateCSVRequest struct {
	StartTime            *timestamp.Timestamp `protobuf:"bytes,1,opt,name=start_time,json=startTime" json:"start_time,omitempty"`
	EndTime              *timestamp.Timestamp `protobuf:"bytes,2,opt,name=end_time,json=endTime" json:"end_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *GenerateCSVRequest) Reset()         { *m = GenerateCSVRequest{} }
func (m *GenerateCSVRequest) String() string { return proto.CompactTextString(m) }
func (*GenerateCSVRequest) ProtoMessage()    {}
func (*GenerateCSVRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_payments_f4730bcfc665723c, []int{6}
}
func (m *GenerateCSVRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenerateCSVRequest.Unmarshal(m, b)
}
func (m *GenerateCSVRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenerateCSVRequest.Marshal(b, m, deterministic)
}
func (dst *GenerateCSVRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenerateCSVRequest.Merge(dst, src)
}
func (m *GenerateCSVRequest) XXX_Size() int {
	return xxx_messageInfo_GenerateCSVRequest.Size(m)
}
func (m *GenerateCSVRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GenerateCSVRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GenerateCSVRequest proto.InternalMessageInfo

func (m *GenerateCSVRequest) GetStartTime() *timestamp.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *GenerateCSVRequest) GetEndTime() *timestamp.Timestamp {
	if m != nil {
		return m.EndTime
	}
	return nil
}

// The response message for querying the data needed to generate a payments CSV
type GenerateCSVResponse struct {
	Filepath             string   `protobuf:"bytes,1,opt,name=filepath,proto3" json:"filepath,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenerateCSVResponse) Reset()         { *m = GenerateCSVResponse{} }
func (m *GenerateCSVResponse) String() string { return proto.CompactTextString(m) }
func (*GenerateCSVResponse) ProtoMessage()    {}
func (*GenerateCSVResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_payments_f4730bcfc665723c, []int{7}
}
func (m *GenerateCSVResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenerateCSVResponse.Unmarshal(m, b)
}
func (m *GenerateCSVResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenerateCSVResponse.Marshal(b, m, deterministic)
}
func (dst *GenerateCSVResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenerateCSVResponse.Merge(dst, src)
}
func (m *GenerateCSVResponse) XXX_Size() int {
	return xxx_messageInfo_GenerateCSVResponse.Size(m)
}
func (m *GenerateCSVResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GenerateCSVResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GenerateCSVResponse proto.InternalMessageInfo

func (m *GenerateCSVResponse) GetFilepath() string {
	if m != nil {
		return m.Filepath
	}
	return ""
}

func init() {
	proto.RegisterType((*PaymentRequest)(nil), "PaymentRequest")
	proto.RegisterType((*PaymentResponse)(nil), "PaymentResponse")
	proto.RegisterType((*CalculateRequest)(nil), "CalculateRequest")
	proto.RegisterType((*CalculateResponse)(nil), "CalculateResponse")
	proto.RegisterType((*AdjustPricesRequest)(nil), "AdjustPricesRequest")
	proto.RegisterType((*AdjustPricesResponse)(nil), "AdjustPricesResponse")
	proto.RegisterType((*GenerateCSVRequest)(nil), "GenerateCSVRequest")
	proto.RegisterType((*GenerateCSVResponse)(nil), "GenerateCSVResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PaymentsClient is the client API for Payments service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PaymentsClient interface {
	// Pay creates a payment to a single storage node
	Pay(ctx context.Context, in *PaymentRequest, opts ...grpc.CallOption) (*PaymentResponse, error)
	// Calculate determines the outstanding balance for a given storage node
	Calculate(ctx context.Context, in *CalculateRequest, opts ...grpc.CallOption) (*CalculateResponse, error)
	// AdjustPrices sets the prices paid by a satellite for data at rest and bandwidth
	AdjustPrices(ctx context.Context, in *AdjustPricesRequest, opts ...grpc.CallOption) (*AdjustPricesResponse, error)
	// GenerateCSV creates a csv file for payment purposes
	GenerateCSV(ctx context.Context, in *GenerateCSVRequest, opts ...grpc.CallOption) (*GenerateCSVResponse, error)
}

type paymentsClient struct {
	cc *grpc.ClientConn
}

func NewPaymentsClient(cc *grpc.ClientConn) PaymentsClient {
	return &paymentsClient{cc}
}

func (c *paymentsClient) Pay(ctx context.Context, in *PaymentRequest, opts ...grpc.CallOption) (*PaymentResponse, error) {
	out := new(PaymentResponse)
	err := c.cc.Invoke(ctx, "/Payments/Pay", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentsClient) Calculate(ctx context.Context, in *CalculateRequest, opts ...grpc.CallOption) (*CalculateResponse, error) {
	out := new(CalculateResponse)
	err := c.cc.Invoke(ctx, "/Payments/Calculate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentsClient) AdjustPrices(ctx context.Context, in *AdjustPricesRequest, opts ...grpc.CallOption) (*AdjustPricesResponse, error) {
	out := new(AdjustPricesResponse)
	err := c.cc.Invoke(ctx, "/Payments/AdjustPrices", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentsClient) GenerateCSV(ctx context.Context, in *GenerateCSVRequest, opts ...grpc.CallOption) (*GenerateCSVResponse, error) {
	out := new(GenerateCSVResponse)
	err := c.cc.Invoke(ctx, "/Payments/GenerateCSV", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentsServer is the server API for Payments service.
type PaymentsServer interface {
	// Pay creates a payment to a single storage node
	Pay(context.Context, *PaymentRequest) (*PaymentResponse, error)
	// Calculate determines the outstanding balance for a given storage node
	Calculate(context.Context, *CalculateRequest) (*CalculateResponse, error)
	// AdjustPrices sets the prices paid by a satellite for data at rest and bandwidth
	AdjustPrices(context.Context, *AdjustPricesRequest) (*AdjustPricesResponse, error)
	// GenerateCSV creates a csv file for payment purposes
	GenerateCSV(context.Context, *GenerateCSVRequest) (*GenerateCSVResponse, error)
}

func RegisterPaymentsServer(s *grpc.Server, srv PaymentsServer) {
	s.RegisterService(&_Payments_serviceDesc, srv)
}

func _Payments_Pay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServer).Pay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Payments/Pay",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServer).Pay(ctx, req.(*PaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payments_Calculate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CalculateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServer).Calculate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Payments/Calculate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServer).Calculate(ctx, req.(*CalculateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payments_AdjustPrices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdjustPricesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServer).AdjustPrices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Payments/AdjustPrices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServer).AdjustPrices(ctx, req.(*AdjustPricesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payments_GenerateCSV_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateCSVRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServer).GenerateCSV(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Payments/GenerateCSV",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServer).GenerateCSV(ctx, req.(*GenerateCSVRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Payments_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Payments",
	HandlerType: (*PaymentsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Pay",
			Handler:    _Payments_Pay_Handler,
		},
		{
			MethodName: "Calculate",
			Handler:    _Payments_Calculate_Handler,
		},
		{
			MethodName: "AdjustPrices",
			Handler:    _Payments_AdjustPrices_Handler,
		},
		{
			MethodName: "GenerateCSV",
			Handler:    _Payments_GenerateCSV_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "payments.proto",
}

func init() { proto.RegisterFile("payments.proto", fileDescriptor_payments_f4730bcfc665723c) }

var fileDescriptor_payments_f4730bcfc665723c = []byte{
	// 376 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0xcd, 0x4a, 0xfb, 0x40,
	0x14, 0xc5, 0x49, 0xfb, 0xff, 0xb7, 0xcd, 0xad, 0xf4, 0x63, 0x1a, 0xb5, 0x04, 0xc1, 0x92, 0x55,
	0x45, 0x98, 0x62, 0x45, 0x50, 0x5c, 0xd9, 0x2e, 0xc4, 0x85, 0x50, 0xa2, 0xb8, 0x70, 0x53, 0x26,
	0x9d, 0xdb, 0x1a, 0x49, 0x33, 0x31, 0x33, 0x41, 0xfa, 0x02, 0xbe, 0xa6, 0xaf, 0x22, 0xf9, 0xaa,
	0xfd, 0x92, 0x2e, 0xcf, 0x70, 0xee, 0xbd, 0x73, 0x7e, 0x07, 0x6a, 0x01, 0x5b, 0xcc, 0xd1, 0x57,
	0x92, 0x06, 0xa1, 0x50, 0xc2, 0x3c, 0x9d, 0x09, 0x31, 0xf3, 0xb0, 0x97, 0x28, 0x27, 0x9a, 0xf6,
	0x94, 0x3b, 0x47, 0xa9, 0xd8, 0x3c, 0x48, 0x0d, 0xd6, 0x19, 0xd4, 0x46, 0xe9, 0x88, 0x8d, 0x1f,
	0x11, 0x4a, 0x45, 0x8e, 0xa1, 0xec, 0x0b, 0x8e, 0x63, 0x97, 0xb7, 0xb5, 0x8e, 0xd6, 0xd5, 0xed,
	0x52, 0x2c, 0x1f, 0xb8, 0xd5, 0x84, 0xfa, 0xd2, 0x2a, 0x03, 0xe1, 0x4b, 0xb4, 0xce, 0xa1, 0x31,
	0x64, 0xde, 0x24, 0xf2, 0x98, 0xc2, 0xbd, 0xf3, 0x03, 0x68, 0xae, 0x98, 0xd3, 0x0d, 0x7f, 0xba,
	0x89, 0x01, 0xff, 0x95, 0x50, 0xcc, 0x6b, 0x17, 0x3a, 0x5a, 0xb7, 0x68, 0xa7, 0xc2, 0x7a, 0x84,
	0xd6, 0x1d, 0x7f, 0x8f, 0xa4, 0x1a, 0x85, 0xee, 0x04, 0x65, 0x7e, 0xf3, 0x04, 0x74, 0x87, 0xf9,
	0xfc, 0xd3, 0xe5, 0xea, 0x2d, 0xd9, 0x53, 0xb4, 0x7f, 0x1f, 0x48, 0x1b, 0xca, 0x52, 0x89, 0x90,
	0xcd, 0x30, 0x5b, 0x96, 0x4b, 0xeb, 0x08, 0x8c, 0xf5, 0x75, 0x59, 0xae, 0x2f, 0x0d, 0xc8, 0x3d,
	0xfa, 0x18, 0x32, 0x85, 0xc3, 0xa7, 0x97, 0xfc, 0xcc, 0x0d, 0x80, 0x54, 0x2c, 0x54, 0xe3, 0x98,
	0x62, 0x72, 0xa7, 0xda, 0x37, 0x69, 0x8a, 0x98, 0xe6, 0x88, 0xe9, 0x73, 0x8e, 0xd8, 0xd6, 0x13,
	0x77, 0xac, 0xc9, 0x15, 0x54, 0xd0, 0xe7, 0xe9, 0x60, 0x61, 0xef, 0x60, 0x19, 0x7d, 0x1e, 0x2b,
	0xeb, 0x02, 0x5a, 0x6b, 0xff, 0xc8, 0xa8, 0x99, 0x50, 0x99, 0xba, 0x1e, 0x06, 0x2c, 0x8b, 0xab,
	0xdb, 0x4b, 0xdd, 0xff, 0xd6, 0xa0, 0x92, 0xf5, 0x24, 0x49, 0x17, 0x8a, 0x23, 0xb6, 0x20, 0x75,
	0xba, 0x5e, 0xb2, 0xd9, 0xa0, 0x1b, 0x55, 0x92, 0x3e, 0xe8, 0xcb, 0x76, 0x48, 0x93, 0x6e, 0xd6,
	0x6a, 0x12, 0xba, 0x5d, 0xde, 0x2d, 0x1c, 0xac, 0xe2, 0x23, 0x06, 0xdd, 0x51, 0x8e, 0x79, 0x48,
	0x77, 0x31, 0x26, 0xd7, 0x50, 0x5d, 0x89, 0x46, 0x5a, 0x74, 0x1b, 0xb8, 0x69, 0xd0, 0x1d, 0xe9,
	0x07, 0xff, 0x5e, 0x0b, 0x81, 0xe3, 0x94, 0x12, 0x6e, 0x97, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff,
	0x52, 0x02, 0x1c, 0x61, 0xf3, 0x02, 0x00, 0x00,
}
