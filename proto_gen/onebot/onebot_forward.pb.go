// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: onebot_forward.proto

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

type ForwardMessageNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message []*Message `protobuf:"bytes,1,rep,name=message,proto3" json:"message,omitempty"`
}

func (x *ForwardMessageNode) Reset() {
	*x = ForwardMessageNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_onebot_forward_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ForwardMessageNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForwardMessageNode) ProtoMessage() {}

func (x *ForwardMessageNode) ProtoReflect() protoreflect.Message {
	mi := &file_onebot_forward_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForwardMessageNode.ProtoReflect.Descriptor instead.
func (*ForwardMessageNode) Descriptor() ([]byte, []int) {
	return file_onebot_forward_proto_rawDescGZIP(), []int{0}
}

func (x *ForwardMessageNode) GetMessage() []*Message {
	if x != nil {
		return x.Message
	}
	return nil
}

type ForwardContent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderId   int64  `protobuf:"varint,1,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	Time       int32  `protobuf:"varint,2,opt,name=time,proto3" json:"time,omitempty"`
	SenderName string `protobuf:"bytes,3,opt,name=sender_name,json=senderName,proto3" json:"sender_name,omitempty"`
	// Types that are assignable to Content:
	//
	//	*ForwardContent_MessageNode
	//	*ForwardContent_ForwardNode
	Content isForwardContent_Content `protobuf_oneof:"Content"`
}

func (x *ForwardContent) Reset() {
	*x = ForwardContent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_onebot_forward_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ForwardContent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForwardContent) ProtoMessage() {}

func (x *ForwardContent) ProtoReflect() protoreflect.Message {
	mi := &file_onebot_forward_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForwardContent.ProtoReflect.Descriptor instead.
func (*ForwardContent) Descriptor() ([]byte, []int) {
	return file_onebot_forward_proto_rawDescGZIP(), []int{1}
}

func (x *ForwardContent) GetSenderId() int64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *ForwardContent) GetTime() int32 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *ForwardContent) GetSenderName() string {
	if x != nil {
		return x.SenderName
	}
	return ""
}

func (m *ForwardContent) GetContent() isForwardContent_Content {
	if m != nil {
		return m.Content
	}
	return nil
}

func (x *ForwardContent) GetMessageNode() *ForwardMessageNode {
	if x, ok := x.GetContent().(*ForwardContent_MessageNode); ok {
		return x.MessageNode
	}
	return nil
}

func (x *ForwardContent) GetForwardNode() int32 {
	if x, ok := x.GetContent().(*ForwardContent_ForwardNode); ok {
		return x.ForwardNode
	}
	return 0
}

type isForwardContent_Content interface {
	isForwardContent_Content()
}

type ForwardContent_MessageNode struct {
	MessageNode *ForwardMessageNode `protobuf:"bytes,101,opt,name=message_node,json=messageNode,proto3,oneof"`
}

type ForwardContent_ForwardNode struct {
	ForwardNode int32 `protobuf:"varint,102,opt,name=forward_node,json=forwardNode,proto3,oneof"`
}

func (*ForwardContent_MessageNode) isForwardContent_Content() {}

func (*ForwardContent_ForwardNode) isForwardContent_Content() {}

type ForwardChain struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content []*ForwardContent `protobuf:"bytes,1,rep,name=content,proto3" json:"content,omitempty"`
}

func (x *ForwardChain) Reset() {
	*x = ForwardChain{}
	if protoimpl.UnsafeEnabled {
		mi := &file_onebot_forward_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ForwardChain) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForwardChain) ProtoMessage() {}

func (x *ForwardChain) ProtoReflect() protoreflect.Message {
	mi := &file_onebot_forward_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForwardChain.ProtoReflect.Descriptor instead.
func (*ForwardChain) Descriptor() ([]byte, []int) {
	return file_onebot_forward_proto_rawDescGZIP(), []int{2}
}

func (x *ForwardChain) GetContent() []*ForwardContent {
	if x != nil {
		return x.Content
	}
	return nil
}

type ForwardMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chain *ForwardChain           `protobuf:"bytes,1,opt,name=chain,proto3" json:"chain,omitempty"`
	Data  map[int32]*ForwardChain `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ForwardMessage) Reset() {
	*x = ForwardMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_onebot_forward_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ForwardMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForwardMessage) ProtoMessage() {}

