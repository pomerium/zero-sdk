// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: connect/connect.proto

package connect

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
	Connect_Subscribe_FullMethodName = "/pomerium.zero.Connect/Subscribe"
)

// ConnectClient is the client API for Connect service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConnectClient interface {
	// Subscribe is used to send a stream of messages from the Zero Cloud to the Pomerium Core in managed mode.
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (Connect_SubscribeClient, error)
}

type connectClient struct {
	cc grpc.ClientConnInterface
}

func NewConnectClient(cc grpc.ClientConnInterface) ConnectClient {
	return &connectClient{cc}
}

func (c *connectClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (Connect_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &Connect_ServiceDesc.Streams[0], Connect_Subscribe_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &connectSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Connect_SubscribeClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type connectSubscribeClient struct {
	grpc.ClientStream
}

func (x *connectSubscribeClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ConnectServer is the server API for Connect service.
// All implementations should embed UnimplementedConnectServer
// for forward compatibility
type ConnectServer interface {
	// Subscribe is used to send a stream of messages from the Zero Cloud to the Pomerium Core in managed mode.
	Subscribe(*SubscribeRequest, Connect_SubscribeServer) error
}

// UnimplementedConnectServer should be embedded to have forward compatible implementations.
type UnimplementedConnectServer struct {
}

func (UnimplementedConnectServer) Subscribe(*SubscribeRequest, Connect_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}

// UnsafeConnectServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConnectServer will
// result in compilation errors.
type UnsafeConnectServer interface {
	mustEmbedUnimplementedConnectServer()
}

func RegisterConnectServer(s grpc.ServiceRegistrar, srv ConnectServer) {
	s.RegisterService(&Connect_ServiceDesc, srv)
}

func _Connect_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConnectServer).Subscribe(m, &connectSubscribeServer{stream})
}

type Connect_SubscribeServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type connectSubscribeServer struct {
	grpc.ServerStream
}

func (x *connectSubscribeServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

// Connect_ServiceDesc is the grpc.ServiceDesc for Connect service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Connect_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pomerium.zero.Connect",
	HandlerType: (*ConnectServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _Connect_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "connect/connect.proto",
}
