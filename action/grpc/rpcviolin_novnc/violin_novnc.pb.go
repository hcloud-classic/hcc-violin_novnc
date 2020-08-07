// Copyright 2015 gRPC authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.1
// source: violin_novnc.proto

package rpcviolin_novnc

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	rpcmsgType "hcc/violin-novnc/action/grpc/rpcmsgType"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Symbols defined in public import of msgType.proto.

type Empty = rpcmsgType.Empty
type Node = rpcmsgType.Node
type NodeDetail = rpcmsgType.NodeDetail
type Server = rpcmsgType.Server
type ServerNode = rpcmsgType.ServerNode
type VNC = rpcmsgType.VNC
type Volume = rpcmsgType.Volume
type VolumeAttachment = rpcmsgType.VolumeAttachment
type AdaptiveIPSetting = rpcmsgType.AdaptiveIPSetting
type AdaptiveIPAvailableIPList = rpcmsgType.AdaptiveIPAvailableIPList
type AdaptiveIPServer = rpcmsgType.AdaptiveIPServer
type Subnet = rpcmsgType.Subnet
type Series = rpcmsgType.Series
type Telegraf = rpcmsgType.Telegraf
type NormalAction = rpcmsgType.NormalAction
type HccAction = rpcmsgType.HccAction
type Action = rpcmsgType.Action
type Control = rpcmsgType.Control
type Controls = rpcmsgType.Controls

type ReqNoVNC struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vncs []*rpcmsgType.VNC `protobuf:"bytes,1,rep,name=vncs,proto3" json:"vncs,omitempty"`
}

func (x *ReqNoVNC) Reset() {
	*x = ReqNoVNC{}
	if protoimpl.UnsafeEnabled {
		mi := &file_violin_novnc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqNoVNC) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqNoVNC) ProtoMessage() {}

func (x *ReqNoVNC) ProtoReflect() protoreflect.Message {
	mi := &file_violin_novnc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqNoVNC.ProtoReflect.Descriptor instead.
func (*ReqNoVNC) Descriptor() ([]byte, []int) {
	return file_violin_novnc_proto_rawDescGZIP(), []int{0}
}

func (x *ReqNoVNC) GetVncs() []*rpcmsgType.VNC {
	if x != nil {
		return x.Vncs
	}
	return nil
}

type ResNoVNC struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vncs []*rpcmsgType.VNC `protobuf:"bytes,1,rep,name=vncs,proto3" json:"vncs,omitempty"`
}

func (x *ResNoVNC) Reset() {
	*x = ResNoVNC{}
	if protoimpl.UnsafeEnabled {
		mi := &file_violin_novnc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResNoVNC) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResNoVNC) ProtoMessage() {}

