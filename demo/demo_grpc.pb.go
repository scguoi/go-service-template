// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: demo.proto

package demo

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	DemoService_OneWay_FullMethodName = "/example.DemoService/OneWay"
	DemoService_Stream_FullMethodName = "/example.DemoService/Stream"
)

// DemoServiceClient is the client API for DemoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DemoServiceClient interface {
	// 单次调用的方法 支持http和grpc调用
	OneWay(ctx context.Context, in *ReqPkg, opts ...grpc.CallOption) (*RespPkg, error)
	// 流式调用的方法 支持http chunked和grpc调用
	Stream(ctx context.Context, opts ...grpc.CallOption) (DemoService_StreamClient, error)
}

type demoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDemoServiceClient(cc grpc.ClientConnInterface) DemoServiceClient {
	return &demoServiceClient{cc}
}

func (c *demoServiceClient) OneWay(ctx context.Context, in *ReqPkg, opts ...grpc.CallOption) (*RespPkg, error) {
	out := new(RespPkg)
	err := c.cc.Invoke(ctx, DemoService_OneWay_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *demoServiceClient) Stream(ctx context.Context, opts ...grpc.CallOption) (DemoService_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &DemoService_ServiceDesc.Streams[0], DemoService_Stream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &demoServiceStreamClient{stream}
	return x, nil
}

type DemoService_StreamClient interface {
	Send(*ReqPkg) error
	Recv() (*RespPkg, error)
	grpc.ClientStream
}

type demoServiceStreamClient struct {
	grpc.ClientStream
}

func (x *demoServiceStreamClient) Send(m *ReqPkg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *demoServiceStreamClient) Recv() (*RespPkg, error) {
	m := new(RespPkg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DemoServiceServer is the server API for DemoService service.
// All implementations must embed UnimplementedDemoServiceServer
// for forward compatibility
type DemoServiceServer interface {
	// 单次调用的方法 支持http和grpc调用
	OneWay(context.Context, *ReqPkg) (*RespPkg, error)
	// 流式调用的方法 支持http chunked和grpc调用
	Stream(DemoService_StreamServer) error
	mustEmbedUnimplementedDemoServiceServer()
}

// UnimplementedDemoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDemoServiceServer struct {
}

func (UnimplementedDemoServiceServer) OneWay(context.Context, *ReqPkg) (*RespPkg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OneWay not implemented")
}
func (UnimplementedDemoServiceServer) Stream(DemoService_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedDemoServiceServer) mustEmbedUnimplementedDemoServiceServer() {}

// UnsafeDemoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DemoServiceServer will
// result in compilation errors.
type UnsafeDemoServiceServer interface {
	mustEmbedUnimplementedDemoServiceServer()
}

func RegisterDemoServiceServer(s grpc.ServiceRegistrar, srv DemoServiceServer) {
	s.RegisterService(&DemoService_ServiceDesc, srv)
}

func _DemoService_OneWay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqPkg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DemoServiceServer).OneWay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DemoService_OneWay_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DemoServiceServer).OneWay(ctx, req.(*ReqPkg))
	}
	return interceptor(ctx, in, info, handler)
}

func _DemoService_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DemoServiceServer).Stream(&demoServiceStreamServer{stream})
}

type DemoService_StreamServer interface {
	Send(*RespPkg) error
	Recv() (*ReqPkg, error)
	grpc.ServerStream
}

type demoServiceStreamServer struct {
	grpc.ServerStream
}

func (x *demoServiceStreamServer) Send(m *RespPkg) error {
	return x.ServerStream.SendMsg(m)
}

func (x *demoServiceStreamServer) Recv() (*ReqPkg, error) {
	m := new(ReqPkg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DemoService_ServiceDesc is the grpc.ServiceDesc for DemoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DemoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "example.DemoService",
	HandlerType: (*DemoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OneWay",
			Handler:    _DemoService_OneWay_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _DemoService_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "demo.proto",
}
