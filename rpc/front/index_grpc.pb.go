// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: proto/front/index.proto

package front

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

// FrontServiceClient is the client API for FrontService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FrontServiceClient interface {
	Health(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthResponse, error)
	Sample(ctx context.Context, in *SampleRequest, opts ...grpc.CallOption) (*SampleResponse, error)
}

type frontServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFrontServiceClient(cc grpc.ClientConnInterface) FrontServiceClient {
	return &frontServiceClient{cc}
}

func (c *frontServiceClient) Health(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthResponse, error) {
	out := new(HealthResponse)
	err := c.cc.Invoke(ctx, "/front.FrontService/Health", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *frontServiceClient) Sample(ctx context.Context, in *SampleRequest, opts ...grpc.CallOption) (*SampleResponse, error) {
	out := new(SampleResponse)
	err := c.cc.Invoke(ctx, "/front.FrontService/Sample", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FrontServiceServer is the server API for FrontService service.
// All implementations must embed UnimplementedFrontServiceServer
// for forward compatibility
type FrontServiceServer interface {
	Health(context.Context, *HealthRequest) (*HealthResponse, error)
	Sample(context.Context, *SampleRequest) (*SampleResponse, error)
	mustEmbedUnimplementedFrontServiceServer()
}

// UnimplementedFrontServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFrontServiceServer struct {
}

func (UnimplementedFrontServiceServer) Health(context.Context, *HealthRequest) (*HealthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Health not implemented")
}
func (UnimplementedFrontServiceServer) Sample(context.Context, *SampleRequest) (*SampleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sample not implemented")
}
func (UnimplementedFrontServiceServer) mustEmbedUnimplementedFrontServiceServer() {}

// UnsafeFrontServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FrontServiceServer will
// result in compilation errors.
type UnsafeFrontServiceServer interface {
	mustEmbedUnimplementedFrontServiceServer()
}

func RegisterFrontServiceServer(s grpc.ServiceRegistrar, srv FrontServiceServer) {
	s.RegisterService(&FrontService_ServiceDesc, srv)
}

func _FrontService_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrontServiceServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/front.FrontService/Health",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrontServiceServer).Health(ctx, req.(*HealthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FrontService_Sample_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SampleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrontServiceServer).Sample(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/front.FrontService/Sample",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrontServiceServer).Sample(ctx, req.(*SampleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FrontService_ServiceDesc is the grpc.ServiceDesc for FrontService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FrontService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "front.FrontService",
	HandlerType: (*FrontServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Health",
			Handler:    _FrontService_Health_Handler,
		},
		{
			MethodName: "Sample",
			Handler:    _FrontService_Sample_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/front/index.proto",
}