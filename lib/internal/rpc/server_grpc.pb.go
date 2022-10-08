// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: server.proto

package rpc

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

// FsSnapshotClient is the client API for FsSnapshot service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FsSnapshotClient interface {
	CanCreateSnapshots(ctx context.Context, in *CanCreateSnapshotsRequest, opts ...grpc.CallOption) (*CanCreateSnapshotsReply, error)
	ListProviders(ctx context.Context, in *ListProvidersRequest, opts ...grpc.CallOption) (*ListProvidersReply, error)
	ListSets(ctx context.Context, in *ListSetsRequest, opts ...grpc.CallOption) (*ListSetsReply, error)
	ListSnapshots(ctx context.Context, in *ListSnapshotsRequest, opts ...grpc.CallOption) (*ListSnapshotsReply, error)
	SimplifyId(ctx context.Context, in *SimplifyIdRequest, opts ...grpc.CallOption) (*SimplifyIdReply, error)
	DeleteSet(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error)
	DeleteSnapshot(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error)
	ListMountPoints(ctx context.Context, in *ListMountPointsRequest, opts ...grpc.CallOption) (*ListMountPointsReply, error)
	StartBackup(ctx context.Context, in *StartBackupRequest, opts ...grpc.CallOption) (FsSnapshot_StartBackupClient, error)
	TryToCreateTemporarySnapshot(ctx context.Context, in *TryToCreateTemporarySnapshotRequest, opts ...grpc.CallOption) (FsSnapshot_TryToCreateTemporarySnapshotClient, error)
	CloseBackup(ctx context.Context, in *CloseBackupRequest, opts ...grpc.CallOption) (FsSnapshot_CloseBackupClient, error)
}

type fsSnapshotClient struct {
	cc grpc.ClientConnInterface
}

func NewFsSnapshotClient(cc grpc.ClientConnInterface) FsSnapshotClient {
	return &fsSnapshotClient{cc}
}