func (x *ForwardMessage) ProtoReflect() protoreflect.Message {
	mi := &file_onebot_forward_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForwardMessage.ProtoReflect.Descriptor instead.
func (*ForwardMessage) Descriptor() ([]byte, []int) {
	return file_onebot_forward_proto_rawDescGZIP(), []int{3}
}

func (x *ForwardMessage) GetChain() *ForwardChain {
	if x != nil {
		return x.Chain
	}
	return nil
}

func (x *ForwardMessage) GetData() map[int32]*ForwardChain {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_onebot_forward_proto protoreflect.FileDescriptor

var file_onebot_forward_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6f, 0x6e, 0x65, 0x62, 0x6f, 0x74, 0x5f, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6f, 0x6e, 0x65, 0x62, 0x6f, 0x74, 0x1a, 0x11,
	0x6f, 0x6e, 0x65, 0x62, 0x6f, 0x74, 0x5f, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x3f, 0x0a, 0x12, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x29, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6f, 0x6e, 0x65, 0x62, 0x6f,
	0x74, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0xd3, 0x01, 0x0a, 0x0e, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x3f, 0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x65, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x6f, 0x6e, 0x65, 0x62, 0x6f, 0x74, 0x2e, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x23, 0x0a, 0x0c, 0x66, 0x6f, 0x72, 0x77,
	0x61, 0x72, 0x64, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x66, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00,
	0x52, 0x0b, 0x66, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x42, 0x09, 0x0a,
	0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x40, 0x0a, 0x0c, 0x46, 0x6f, 0x72, 0x77,
	0x61, 0x72, 0x64, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x12, 0x30, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6f, 0x6e, 0x65, 0x62,
	0x6f, 0x74, 0x2e, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0xc1, 0x01, 0x0a, 0x0e, 0x46,
	0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2a, 0x0a,
	0x05, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6f,
	0x6e, 0x65, 0x62, 0x6f, 0x74, 0x2e, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x43, 0x68, 0x61,
	0x69, 0x6e, 0x52, 0x05, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x12, 0x34, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x6f, 0x6e, 0x65, 0x62, 0x6f, 0x74,
	0x2e, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a,
	0x4d, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2a,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x6f, 0x6e, 0x65, 0x62, 0x6f, 0x74, 0x2e, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x43, 0x68,
	0x61, 0x69, 0x6e, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x0a,
	0x5a, 0x08, 0x2e, 0x2f, 0x6f, 0x6e, 0x65, 0x62, 0x6f, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_onebot_forward_proto_rawDescOnce sync.Once
	file_onebot_forward_proto_rawDescData = file_onebot_forward_proto_rawDesc
)

func file_onebot_forward_proto_rawDescGZIP() []byte {
	file_onebot_forward_proto_rawDescOnce.Do(func() {
		file_onebot_forward_proto_rawDescData = protoimpl.X.CompressGZIP(file_onebot_forward_proto_rawDescData)
	})
	return file_onebot_forward_proto_rawDescData
}

var file_onebot_forward_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_onebot_forward_proto_goTypes = []interface{}{
	(*ForwardMessageNode)(nil), // 0: onebot.ForwardMessageNode
	(*ForwardContent)(nil),     // 1: onebot.ForwardContent
	(*ForwardChain)(nil),       // 2: onebot.ForwardChain
	(*ForwardMessage)(nil),     // 3: onebot.ForwardMessage
	nil,                        // 4: onebot.ForwardMessage.DataEntry
	(*Message)(nil),            // 5: onebot.Message
}
var file_onebot_forward_proto_depIdxs = []int32{
	5, // 0: onebot.ForwardMessageNode.message:type_name -> onebot.Message
	0, // 1: onebot.ForwardContent.message_node:type_name -> onebot.ForwardMessageNode
	1, // 2: onebot.ForwardChain.content:type_name -> onebot.ForwardContent
	2, // 3: onebot.ForwardMessage.chain:type_name -> onebot.ForwardChain
	4, // 4: onebot.ForwardMessage.data:type_name -> onebot.ForwardMessage.DataEntry
	2, // 5: onebot.ForwardMessage.DataEntry.value:type_name -> onebot.ForwardChain
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_onebot_forward_proto_init() }
func file_onebot_forward_proto_init() {
	if File_onebot_forward_proto != nil {
		return
	}
	file_onebot_base_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_onebot_forward_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ForwardMessageNode); i {
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
		file_onebot_forward_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ForwardContent); i {
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
		file_onebot_forward_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ForwardChain); i {
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
		file_onebot_forward_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ForwardMessage); i {
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
	file_onebot_forward_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*ForwardContent_MessageNode)(nil),
		(*ForwardContent_ForwardNode)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_onebot_forward_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_onebot_forward_proto_goTypes,
		DependencyIndexes: file_onebot_forward_proto_depIdxs,
		MessageInfos:      file_onebot_forward_proto_msgTypes,
	}.Build()
	File_onebot_forward_proto = out.File
	file_onebot_forward_proto_rawDesc = nil
	file_onebot_forward_proto_goTypes = nil
	file_onebot_forward_proto_depIdxs = nil
}