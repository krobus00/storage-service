// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.22.2
// source: pb/storage/storage.proto

package storage

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

type Object struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	FileName   string `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name"`
	Type       string `protobuf:"bytes,3,opt,name=type,proto3" json:"type"`
	SignedUrl  string `protobuf:"bytes,4,opt,name=signed_url,json=signedUrl,proto3" json:"signed_url"`
	ExpiredAt  string `protobuf:"bytes,5,opt,name=expired_at,json=expiredAt,proto3" json:"expired_at"`
	IsPublic   bool   `protobuf:"varint,6,opt,name=is_public,json=isPublic,proto3" json:"is_public"`
	UploadedBy string `protobuf:"bytes,7,opt,name=uploaded_by,json=uploadedBy,proto3" json:"uploaded_by"`
	CreatedAt  string `protobuf:"bytes,8,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
}

func (x *Object) Reset() {
	*x = Object{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_storage_storage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Object) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Object) ProtoMessage() {}

func (x *Object) ProtoReflect() protoreflect.Message {
	mi := &file_pb_storage_storage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Object.ProtoReflect.Descriptor instead.
func (*Object) Descriptor() ([]byte, []int) {
	return file_pb_storage_storage_proto_rawDescGZIP(), []int{0}
}

func (x *Object) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Object) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *Object) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Object) GetSignedUrl() string {
	if x != nil {
		return x.SignedUrl
	}
	return ""
}

func (x *Object) GetExpiredAt() string {
	if x != nil {
		return x.ExpiredAt
	}
	return ""
}

func (x *Object) GetIsPublic() bool {
	if x != nil {
		return x.IsPublic
	}
	return false
}

func (x *Object) GetUploadedBy() string {
	if x != nil {
		return x.UploadedBy
	}
	return ""
}

func (x *Object) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

type GetObjectByIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id"`
	ObjectId string `protobuf:"bytes,2,opt,name=object_id,json=objectId,proto3" json:"object_id"`
}

func (x *GetObjectByIDRequest) Reset() {
	*x = GetObjectByIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_storage_storage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetObjectByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetObjectByIDRequest) ProtoMessage() {}

func (x *GetObjectByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_storage_storage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetObjectByIDRequest.ProtoReflect.Descriptor instead.
func (*GetObjectByIDRequest) Descriptor() ([]byte, []int) {
	return file_pb_storage_storage_proto_rawDescGZIP(), []int{1}
}

func (x *GetObjectByIDRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetObjectByIDRequest) GetObjectId() string {
	if x != nil {
		return x.ObjectId
	}
	return ""
}

type DeleteObjectByIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id"`
	ObjectId string `protobuf:"bytes,2,opt,name=object_id,json=objectId,proto3" json:"object_id"`
}

func (x *DeleteObjectByIDRequest) Reset() {
	*x = DeleteObjectByIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_storage_storage_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteObjectByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteObjectByIDRequest) ProtoMessage() {}

func (x *DeleteObjectByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_storage_storage_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteObjectByIDRequest.ProtoReflect.Descriptor instead.
func (*DeleteObjectByIDRequest) Descriptor() ([]byte, []int) {
	return file_pb_storage_storage_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteObjectByIDRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *DeleteObjectByIDRequest) GetObjectId() string {
	if x != nil {
		return x.ObjectId
	}
	return ""
}

var File_pb_storage_storage_proto protoreflect.FileDescriptor

var file_pb_storage_storage_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x62, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x70, 0x62, 0x2e, 0x73,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x22, 0xe4, 0x01, 0x0a, 0x06, 0x4f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x55, 0x72,
	0x6c, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x12, 0x1f, 0x0a,
	0x0b, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x42, 0x79, 0x12, 0x1d,
	0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x4c, 0x0a,
	0x14, 0x47, 0x65, 0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1b,
	0x0a, 0x09, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x22, 0x4f, 0x0a, 0x17, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x42, 0x79, 0x49, 0x44, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x1b, 0x0a, 0x09, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x42, 0x0c, 0x5a, 0x0a,
	0x70, 0x62, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_pb_storage_storage_proto_rawDescOnce sync.Once
	file_pb_storage_storage_proto_rawDescData = file_pb_storage_storage_proto_rawDesc
)

func file_pb_storage_storage_proto_rawDescGZIP() []byte {
	file_pb_storage_storage_proto_rawDescOnce.Do(func() {
		file_pb_storage_storage_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_storage_storage_proto_rawDescData)
	})
	return file_pb_storage_storage_proto_rawDescData
}

var file_pb_storage_storage_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pb_storage_storage_proto_goTypes = []interface{}{
	(*Object)(nil),                  // 0: pb.storage.Object
	(*GetObjectByIDRequest)(nil),    // 1: pb.storage.GetObjectByIDRequest
	(*DeleteObjectByIDRequest)(nil), // 2: pb.storage.DeleteObjectByIDRequest
}
var file_pb_storage_storage_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pb_storage_storage_proto_init() }
func file_pb_storage_storage_proto_init() {
	if File_pb_storage_storage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_storage_storage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Object); i {
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
		file_pb_storage_storage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetObjectByIDRequest); i {
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
		file_pb_storage_storage_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteObjectByIDRequest); i {
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
			RawDescriptor: file_pb_storage_storage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pb_storage_storage_proto_goTypes,
		DependencyIndexes: file_pb_storage_storage_proto_depIdxs,
		MessageInfos:      file_pb_storage_storage_proto_msgTypes,
	}.Build()
	File_pb_storage_storage_proto = out.File
	file_pb_storage_storage_proto_rawDesc = nil
	file_pb_storage_storage_proto_goTypes = nil
	file_pb_storage_storage_proto_depIdxs = nil
}