func (c *fsSnapshotClient) CanCreateSnapshots(ctx context.Context, in *CanCreateSnapshotsRequest, opts ...grpc.CallOption) (*CanCreateSnapshotsReply, error) {
	out := new(CanCreateSnapshotsReply)
	err := c.cc.Invoke(ctx, "/rpc.FsSnapshot/CanCreateSnapshots", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fsSnapshotClient) ListProviders(ctx context.Context, in *ListProvidersRequest, opts ...grpc.CallOption) (*ListProvidersReply, error) {
	out := new(ListProvidersReply)
	err := c.cc.Invoke(ctx, "/rpc.FsSnapshot/ListProviders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fsSnapshotClient) ListSets(ctx context.Context, in *ListSetsRequest, opts ...grpc.CallOption) (*ListSetsReply, error) {
	out := new(ListSetsReply)
	err := c.cc.Invoke(ctx, "/rpc.FsSnapshot/ListSets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fsSnapshotClient) ListSnapshots(ctx context.Context, in *ListSnapshotsRequest, opts ...grpc.CallOption) (*ListSnapshotsReply, error) {
	out := new(ListSnapshotsReply)
	err := c.cc.Invoke(ctx, "/rpc.FsSnapshot/ListSnapshots", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fsSnapshotClient) SimplifyId(ctx context.Context, in *SimplifyIdRequest, opts ...grpc.CallOption) (*SimplifyIdReply, error) {
	out := new(SimplifyIdReply)
	err := c.cc.Invoke(ctx, "/rpc.FsSnapshot/SimplifyId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fsSnapshotClient) DeleteSet(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error) {
	out := new(DeleteReply)
	err := c.cc.Invoke(ctx, "/rpc.FsSnapshot/DeleteSet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fsSnapshotClient) DeleteSnapshot(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error) {
	out := new(DeleteReply)
	err := c.cc.Invoke(ctx, "/rpc.FsSnapshot/DeleteSnapshot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fsSnapshotClient) ListMountPoints(ctx context.Context, in *ListMountPointsRequest, opts ...grpc.CallOption) (*ListMountPointsReply, error) {
	out := new(ListMountPointsReply)
	err := c.cc.Invoke(ctx, "/rpc.FsSnapshot/ListMountPoints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fsSnapshotClient) StartBackup(ctx context.Context, in *StartBackupRequest, opts ...grpc.CallOption) (FsSnapshot_StartBackupClient, error) {
	stream, err := c.cc.NewStream(ctx, &FsSnapshot_ServiceDesc.Streams[0], "/rpc.FsSnapshot/StartBackup", opts...)
	if err != nil {
		return nil, err
	}
	x := &fsSnapshotStartBackupClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FsSnapshot_StartBackupClient interface {
	Recv() (*StartBackupReply, error)
	grpc.ClientStream
}

type fsSnapshotStartBackupClient struct {
	grpc.ClientStream
}

func (x *fsSnapshotStartBackupClient) Recv() (*StartBackupReply, error) {
	m := new(StartBackupReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fsSnapshotClient) TryToCreateTemporarySnapshot(ctx context.Context, in *TryToCreateTemporarySnapshotRequest, opts ...grpc.CallOption) (FsSnapshot_TryToCreateTemporarySnapshotClient, error) {
	stream, err := c.cc.NewStream(ctx, &FsSnapshot_ServiceDesc.Streams[1], "/rpc.FsSnapshot/TryToCreateTemporarySnapshot", opts...)
	if err != nil {
		return nil, err
	}
	x := &fsSnapshotTryToCreateTemporarySnapshotClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FsSnapshot_TryToCreateTemporarySnapshotClient interface {
	Recv() (*TryToCreateTemporarySnapshotReply, error)
	grpc.ClientStream
}

type fsSnapshotTryToCreateTemporarySnapshotClient struct {
	grpc.ClientStream
}

func (x *fsSnapshotTryToCreateTemporarySnapshotClient) Recv() (*TryToCreateTemporarySnapshotReply, error) {
	m := new(TryToCreateTemporarySnapshotReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fsSnapshotClient) CloseBackup(ctx context.Context, in *CloseBackupRequest, opts ...grpc.CallOption) (FsSnapshot_CloseBackupClient, error) {
	stream, err := c.cc.NewStream(ctx, &FsSnapshot_ServiceDesc.Streams[2], "/rpc.FsSnapshot/CloseBackup", opts...)
	if err != nil {
		return nil, err
	}
	x := &fsSnapshotCloseBackupClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FsSnapshot_CloseBackupClient interface {
	Recv() (*CloseBackupReply, error)
	grpc.ClientStream
}

type fsSnapshotCloseBackupClient struct {
	grpc.ClientStream
}

func (x *fsSnapshotCloseBackupClient) Recv() (*CloseBackupReply, error) {
	m := new(CloseBackupReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FsSnapshotServer is the server API for FsSnapshot service.
// All implementations must embed UnimplementedFsSnapshotServer
// for forward compatibility
type FsSnapshotServer interface {
	CanCreateSnapshots(context.Context, *CanCreateSnapshotsRequest) (*CanCreateSnapshotsReply, error)
	ListProviders(context.Context, *ListProvidersRequest) (*ListProvidersReply, error)
	ListSets(context.Context, *ListSetsRequest) (*ListSetsReply, error)
	ListSnapshots(context.Context, *ListSnapshotsRequest) (*ListSnapshotsReply, error)
	SimplifyId(context.Context, *SimplifyIdRequest) (*SimplifyIdReply, error)
	DeleteSet(context.Context, *DeleteRequest) (*DeleteReply, error)
	DeleteSnapshot(context.Context, *DeleteRequest) (*DeleteReply, error)
	ListMountPoints(context.Context, *ListMountPointsRequest) (*ListMountPointsReply, error)
	StartBackup(*StartBackupRequest, FsSnapshot_StartBackupServer) error
	TryToCreateTemporarySnapshot(*TryToCreateTemporarySnapshotRequest, FsSnapshot_TryToCreateTemporarySnapshotServer) error
	CloseBackup(*CloseBackupRequest, FsSnapshot_CloseBackupServer) error
	mustEmbedUnimplementedFsSnapshotServer()
}

// UnimplementedFsSnapshotServer must be embedded to have forward compatible implementations.
type UnimplementedFsSnapshotServer struct {
}

func (UnimplementedFsSnapshotServer) CanCreateSnapshots(context.Context, *CanCreateSnapshotsRequest) (*CanCreateSnapshotsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CanCreateSnapshots not implemented")
}
func (UnimplementedFsSnapshotServer) ListProviders(context.Context, *ListProvidersRequest) (*ListProvidersReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProviders not implemented")
}
func (UnimplementedFsSnapshotServer) ListSets(context.Context, *ListSetsRequest) (*ListSetsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSets not implemented")
}
func (UnimplementedFsSnapshotServer) ListSnapshots(context.Context, *ListSnapshotsRequest) (*ListSnapshotsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSnapshots not implemented")
}
func (UnimplementedFsSnapshotServer) SimplifyId(context.Context, *SimplifyIdRequest) (*SimplifyIdReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SimplifyId not implemented")
}
func (UnimplementedFsSnapshotServer) DeleteSet(context.Context, *DeleteRequest) (*DeleteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSet not implemented")
}
func (UnimplementedFsSnapshotServer) DeleteSnapshot(context.Context, *DeleteRequest) (*DeleteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSnapshot not implemented")
}
func (UnimplementedFsSnapshotServer) ListMountPoints(context.Context, *ListMountPointsRequest) (*ListMountPointsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMountPoints not implemented")
}
func (UnimplementedFsSnapshotServer) StartBackup(*StartBackupRequest, FsSnapshot_StartBackupServer) error {
	return status.Errorf(codes.Unimplemented, "method StartBackup not implemented")
}
func (UnimplementedFsSnapshotServer) TryToCreateTemporarySnapshot(*TryToCreateTemporarySnapshotRequest, FsSnapshot_TryToCreateTemporarySnapshotServer) error {
	return status.Errorf(codes.Unimplemented, "method TryToCreateTemporarySnapshot not implemented")
}
func (UnimplementedFsSnapshotServer) CloseBackup(*CloseBackupRequest, FsSnapshot_CloseBackupServer) error {
	return status.Errorf(codes.Unimplemented, "method CloseBackup not implemented")
}
func (UnimplementedFsSnapshotServer) mustEmbedUnimplementedFsSnapshotServer() {}

// UnsafeFsSnapshotServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FsSnapshotServer will
// result in compilation errors.
type UnsafeFsSnapshotServer interface {
	mustEmbedUnimplementedFsSnapshotServer()
}

func RegisterFsSnapshotServer(s grpc.ServiceRegistrar, srv FsSnapshotServer) {
	s.RegisterService(&FsSnapshot_ServiceDesc, srv)
}

func _FsSnapshot_CanCreateSnapshots_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CanCreateSnapshotsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FsSnapshotServer).CanCreateSnapshots(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.FsSnapshot/CanCreateSnapshots",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FsSnapshotServer).CanCreateSnapshots(ctx, req.(*CanCreateSnapshotsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FsSnapshot_ListProviders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProvidersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FsSnapshotServer).ListProviders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.FsSnapshot/ListProviders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FsSnapshotServer).ListProviders(ctx, req.(*ListProvidersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FsSnapshot_ListSets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FsSnapshotServer).ListSets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.FsSnapshot/ListSets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FsSnapshotServer).ListSets(ctx, req.(*ListSetsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FsSnapshot_ListSnapshots_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSnapshotsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FsSnapshotServer).ListSnapshots(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.FsSnapshot/ListSnapshots",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FsSnapshotServer).ListSnapshots(ctx, req.(*ListSnapshotsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FsSnapshot_SimplifyId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SimplifyIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FsSnapshotServer).SimplifyId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.FsSnapshot/SimplifyId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FsSnapshotServer).SimplifyId(ctx, req.(*SimplifyIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FsSnapshot_DeleteSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FsSnapshotServer).DeleteSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.FsSnapshot/DeleteSet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FsSnapshotServer).DeleteSet(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FsSnapshot_DeleteSnapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FsSnapshotServer).DeleteSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.FsSnapshot/DeleteSnapshot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FsSnapshotServer).DeleteSnapshot(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FsSnapshot_ListMountPoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMountPointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FsSnapshotServer).ListMountPoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.FsSnapshot/ListMountPoints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FsSnapshotServer).ListMountPoints(ctx, req.(*ListMountPointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FsSnapshot_StartBackup_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StartBackupRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FsSnapshotServer).StartBackup(m, &fsSnapshotStartBackupServer{stream})
}

type FsSnapshot_StartBackupServer interface {
	Send(*StartBackupReply) error
	grpc.ServerStream
}

type fsSnapshotStartBackupServer struct {
	grpc.ServerStream
}

func (x *fsSnapshotStartBackupServer) Send(m *StartBackupReply) error {
	return x.ServerStream.SendMsg(m)
}

func _FsSnapshot_TryToCreateTemporarySnapshot_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TryToCreateTemporarySnapshotRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FsSnapshotServer).TryToCreateTemporarySnapshot(m, &fsSnapshotTryToCreateTemporarySnapshotServer{stream})
}

type FsSnapshot_TryToCreateTemporarySnapshotServer interface {
	Send(*TryToCreateTemporarySnapshotReply) error
	grpc.ServerStream
}

type fsSnapshotTryToCreateTemporarySnapshotServer struct {
	grpc.ServerStream
}

func (x *fsSnapshotTryToCreateTemporarySnapshotServer) Send(m *TryToCreateTemporarySnapshotReply) error {
	return x.ServerStream.SendMsg(m)
}

func _FsSnapshot_CloseBackup_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CloseBackupRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FsSnapshotServer).CloseBackup(m, &fsSnapshotCloseBackupServer{stream})
}

type FsSnapshot_CloseBackupServer interface {
	Send(*CloseBackupReply) error
	grpc.ServerStream
}

type fsSnapshotCloseBackupServer struct {
	grpc.ServerStream
}

func (x *fsSnapshotCloseBackupServer) Send(m *CloseBackupReply) error {
	return x.ServerStream.SendMsg(m)
}

// FsSnapshot_ServiceDesc is the grpc.ServiceDesc for FsSnapshot service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FsSnapshot_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.FsSnapshot",
	HandlerType: (*FsSnapshotServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CanCreateSnapshots",
			Handler:    _FsSnapshot_CanCreateSnapshots_Handler,
		},
		{
			MethodName: "ListProviders",
			Handler:    _FsSnapshot_ListProviders_Handler,
		},
		{
			MethodName: "ListSets",
			Handler:    _FsSnapshot_ListSets_Handler,
		},
		{
			MethodName: "ListSnapshots",
			Handler:    _FsSnapshot_ListSnapshots_Handler,
		},
		{
			MethodName: "SimplifyId",
			Handler:    _FsSnapshot_SimplifyId_Handler,
		},
		{
			MethodName: "DeleteSet",
			Handler:    _FsSnapshot_DeleteSet_Handler,
		},
		{
			MethodName: "DeleteSnapshot",
			Handler:    _FsSnapshot_DeleteSnapshot_Handler,
		},
		{
			MethodName: "ListMountPoints",
			Handler:    _FsSnapshot_ListMountPoints_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StartBackup",
			Handler:       _FsSnapshot_StartBackup_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "TryToCreateTemporarySnapshot",
			Handler:       _FsSnapshot_TryToCreateTemporarySnapshot_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "CloseBackup",
			Handler:       _FsSnapshot_CloseBackup_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "server.proto",
}
