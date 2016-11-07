// Code generated by protoc-gen-go.
// source: api/api.proto
// DO NOT EDIT!

/*
Package api is a generated protocol buffer package.

It is generated from these files:
	api/api.proto

It has these top-level messages:
	Client
	CreateClientReq
	CreateClientResp
	DeleteClientReq
	DeleteClientResp
	Password
	CreatePasswordReq
	CreatePasswordResp
	UpdatePasswordReq
	UpdatePasswordResp
	DeletePasswordReq
	DeletePasswordResp
	Version
	VersionReq
*/
package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Client represents an OAuth2 client.
type Client struct {
	Id           string   `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Secret       string   `protobuf:"bytes,2,opt,name=secret" json:"secret,omitempty"`
	RedirectUris []string `protobuf:"bytes,3,rep,name=redirect_uris,json=redirectUris" json:"redirect_uris,omitempty"`
	TrustedPeers []string `protobuf:"bytes,4,rep,name=trusted_peers,json=trustedPeers" json:"trusted_peers,omitempty"`
	Public       bool     `protobuf:"varint,5,opt,name=public" json:"public,omitempty"`
	Name         string   `protobuf:"bytes,6,opt,name=name" json:"name,omitempty"`
	LogoUrl      string   `protobuf:"bytes,7,opt,name=logo_url,json=logoUrl" json:"logo_url,omitempty"`
}

func (m *Client) Reset()                    { *m = Client{} }
func (m *Client) String() string            { return proto.CompactTextString(m) }
func (*Client) ProtoMessage()               {}
func (*Client) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// CreateClientReq is a request to make a client.
type CreateClientReq struct {
	Client *Client `protobuf:"bytes,1,opt,name=client" json:"client,omitempty"`
}

func (m *CreateClientReq) Reset()                    { *m = CreateClientReq{} }
func (m *CreateClientReq) String() string            { return proto.CompactTextString(m) }
func (*CreateClientReq) ProtoMessage()               {}
func (*CreateClientReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CreateClientReq) GetClient() *Client {
	if m != nil {
		return m.Client
	}
	return nil
}

// CreateClientResp returns the response from creating a client.
type CreateClientResp struct {
	AlreadyExists bool    `protobuf:"varint,1,opt,name=already_exists,json=alreadyExists" json:"already_exists,omitempty"`
	Client        *Client `protobuf:"bytes,2,opt,name=client" json:"client,omitempty"`
}

func (m *CreateClientResp) Reset()                    { *m = CreateClientResp{} }
func (m *CreateClientResp) String() string            { return proto.CompactTextString(m) }
func (*CreateClientResp) ProtoMessage()               {}
func (*CreateClientResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *CreateClientResp) GetClient() *Client {
	if m != nil {
		return m.Client
	}
	return nil
}

// DeleteClientReq is a request to delete a client.
type DeleteClientReq struct {
	// The ID of the client.
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *DeleteClientReq) Reset()                    { *m = DeleteClientReq{} }
func (m *DeleteClientReq) String() string            { return proto.CompactTextString(m) }
func (*DeleteClientReq) ProtoMessage()               {}
func (*DeleteClientReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

// DeleteClientResp determines if the.
type DeleteClientResp struct {
	NotFound bool `protobuf:"varint,1,opt,name=not_found,json=notFound" json:"not_found,omitempty"`
}

func (m *DeleteClientResp) Reset()                    { *m = DeleteClientResp{} }
func (m *DeleteClientResp) String() string            { return proto.CompactTextString(m) }
func (*DeleteClientResp) ProtoMessage()               {}
func (*DeleteClientResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

// Password is an email for password mapping managed by the storage.
type Password struct {
	Email string `protobuf:"bytes,1,opt,name=email" json:"email,omitempty"`
	// Currently we do not accept plain text passwords. Could be an option in the future.
	Hash     []byte `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	Username string `protobuf:"bytes,3,opt,name=username" json:"username,omitempty"`
	UserId   string `protobuf:"bytes,4,opt,name=user_id,json=userId" json:"user_id,omitempty"`
}

func (m *Password) Reset()                    { *m = Password{} }
func (m *Password) String() string            { return proto.CompactTextString(m) }
func (*Password) ProtoMessage()               {}
func (*Password) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

// CreatePasswordReq is a request to make a password.
type CreatePasswordReq struct {
	Password *Password `protobuf:"bytes,1,opt,name=password" json:"password,omitempty"`
}

