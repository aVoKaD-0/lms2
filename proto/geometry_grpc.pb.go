// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.11
// source: proto/geometry.proto

package proto

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

// KalkulatorServiceClient is the client API for KalkulatorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KalkulatorServiceClient interface {
	Reception(ctx context.Context, in *ExpressionRequest, opts ...grpc.CallOption) (*ExpressionResponse, error)
}

type kalkulatorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKalkulatorServiceClient(cc grpc.ClientConnInterface) KalkulatorServiceClient {
	return &kalkulatorServiceClient{cc}
}

func (c *kalkulatorServiceClient) Reception(ctx context.Context, in *ExpressionRequest, opts ...grpc.CallOption) (*ExpressionResponse, error) {
	out := new(ExpressionResponse)
	err := c.cc.Invoke(ctx, "/geometry.KalkulatorService/Reception", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KalkulatorServiceServer is the server API for KalkulatorService service.
// All implementations must embed UnimplementedKalkulatorServiceServer
// for forward compatibility
type KalkulatorServiceServer interface {
	Reception(context.Context, *ExpressionRequest) (*ExpressionResponse, error)
	mustEmbedUnimplementedKalkulatorServiceServer()
}

// UnimplementedKalkulatorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedKalkulatorServiceServer struct {
}

func (UnimplementedKalkulatorServiceServer) Reception(context.Context, *ExpressionRequest) (*ExpressionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reception not implemented")
}
func (UnimplementedKalkulatorServiceServer) mustEmbedUnimplementedKalkulatorServiceServer() {}

// UnsafeKalkulatorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KalkulatorServiceServer will
// result in compilation errors.
type UnsafeKalkulatorServiceServer interface {
	mustEmbedUnimplementedKalkulatorServiceServer()
}

func RegisterKalkulatorServiceServer(s grpc.ServiceRegistrar, srv KalkulatorServiceServer) {
	s.RegisterService(&KalkulatorService_ServiceDesc, srv)
}

func _KalkulatorService_Reception_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExpressionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KalkulatorServiceServer).Reception(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/geometry.KalkulatorService/Reception",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KalkulatorServiceServer).Reception(ctx, req.(*ExpressionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KalkulatorService_ServiceDesc is the grpc.ServiceDesc for KalkulatorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KalkulatorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "geometry.KalkulatorService",
	HandlerType: (*KalkulatorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Reception",
			Handler:    _KalkulatorService_Reception_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/geometry.proto",
}