func (x *ResNoVNC) ProtoReflect() protoreflect.Message {
	mi := &file_violin_novnc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResNoVNC.ProtoReflect.Descriptor instead.
func (*ResNoVNC) Descriptor() ([]byte, []int) {
	return file_violin_novnc_proto_rawDescGZIP(), []int{1}
}

func (x *ResNoVNC) GetVncs() []*rpcmsgType.VNC {
	if x != nil {
		return x.Vncs
	}
	return nil
}

var File_violin_novnc_proto protoreflect.FileDescriptor

var file_violin_novnc_proto_rawDesc = []byte{
	0x0a, 0x12, 0x76, 0x69, 0x6f, 0x6c, 0x69, 0x6e, 0x5f, 0x6e, 0x6f, 0x76, 0x6e, 0x63, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x52, 0x70, 0x63, 0x4e, 0x6f, 0x56, 0x4e, 0x43, 0x1a, 0x0d,
	0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2c, 0x0a,
	0x08, 0x52, 0x65, 0x71, 0x4e, 0x6f, 0x56, 0x4e, 0x43, 0x12, 0x20, 0x0a, 0x04, 0x76, 0x6e, 0x63,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70,
	0x65, 0x2e, 0x56, 0x4e, 0x43, 0x52, 0x04, 0x76, 0x6e, 0x63, 0x73, 0x22, 0x2c, 0x0a, 0x08, 0x52,
	0x65, 0x73, 0x4e, 0x6f, 0x56, 0x4e, 0x43, 0x12, 0x20, 0x0a, 0x04, 0x76, 0x6e, 0x63, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x2e,
	0x56, 0x4e, 0x43, 0x52, 0x04, 0x76, 0x6e, 0x63, 0x73, 0x32, 0x76, 0x0a, 0x05, 0x6e, 0x6f, 0x76,
	0x6e, 0x63, 0x12, 0x35, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x56, 0x4e, 0x43, 0x12,
	0x12, 0x2e, 0x52, 0x70, 0x63, 0x4e, 0x6f, 0x56, 0x4e, 0x43, 0x2e, 0x52, 0x65, 0x71, 0x4e, 0x6f,
	0x56, 0x4e, 0x43, 0x1a, 0x12, 0x2e, 0x52, 0x70, 0x63, 0x4e, 0x6f, 0x56, 0x4e, 0x43, 0x2e, 0x52,
	0x65, 0x73, 0x4e, 0x6f, 0x56, 0x4e, 0x43, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0a, 0x43, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x56, 0x4e, 0x43, 0x12, 0x12, 0x2e, 0x52, 0x70, 0x63, 0x4e, 0x6f, 0x56,
	0x4e, 0x43, 0x2e, 0x52, 0x65, 0x71, 0x4e, 0x6f, 0x56, 0x4e, 0x43, 0x1a, 0x12, 0x2e, 0x52, 0x70,
	0x63, 0x4e, 0x6f, 0x56, 0x4e, 0x43, 0x2e, 0x52, 0x65, 0x73, 0x4e, 0x6f, 0x56, 0x4e, 0x43, 0x22,
	0x00, 0x42, 0x2e, 0x5a, 0x2c, 0x68, 0x63, 0x63, 0x2f, 0x76, 0x69, 0x6f, 0x6c, 0x69, 0x6e, 0x2d,
	0x6e, 0x6f, 0x76, 0x6e, 0x63, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x67, 0x72, 0x70,
	0x63, 0x2f, 0x72, 0x70, 0x63, 0x76, 0x69, 0x6f, 0x6c, 0x69, 0x6e, 0x5f, 0x6e, 0x6f, 0x76, 0x6e,
	0x63, 0x50, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_violin_novnc_proto_rawDescOnce sync.Once
	file_violin_novnc_proto_rawDescData = file_violin_novnc_proto_rawDesc
)

func file_violin_novnc_proto_rawDescGZIP() []byte {
	file_violin_novnc_proto_rawDescOnce.Do(func() {
		file_violin_novnc_proto_rawDescData = protoimpl.X.CompressGZIP(file_violin_novnc_proto_rawDescData)
	})
	return file_violin_novnc_proto_rawDescData
}

var file_violin_novnc_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_violin_novnc_proto_goTypes = []interface{}{
	(*ReqNoVNC)(nil),       // 0: RpcNoVNC.ReqNoVNC
	(*ResNoVNC)(nil),       // 1: RpcNoVNC.ResNoVNC
	(*rpcmsgType.VNC)(nil), // 2: MsgType.VNC
}
var file_violin_novnc_proto_depIdxs = []int32{
	2, // 0: RpcNoVNC.ReqNoVNC.vncs:type_name -> MsgType.VNC
	2, // 1: RpcNoVNC.ResNoVNC.vncs:type_name -> MsgType.VNC
	0, // 2: RpcNoVNC.novnc.CreateVNC:input_type -> RpcNoVNC.ReqNoVNC
	0, // 3: RpcNoVNC.novnc.ControlVNC:input_type -> RpcNoVNC.ReqNoVNC
	1, // 4: RpcNoVNC.novnc.CreateVNC:output_type -> RpcNoVNC.ResNoVNC
	1, // 5: RpcNoVNC.novnc.ControlVNC:output_type -> RpcNoVNC.ResNoVNC
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_violin_novnc_proto_init() }
func file_violin_novnc_proto_init() {
	if File_violin_novnc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_violin_novnc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqNoVNC); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_violin_novnc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResNoVNC); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_violin_novnc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_violin_novnc_proto_goTypes,
		DependencyIndexes: file_violin_novnc_proto_depIdxs,
		MessageInfos:      file_violin_novnc_proto_msgTypes,
	}.Build()
	File_violin_novnc_proto = out.File
	file_violin_novnc_proto_rawDesc = nil
	file_violin_novnc_proto_goTypes = nil
	file_violin_novnc_proto_depIdxs = nil
}
