// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.20.3
// source: protos/common/time.proto

package gopb

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

// Type that we'll use to encode nano-second epochs because the Google variant for Timestamp doesn't
// support JSON deserialization from a numeric value
type UnixTimestamp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Represents seconds of UTC time since Unix epoch 1970-01-01T00:00:00Z. Must be from 0001-01-01T00:00:00Z
	// to 9999-12-31T23:59:59Z inclusive.
	Seconds int64 `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
	// Non-negative fractions of a second at nanosecond resolution. Negative second values with fractions
	// must still have non-negative nanos values that count forward in time. Must be from 0 to 999,999,999
	// inclusive.
	Nanoseconds int32 `protobuf:"varint,2,opt,name=nanoseconds,proto3" json:"nanoseconds,omitempty"`
}

func (x *UnixTimestamp) Reset() {
	*x = UnixTimestamp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_common_time_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnixTimestamp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnixTimestamp) ProtoMessage() {}

func (x *UnixTimestamp) ProtoReflect() protoreflect.Message {
	mi := &file_protos_common_time_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnixTimestamp.ProtoReflect.Descriptor instead.
func (*UnixTimestamp) Descriptor() ([]byte, []int) {
	return file_protos_common_time_proto_rawDescGZIP(), []int{0}
}

func (x *UnixTimestamp) GetSeconds() int64 {
	if x != nil {
		return x.Seconds
	}
	return 0
}

func (x *UnixTimestamp) GetNanoseconds() int32 {
	if x != nil {
		return x.Nanoseconds
	}
	return 0
}

// Duration represents a signed, fixed-length span of time represented as a count of seconds and fractions
// of seconds at nanosecond resolution. It is independent of any calendar and concepts like "day" or
// "month". It is related to Timestamp in that the difference between two Timestamp values is a Duration
// and it can be added or subtracted from a Timestamp. Range is approximately +-10,000 years.
type UnixDuration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Signed seconds of the span of time. Must be from -315,576,000,000 to +315,576,000,000 inclusive.
	// Note: these bounds are computed from: 60 sec/min * 60 min/hr * 24 hr/day * 365.25 days/year * 10000 years
	Seconds int64 `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
	// Signed fractions of a second at nanosecond resolution of the span of time. Durations less than
	// one second are represented with a 0 `seconds` field and a positive or negative `nanos` field.
	// For durations of one second or more, a non-zero value for the `nanos` field must be of the same
	// sign as the `seconds` field. Must be from -999,999,999 to +999,999,999 inclusive.
	Nanoseconds int32 `protobuf:"varint,2,opt,name=nanoseconds,proto3" json:"nanoseconds,omitempty"`
}

func (x *UnixDuration) Reset() {
	*x = UnixDuration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_common_time_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnixDuration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnixDuration) ProtoMessage() {}

func (x *UnixDuration) ProtoReflect() protoreflect.Message {
	mi := &file_protos_common_time_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnixDuration.ProtoReflect.Descriptor instead.
func (*UnixDuration) Descriptor() ([]byte, []int) {
	return file_protos_common_time_proto_rawDescGZIP(), []int{1}
}

func (x *UnixDuration) GetSeconds() int64 {
	if x != nil {
		return x.Seconds
	}
	return 0
}

func (x *UnixDuration) GetNanoseconds() int32 {
	if x != nil {
		return x.Nanoseconds
	}
	return 0
}

var File_protos_common_time_proto protoreflect.FileDescriptor

var file_protos_common_time_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x22, 0x4b, 0x0a, 0x0d, 0x55, 0x6e, 0x69,
	0x78, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65,
	0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x73, 0x65, 0x63,
	0x6f, 0x6e, 0x64, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x61, 0x6e, 0x6f, 0x73, 0x65, 0x63, 0x6f,
	0x6e, 0x64, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x6e, 0x61, 0x6e, 0x6f, 0x73,
	0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x22, 0x4a, 0x0a, 0x0c, 0x55, 0x6e, 0x69, 0x78, 0x44, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73,
	0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x61, 0x6e, 0x6f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x6e, 0x61, 0x6e, 0x6f, 0x73, 0x65, 0x63, 0x6f, 0x6e,
	0x64, 0x73, 0x42, 0x28, 0x5a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x78, 0x65, 0x66, 0x69, 0x6e, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x67, 0x6f, 0x2f, 0x67, 0x6f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_common_time_proto_rawDescOnce sync.Once
	file_protos_common_time_proto_rawDescData = file_protos_common_time_proto_rawDesc
)

func file_protos_common_time_proto_rawDescGZIP() []byte {
	file_protos_common_time_proto_rawDescOnce.Do(func() {
		file_protos_common_time_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_common_time_proto_rawDescData)
	})
	return file_protos_common_time_proto_rawDescData
}

var file_protos_common_time_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_protos_common_time_proto_goTypes = []interface{}{
	(*UnixTimestamp)(nil), // 0: protos.common.UnixTimestamp
	(*UnixDuration)(nil),  // 1: protos.common.UnixDuration
}
var file_protos_common_time_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protos_common_time_proto_init() }
func file_protos_common_time_proto_init() {
	if File_protos_common_time_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_common_time_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnixTimestamp); i {
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
		file_protos_common_time_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnixDuration); i {
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
			RawDescriptor: file_protos_common_time_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protos_common_time_proto_goTypes,
		DependencyIndexes: file_protos_common_time_proto_depIdxs,
		MessageInfos:      file_protos_common_time_proto_msgTypes,
	}.Build()
	File_protos_common_time_proto = out.File
	file_protos_common_time_proto_rawDesc = nil
	file_protos_common_time_proto_goTypes = nil
	file_protos_common_time_proto_depIdxs = nil
}
