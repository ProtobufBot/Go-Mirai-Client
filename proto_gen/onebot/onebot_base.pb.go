// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.1
// source: onebot_base.proto

package onebot

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

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type string            `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Data map[string]string `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_onebot_base_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_onebot_base_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_onebot_base_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Message) GetData() map[string]string {
	if x != nil {
		return x.Data
	}
	return nil
}

type MessageReceipt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderId int64   `protobuf:"varint,1,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	Time     int64   `protobuf:"varint,2,opt,name=time,proto3" json:"time,omitempty"`
	Seqs     []int32 `protobuf:"varint,3,rep,packed,name=seqs,proto3" json:"seqs,omitempty"`
	Rands    []int32 `protobuf:"varint,4,rep,packed,name=rands,proto3" json:"rands,omitempty"`
	GroupId  int64   `protobuf:"varint,5,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
}

func (x *MessageReceipt) Reset() {
	*x = MessageReceipt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_onebot_base_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageReceipt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageReceipt) ProtoMessage() {}

func (x *MessageReceipt) ProtoReflect() protoreflect.Message {
	mi := &file_onebot_base_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageReceipt.ProtoReflect.Descriptor instead.
func (*MessageReceipt) Descriptor() ([]byte, []int) {
	return file_onebot_base_proto_rawDescGZIP(), []int{1}
}

func (x *MessageReceipt) GetSenderId() int64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *MessageReceipt) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *MessageReceipt) GetSeqs() []int32 {
	if x != nil {
		return x.Seqs
	}
	return nil
}

func (x *MessageReceipt) GetRands() []int32 {
	if x != nil {
		return x.Rands
	}
	return nil
}

func (x *MessageReceipt) GetGroupId() int64 {
	if x != nil {
		return x.GroupId
	}
	return 0
}

var File_onebot_base_proto protoreflect.FileDescriptor

var file_onebot_base_proto_rawDesc = []byte{
	0x0a, 0x11, 0x6f, 0x6e, 0x65, 0x62, 0x6f, 0x74, 0x5f, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6f, 0x6e, 0x65, 0x62, 0x6f, 0x74, 0x22, 0x85, 0x01, 0x0a, 0x07,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x2d, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x6f, 0x6e, 0x65, 0x62,
	0x6f, 0x74, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x37, 0x0a, 0x09, 0x44, 0x61,
	0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x22, 0x86, 0x01, 0x0a, 0x0e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x65, 0x71, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x05, 0x52, 0x04, 0x73, 0x65, 0x71, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x72,
	0x61, 0x6e, 0x64, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x05, 0x52, 0x05, 0x72, 0x61, 0x6e, 0x64,
	0x73, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x07, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x42, 0x0a, 0x5a, 0x08,
	0x2e, 0x2f, 0x6f, 0x6e, 0x65, 0x62, 0x6f, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_onebot_base_proto_rawDescOnce sync.Once
	file_onebot_base_proto_rawDescData = file_onebot_base_proto_rawDesc
)

func file_onebot_base_proto_rawDescGZIP() []byte {
	file_onebot_base_proto_rawDescOnce.Do(func() {
		file_onebot_base_proto_rawDescData = protoimpl.X.CompressGZIP(file_onebot_base_proto_rawDescData)
	})
	return file_onebot_base_proto_rawDescData
}

var file_onebot_base_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_onebot_base_proto_goTypes = []interface{}{
	(*Message)(nil),        // 0: onebot.Message
	(*MessageReceipt)(nil), // 1: onebot.MessageReceipt
	nil,                    // 2: onebot.Message.DataEntry
}
var file_onebot_base_proto_depIdxs = []int32{
	2, // 0: onebot.Message.data:type_name -> onebot.Message.DataEntry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_onebot_base_proto_init() }
func file_onebot_base_proto_init() {
	if File_onebot_base_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_onebot_base_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_onebot_base_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageReceipt); i {
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
			RawDescriptor: file_onebot_base_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_onebot_base_proto_goTypes,
		DependencyIndexes: file_onebot_base_proto_depIdxs,
		MessageInfos:      file_onebot_base_proto_msgTypes,
	}.Build()
	File_onebot_base_proto = out.File
	file_onebot_base_proto_rawDesc = nil
	file_onebot_base_proto_goTypes = nil
	file_onebot_base_proto_depIdxs = nil
}
