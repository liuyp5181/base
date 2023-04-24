// protoc --go_out=.  --go-grpc_out=. *.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: monitor.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type WatchType int32

const (
	WatchType_PUT    WatchType = 0
	WatchType_DELETE WatchType = 1
)

// Enum value maps for WatchType.
var (
	WatchType_name = map[int32]string{
		0: "PUT",
		1: "DELETE",
	}
	WatchType_value = map[string]int32{
		"PUT":    0,
		"DELETE": 1,
	}
)

func (x WatchType) Enum() *WatchType {
	p := new(WatchType)
	*p = x
	return p
}

func (x WatchType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WatchType) Descriptor() protoreflect.EnumDescriptor {
	return file_config_proto_enumTypes[0].Descriptor()
}

func (WatchType) Type() protoreflect.EnumType {
	return &file_config_proto_enumTypes[0]
}

func (x WatchType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WatchType.Descriptor instead.
func (WatchType) EnumDescriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{0}
}

type GetConfigReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *GetConfigReq) Reset() {
	*x = GetConfigReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetConfigReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetConfigReq) ProtoMessage() {}

func (x *GetConfigReq) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetConfigReq.ProtoReflect.Descriptor instead.
func (*GetConfigReq) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{0}
}

func (x *GetConfigReq) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type GetConfigRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Val string `protobuf:"bytes,1,opt,name=val,proto3" json:"val,omitempty"` // json
}

func (x *GetConfigRes) Reset() {
	*x = GetConfigRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetConfigRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetConfigRes) ProtoMessage() {}

func (x *GetConfigRes) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetConfigRes.ProtoReflect.Descriptor instead.
func (*GetConfigRes) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{1}
}

func (x *GetConfigRes) GetVal() string {
	if x != nil {
		return x.Val
	}
	return ""
}

type WatchReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *WatchReq) Reset() {
	*x = WatchReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WatchReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WatchReq) ProtoMessage() {}

func (x *WatchReq) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WatchReq.ProtoReflect.Descriptor instead.
func (*WatchReq) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{2}
}

func (x *WatchReq) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type WatchRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type WatchType `protobuf:"varint,1,opt,name=type,proto3,enum=configmgr.WatchType" json:"type,omitempty"`
	Key  string    `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Val  []byte    `protobuf:"bytes,3,opt,name=val,proto3" json:"val,omitempty"`
}

func (x *WatchRes) Reset() {
	*x = WatchRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WatchRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WatchRes) ProtoMessage() {}

func (x *WatchRes) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WatchRes.ProtoReflect.Descriptor instead.
func (*WatchRes) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{3}
}

func (x *WatchRes) GetType() WatchType {
	if x != nil {
		return x.Type
	}
	return WatchType_PUT
}

func (x *WatchRes) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *WatchRes) GetVal() []byte {
	if x != nil {
		return x.Val
	}
	return nil
}

var File_config_proto protoreflect.FileDescriptor

var file_config_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x6d, 0x67, 0x72, 0x22, 0x20, 0x0a, 0x0c, 0x47, 0x65, 0x74,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x20, 0x0a, 0x0c, 0x47,
	0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x76,
	0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x22, 0x1c, 0x0a,
	0x08, 0x57, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x58, 0x0a, 0x08, 0x57,
	0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x12, 0x28, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x6d, 0x67,
	0x72, 0x2e, 0x57, 0x61, 0x74, 0x63, 0x68, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x03, 0x76, 0x61, 0x6c, 0x2a, 0x20, 0x0a, 0x09, 0x57, 0x61, 0x74, 0x63, 0x68, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x07, 0x0a, 0x03, 0x50, 0x55, 0x54, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x44,
	0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x01, 0x32, 0x81, 0x01, 0x0a, 0x07, 0x47, 0x72, 0x65, 0x65,
	0x74, 0x65, 0x72, 0x12, 0x3f, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x12, 0x17, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x6d, 0x67, 0x72, 0x2e, 0x47, 0x65, 0x74,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x6d, 0x67, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52,
	0x65, 0x73, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x05, 0x57, 0x61, 0x74, 0x63, 0x68, 0x12, 0x13, 0x2e,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x6d, 0x67, 0x72, 0x2e, 0x57, 0x61, 0x74, 0x63, 0x68, 0x52,
	0x65, 0x71, 0x1a, 0x13, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x6d, 0x67, 0x72, 0x2e, 0x57,
	0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x22, 0x00, 0x30, 0x01, 0x42, 0x08, 0x5a, 0x06, 0x2e,
	0x2f, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_config_proto_rawDescOnce sync.Once
	file_config_proto_rawDescData = file_config_proto_rawDesc
)

func file_config_proto_rawDescGZIP() []byte {
	file_config_proto_rawDescOnce.Do(func() {
		file_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_config_proto_rawDescData)
	})
	return file_config_proto_rawDescData
}

var file_config_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_config_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_config_proto_goTypes = []interface{}{
	(WatchType)(0),       // 0: configmgr.WatchType
	(*GetConfigReq)(nil), // 1: configmgr.GetConfigReq
	(*GetConfigRes)(nil), // 2: configmgr.GetConfigRes
	(*WatchReq)(nil),     // 3: configmgr.WatchReq
	(*WatchRes)(nil),     // 4: configmgr.WatchRes
}
var file_config_proto_depIdxs = []int32{
	0, // 0: configmgr.WatchRes.type:type_name -> configmgr.WatchType
	1, // 1: configmgr.Greeter.GetConfig:input_type -> configmgr.GetConfigReq
	3, // 2: configmgr.Greeter.Watch:input_type -> configmgr.WatchReq
	2, // 3: configmgr.Greeter.GetConfig:output_type -> configmgr.GetConfigRes
	4, // 4: configmgr.Greeter.Watch:output_type -> configmgr.WatchRes
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_config_proto_init() }
func file_config_proto_init() {
	if File_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetConfigReq); i {
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
		file_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetConfigRes); i {
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
		file_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WatchReq); i {
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
		file_config_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WatchRes); i {
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
			RawDescriptor: file_config_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_config_proto_goTypes,
		DependencyIndexes: file_config_proto_depIdxs,
		EnumInfos:         file_config_proto_enumTypes,
		MessageInfos:      file_config_proto_msgTypes,
	}.Build()
	File_config_proto = out.File
	file_config_proto_rawDesc = nil
	file_config_proto_goTypes = nil
	file_config_proto_depIdxs = nil
}