func (m *CreatePasswordReq) Reset()                    { *m = CreatePasswordReq{} }
func (m *CreatePasswordReq) String() string            { return proto.CompactTextString(m) }
func (*CreatePasswordReq) ProtoMessage()               {}
func (*CreatePasswordReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *CreatePasswordReq) GetPassword() *Password {
	if m != nil {
		return m.Password
	}
	return nil
}

// CreatePasswordResp returns the response from creating a password.
type CreatePasswordResp struct {
	AlreadyExists bool `protobuf:"varint,1,opt,name=already_exists,json=alreadyExists" json:"already_exists,omitempty"`
}

func (m *CreatePasswordResp) Reset()                    { *m = CreatePasswordResp{} }
func (m *CreatePasswordResp) String() string            { return proto.CompactTextString(m) }
func (*CreatePasswordResp) ProtoMessage()               {}
func (*CreatePasswordResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

// UpdatePasswordReq is a request to modify an existing password.
type UpdatePasswordReq struct {
	// The email used to lookup the password. This field cannot be modified
	Email       string `protobuf:"bytes,1,opt,name=email" json:"email,omitempty"`
	NewHash     []byte `protobuf:"bytes,2,opt,name=new_hash,json=newHash,proto3" json:"new_hash,omitempty"`
	NewUsername string `protobuf:"bytes,3,opt,name=new_username,json=newUsername" json:"new_username,omitempty"`
}

func (m *UpdatePasswordReq) Reset()                    { *m = UpdatePasswordReq{} }
func (m *UpdatePasswordReq) String() string            { return proto.CompactTextString(m) }
func (*UpdatePasswordReq) ProtoMessage()               {}
func (*UpdatePasswordReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

// UpdatePasswordResp returns the response from modifying an existing password.
type UpdatePasswordResp struct {
	NotFound bool `protobuf:"varint,1,opt,name=not_found,json=notFound" json:"not_found,omitempty"`
}

func (m *UpdatePasswordResp) Reset()                    { *m = UpdatePasswordResp{} }
func (m *UpdatePasswordResp) String() string            { return proto.CompactTextString(m) }
func (*UpdatePasswordResp) ProtoMessage()               {}
func (*UpdatePasswordResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

// DeletePasswordReq is a request to delete a password.
type DeletePasswordReq struct {
	Email string `protobuf:"bytes,1,opt,name=email" json:"email,omitempty"`
}

func (m *DeletePasswordReq) Reset()                    { *m = DeletePasswordReq{} }
func (m *DeletePasswordReq) String() string            { return proto.CompactTextString(m) }
func (*DeletePasswordReq) ProtoMessage()               {}
func (*DeletePasswordReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

// DeletePasswordResp returns the response from deleting a password.
type DeletePasswordResp struct {
	NotFound bool `protobuf:"varint,1,opt,name=not_found,json=notFound" json:"not_found,omitempty"`
}

func (m *DeletePasswordResp) Reset()                    { *m = DeletePasswordResp{} }
func (m *DeletePasswordResp) String() string            { return proto.CompactTextString(m) }
func (*DeletePasswordResp) ProtoMessage()               {}
func (*DeletePasswordResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

// Version holds the version info of components.
type Version struct {
	// Semantic version of the server.
	Server string `protobuf:"bytes,1,opt,name=server" json:"server,omitempty"`
	// Numeric version of the API.
	Api int32 `protobuf:"varint,2,opt,name=api" json:"api,omitempty"`
}

func (m *Version) Reset()                    { *m = Version{} }
func (m *Version) String() string            { return proto.CompactTextString(m) }
func (*Version) ProtoMessage()               {}
func (*Version) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

// VersionReq is a request to fetch version info.
type VersionReq struct {
}

func (m *VersionReq) Reset()                    { *m = VersionReq{} }
func (m *VersionReq) String() string            { return proto.CompactTextString(m) }
func (*VersionReq) ProtoMessage()               {}
func (*VersionReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func init() {
	proto.RegisterType((*Client)(nil), "api.Client")
	proto.RegisterType((*CreateClientReq)(nil), "api.CreateClientReq")
	proto.RegisterType((*CreateClientResp)(nil), "api.CreateClientResp")
	proto.RegisterType((*DeleteClientReq)(nil), "api.DeleteClientReq")
	proto.RegisterType((*DeleteClientResp)(nil), "api.DeleteClientResp")
	proto.RegisterType((*Password)(nil), "api.Password")
	proto.RegisterType((*CreatePasswordReq)(nil), "api.CreatePasswordReq")
	proto.RegisterType((*CreatePasswordResp)(nil), "api.CreatePasswordResp")
	proto.RegisterType((*UpdatePasswordReq)(nil), "api.UpdatePasswordReq")
	proto.RegisterType((*UpdatePasswordResp)(nil), "api.UpdatePasswordResp")
	proto.RegisterType((*DeletePasswordReq)(nil), "api.DeletePasswordReq")
	proto.RegisterType((*DeletePasswordResp)(nil), "api.DeletePasswordResp")
	proto.RegisterType((*Version)(nil), "api.Version")
	proto.RegisterType((*VersionReq)(nil), "api.VersionReq")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Dex service

type DexClient interface {
	// CreateClient attempts to create the client.
	CreateClient(ctx context.Context, in *CreateClientReq, opts ...grpc.CallOption) (*CreateClientResp, error)
	// DeleteClient attempts to delete the provided client.
	DeleteClient(ctx context.Context, in *DeleteClientReq, opts ...grpc.CallOption) (*DeleteClientResp, error)
	// CreatePassword attempts to create the password.
	CreatePassword(ctx context.Context, in *CreatePasswordReq, opts ...grpc.CallOption) (*CreatePasswordResp, error)
	// UpdatePassword attempts to modify existing password.
	UpdatePassword(ctx context.Context, in *UpdatePasswordReq, opts ...grpc.CallOption) (*UpdatePasswordResp, error)
	// DeletePassword attempts to delete the password.
	DeletePassword(ctx context.Context, in *DeletePasswordReq, opts ...grpc.CallOption) (*DeletePasswordResp, error)
	// Version attempts to get the version info.
	GetVersion(ctx context.Context, in *VersionReq, opts ...grpc.CallOption) (*Version, error)
}

type dexClient struct {
	cc *grpc.ClientConn
}

func NewDexClient(cc *grpc.ClientConn) DexClient {
	return &dexClient{cc}
}

func (c *dexClient) CreateClient(ctx context.Context, in *CreateClientReq, opts ...grpc.CallOption) (*CreateClientResp, error) {
	out := new(CreateClientResp)
	err := grpc.Invoke(ctx, "/api.Dex/CreateClient", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dexClient) DeleteClient(ctx context.Context, in *DeleteClientReq, opts ...grpc.CallOption) (*DeleteClientResp, error) {
	out := new(DeleteClientResp)
	err := grpc.Invoke(ctx, "/api.Dex/DeleteClient", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dexClient) CreatePassword(ctx context.Context, in *CreatePasswordReq, opts ...grpc.CallOption) (*CreatePasswordResp, error) {
	out := new(CreatePasswordResp)
	err := grpc.Invoke(ctx, "/api.Dex/CreatePassword", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dexClient) UpdatePassword(ctx context.Context, in *UpdatePasswordReq, opts ...grpc.CallOption) (*UpdatePasswordResp, error) {
	out := new(UpdatePasswordResp)
	err := grpc.Invoke(ctx, "/api.Dex/UpdatePassword", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dexClient) DeletePassword(ctx context.Context, in *DeletePasswordReq, opts ...grpc.CallOption) (*DeletePasswordResp, error) {
	out := new(DeletePasswordResp)
	err := grpc.Invoke(ctx, "/api.Dex/DeletePassword", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dexClient) GetVersion(ctx context.Context, in *VersionReq, opts ...grpc.CallOption) (*Version, error) {
	out := new(Version)
	err := grpc.Invoke(ctx, "/api.Dex/GetVersion", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Dex service

type DexServer interface {
	// CreateClient attempts to create the client.
	CreateClient(context.Context, *CreateClientReq) (*CreateClientResp, error)
	// DeleteClient attempts to delete the provided client.
	DeleteClient(context.Context, *DeleteClientReq) (*DeleteClientResp, error)
	// CreatePassword attempts to create the password.
	CreatePassword(context.Context, *CreatePasswordReq) (*CreatePasswordResp, error)
	// UpdatePassword attempts to modify existing password.
	UpdatePassword(context.Context, *UpdatePasswordReq) (*UpdatePasswordResp, error)
	// DeletePassword attempts to delete the password.
	DeletePassword(context.Context, *DeletePasswordReq) (*DeletePasswordResp, error)
	// Version attempts to get the version info.
	GetVersion(context.Context, *VersionReq) (*Version, error)
}

func RegisterDexServer(s *grpc.Server, srv DexServer) {
	s.RegisterService(&_Dex_serviceDesc, srv)
}

func _Dex_CreateClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateClientReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DexServer).CreateClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Dex/CreateClient",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DexServer).CreateClient(ctx, req.(*CreateClientReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dex_DeleteClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteClientReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DexServer).DeleteClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Dex/DeleteClient",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DexServer).DeleteClient(ctx, req.(*DeleteClientReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dex_CreatePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePasswordReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DexServer).CreatePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Dex/CreatePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DexServer).CreatePassword(ctx, req.(*CreatePasswordReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dex_UpdatePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePasswordReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DexServer).UpdatePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Dex/UpdatePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DexServer).UpdatePassword(ctx, req.(*UpdatePasswordReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dex_DeletePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePasswordReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DexServer).DeletePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Dex/DeletePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DexServer).DeletePassword(ctx, req.(*DeletePasswordReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dex_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VersionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DexServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Dex/GetVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DexServer).GetVersion(ctx, req.(*VersionReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Dex_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Dex",
	HandlerType: (*DexServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateClient",
			Handler:    _Dex_CreateClient_Handler,
		},
		{
			MethodName: "DeleteClient",
			Handler:    _Dex_DeleteClient_Handler,
		},
		{
			MethodName: "CreatePassword",
			Handler:    _Dex_CreatePassword_Handler,
		},
		{
			MethodName: "UpdatePassword",
			Handler:    _Dex_UpdatePassword_Handler,
		},
		{
			MethodName: "DeletePassword",
			Handler:    _Dex_DeletePassword_Handler,
		},
		{
			MethodName: "GetVersion",
			Handler:    _Dex_GetVersion_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("api/api.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 577 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x54, 0xcb, 0x6e, 0xdb, 0x30,
	0x10, 0xb4, 0x2d, 0x3f, 0xe4, 0xb5, 0xfc, 0x22, 0xd2, 0x44, 0x71, 0x2f, 0x0e, 0x83, 0x02, 0xce,
	0xa1, 0x09, 0x9a, 0x00, 0xbd, 0x14, 0xed, 0xc5, 0xe9, 0xeb, 0x16, 0x08, 0x70, 0x8f, 0x15, 0x14,
	0x6b, 0x9b, 0x10, 0x50, 0x24, 0x85, 0xa4, 0xea, 0xf4, 0xd8, 0x4f, 0xeb, 0x9f, 0x15, 0xa4, 0x68,
	0x57, 0x92, 0x5d, 0x38, 0x37, 0xed, 0x70, 0x77, 0x96, 0xb3, 0xb3, 0x14, 0xf4, 0x83, 0x94, 0x5d,
	0x04, 0x29, 0x3b, 0x4f, 0x79, 0x22, 0x13, 0x62, 0x05, 0x29, 0xa3, 0x7f, 0xea, 0xd0, 0x9e, 0x47,
	0x0c, 0x63, 0x49, 0x06, 0xd0, 0x60, 0xa1, 0x5b, 0x9f, 0xd6, 0x67, 0x5d, 0xaf, 0xc1, 0x42, 0x72,
	0x08, 0x6d, 0x81, 0x4b, 0x8e, 0xd2, 0x6d, 0x68, 0xcc, 0x44, 0xe4, 0x14, 0xfa, 0x1c, 0x43, 0xc6,
	0x71, 0x29, 0xfd, 0x8c, 0x33, 0xe1, 0x5a, 0x53, 0x6b, 0xd6, 0xf5, 0x9c, 0x35, 0xb8, 0xe0, 0x4c,
	0xa8, 0x24, 0xc9, 0x33, 0x21, 0x31, 0xf4, 0x53, 0x44, 0x2e, 0xdc, 0x66, 0x9e, 0x64, 0xc0, 0x1b,
	0x85, 0xa9, 0x0e, 0x69, 0x76, 0x1b, 0xb1, 0xa5, 0xdb, 0x9a, 0xd6, 0x67, 0xb6, 0x67, 0x22, 0x42,
	0xa0, 0x19, 0x07, 0x0f, 0xe8, 0xb6, 0x75, 0x5f, 0xfd, 0x4d, 0x8e, 0xc1, 0x8e, 0x92, 0xbb, 0xc4,
	0xcf, 0x78, 0xe4, 0x76, 0x34, 0xde, 0x51, 0xf1, 0x82, 0x47, 0xf4, 0x2d, 0x0c, 0xe7, 0x1c, 0x03,
	0x89, 0xb9, 0x10, 0x0f, 0x1f, 0xc9, 0x29, 0xb4, 0x97, 0x3a, 0xd0, 0x7a, 0x7a, 0x97, 0xbd, 0x73,
	0xa5, 0xdb, 0x9c, 0x9b, 0x23, 0xfa, 0x1d, 0x46, 0xe5, 0x3a, 0x91, 0x92, 0x57, 0x30, 0x08, 0x22,
	0x8e, 0x41, 0xf8, 0xcb, 0xc7, 0x27, 0x26, 0xa4, 0xd0, 0x04, 0xb6, 0xd7, 0x37, 0xe8, 0x47, 0x0d,
	0x16, 0xf8, 0x1b, 0xff, 0xe7, 0x3f, 0x81, 0xe1, 0x35, 0x46, 0x58, 0xbc, 0x57, 0x65, 0xc6, 0xf4,
	0x02, 0x46, 0xe5, 0x14, 0x91, 0x92, 0x97, 0xd0, 0x8d, 0x13, 0xe9, 0xff, 0x48, 0xb2, 0x38, 0x34,
	0xdd, 0xed, 0x38, 0x91, 0x9f, 0x54, 0x4c, 0x19, 0xd8, 0x37, 0x81, 0x10, 0xab, 0x84, 0x87, 0xe4,
	0x00, 0x5a, 0xf8, 0x10, 0xb0, 0xc8, 0xf0, 0xe5, 0x81, 0x1a, 0xde, 0x7d, 0x20, 0xee, 0xf5, 0xc5,
	0x1c, 0x4f, 0x7f, 0x93, 0x09, 0xd8, 0x99, 0x40, 0xae, 0x87, 0x6a, 0xe9, 0xe4, 0x4d, 0x4c, 0x8e,
	0xa0, 0xa3, 0xbe, 0x7d, 0x16, 0xba, 0xcd, 0xdc, 0x67, 0x15, 0x7e, 0x0d, 0xe9, 0x07, 0x18, 0xe7,
	0xe3, 0x59, 0x37, 0x54, 0x02, 0xce, 0xc0, 0x4e, 0x4d, 0x68, 0x46, 0xdb, 0xd7, 0xd2, 0x37, 0x39,
	0x9b, 0x63, 0xfa, 0x0e, 0x48, 0xb5, 0xfe, 0xd9, 0x03, 0xa6, 0x77, 0x30, 0x5e, 0xa4, 0x61, 0xa5,
	0xf9, 0x6e, 0xc1, 0xc7, 0x60, 0xc7, 0xb8, 0xf2, 0x0b, 0xa2, 0x3b, 0x31, 0xae, 0xbe, 0x28, 0xdd,
	0x27, 0xe0, 0xa8, 0xa3, 0x8a, 0xf6, 0x5e, 0x8c, 0xab, 0x85, 0x81, 0xe8, 0x1b, 0x20, 0xd5, 0x46,
	0xfb, 0x3c, 0x38, 0x83, 0x71, 0x6e, 0xda, 0xde, 0xbb, 0x29, 0xf6, 0x6a, 0xea, 0x3e, 0xf6, 0x2b,
	0xe8, 0x7c, 0x43, 0x2e, 0x58, 0x12, 0xe7, 0x2f, 0x90, 0xff, 0x44, 0x6e, 0x48, 0x4d, 0x44, 0x46,
	0xa0, 0xde, 0xae, 0x16, 0xdb, 0xf2, 0xf4, 0x33, 0x76, 0x00, 0x4c, 0x91, 0x87, 0x8f, 0x97, 0xbf,
	0x2d, 0xb0, 0xae, 0xf1, 0x89, 0xbc, 0x07, 0xa7, 0xb8, 0xe0, 0xe4, 0x20, 0xdf, 0xd2, 0xf2, 0x5b,
	0x99, 0xbc, 0xd8, 0x81, 0x8a, 0x94, 0xd6, 0x54, 0x79, 0x71, 0x39, 0x4d, 0x79, 0x65, 0xa5, 0x4d,
	0x79, 0x75, 0x8b, 0x69, 0x8d, 0xcc, 0x61, 0x50, 0xf6, 0x9f, 0x1c, 0x16, 0x3a, 0x15, 0x66, 0x37,
	0x39, 0xda, 0x89, 0xaf, 0x49, 0xca, 0xf6, 0x18, 0x92, 0xad, 0xe5, 0x30, 0x24, 0xdb, 0x5e, 0xe6,
	0x24, 0x65, 0x17, 0x0c, 0xc9, 0x96, 0x8b, 0x86, 0x64, 0xdb, 0x32, 0x5a, 0x23, 0xaf, 0x01, 0x3e,
	0xa3, 0x5c, 0x5b, 0x33, 0xd4, 0x89, 0xff, 0x66, 0x3e, 0x71, 0x8a, 0x00, 0xad, 0xdd, 0xb6, 0xf5,
	0x4f, 0xf6, 0xea, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x95, 0x88, 0xa3, 0x2b, 0x75, 0x05, 0x00,
	0x00,
}
