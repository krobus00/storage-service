// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: pb/storage/storage_service.proto

package storage

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_pb_storage_storage_service_proto protoreflect.FileDescriptor

var file_pb_storage_storage_service_proto_rawDesc = []byte{
	0x0a, 0x20, 0x70, 0x62, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x70, 0x62, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x1a, 0x18,
	0x70, 0x62, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x5a, 0x0a, 0x0e, 0x53, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x48, 0x0a, 0x0d, 0x47, 0x65,
	0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x42, 0x79, 0x49, 0x44, 0x12, 0x20, 0x2e, 0x70, 0x62,
	0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e,
	0x70, 0x62, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x22, 0x00, 0x42, 0x0c, 0x5a, 0x0a, 0x70, 0x62, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_pb_storage_storage_service_proto_goTypes = []interface{}{
	(*GetObjectByIDRequest)(nil), // 0: pb.storage.GetObjectByIDRequest
	(*Storage)(nil),              // 1: pb.storage.Storage
}
var file_pb_storage_storage_service_proto_depIdxs = []int32{
	0, // 0: pb.storage.StorageService.GetObjectByID:input_type -> pb.storage.GetObjectByIDRequest
	1, // 1: pb.storage.StorageService.GetObjectByID:output_type -> pb.storage.Storage
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pb_storage_storage_service_proto_init() }
func file_pb_storage_storage_service_proto_init() {
	if File_pb_storage_storage_service_proto != nil {
		return
	}
	file_pb_storage_storage_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pb_storage_storage_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_storage_storage_service_proto_goTypes,
		DependencyIndexes: file_pb_storage_storage_service_proto_depIdxs,
	}.Build()
	File_pb_storage_storage_service_proto = out.File
	file_pb_storage_storage_service_proto_rawDesc = nil
	file_pb_storage_storage_service_proto_goTypes = nil
	file_pb_storage_storage_service_proto_depIdxs = nil
}
