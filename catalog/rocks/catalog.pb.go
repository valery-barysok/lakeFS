// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: catalog.proto

package rocks

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

type Entry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address      string                 `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	LastModified *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=last_modified,json=lastModified,proto3" json:"last_modified,omitempty"`
	Size         int64                  `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	ETag         string                 `protobuf:"bytes,4,opt,name=e_tag,json=eTag,proto3" json:"e_tag,omitempty"`
	Metadata     map[string]string      `protobuf:"bytes,5,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Entry) Reset() {
	*x = Entry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_catalog_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entry) ProtoMessage() {}

func (x *Entry) ProtoReflect() protoreflect.Message {
	mi := &file_catalog_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entry.ProtoReflect.Descriptor instead.
func (*Entry) Descriptor() ([]byte, []int) {
	return file_catalog_proto_rawDescGZIP(), []int{0}
}

func (x *Entry) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Entry) GetLastModified() *timestamppb.Timestamp {
	if x != nil {
		return x.LastModified
	}
	return nil
}

func (x *Entry) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *Entry) GetETag() string {
	if x != nil {
		return x.ETag
	}
	return ""
}

func (x *Entry) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

var File_catalog_proto protoreflect.FileDescriptor

var file_catalog_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x82, 0x02, 0x0a, 0x05, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x3f, 0x0a,
	0x0d, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x4d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x69,
	0x7a, 0x65, 0x12, 0x13, 0x0a, 0x05, 0x65, 0x5f, 0x74, 0x61, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x65, 0x54, 0x61, 0x67, 0x12, 0x38, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x61, 0x74, 0x61,
	0x6c, 0x6f, 0x67, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x2a,
	0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x72, 0x65,
	0x65, 0x76, 0x65, 0x73, 0x65, 0x2f, 0x6c, 0x61, 0x6b, 0x65, 0x66, 0x73, 0x2f, 0x63, 0x61, 0x74,
	0x61, 0x6c, 0x6f, 0x67, 0x2f, 0x72, 0x6f, 0x63, 0x6b, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_catalog_proto_rawDescOnce sync.Once
	file_catalog_proto_rawDescData = file_catalog_proto_rawDesc
)

func file_catalog_proto_rawDescGZIP() []byte {
	file_catalog_proto_rawDescOnce.Do(func() {
		file_catalog_proto_rawDescData = protoimpl.X.CompressGZIP(file_catalog_proto_rawDescData)
	})
	return file_catalog_proto_rawDescData
}

var file_catalog_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_catalog_proto_goTypes = []interface{}{
	(*Entry)(nil),                 // 0: catalog.Entry
	nil,                           // 1: catalog.Entry.MetadataEntry
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_catalog_proto_depIdxs = []int32{
	2, // 0: catalog.Entry.last_modified:type_name -> google.protobuf.Timestamp
	1, // 1: catalog.Entry.metadata:type_name -> catalog.Entry.MetadataEntry
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_catalog_proto_init() }
func file_catalog_proto_init() {
	if File_catalog_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_catalog_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entry); i {
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
			RawDescriptor: file_catalog_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_catalog_proto_goTypes,
		DependencyIndexes: file_catalog_proto_depIdxs,
		MessageInfos:      file_catalog_proto_msgTypes,
	}.Build()
	File_catalog_proto = out.File
	file_catalog_proto_rawDesc = nil
	file_catalog_proto_goTypes = nil
	file_catalog_proto_depIdxs = nil
}
