// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.21.12
// source: webchannel.proto

package gvproto

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	_ "embed"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RespChooseServer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GSessionID string `protobuf:"bytes,1,opt,name=gSessionID,proto3" json:"gSessionID,omitempty"`
	UnknownInt int32  `protobuf:"varint,2,opt,name=unknownInt,proto3" json:"unknownInt,omitempty"` // 3
}

func (x *RespChooseServer) Reset() {
	*x = RespChooseServer{}
	mi := &file_webchannel_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RespChooseServer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RespChooseServer) ProtoMessage() {}

func (x *RespChooseServer) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RespChooseServer.ProtoReflect.Descriptor instead.
func (*RespChooseServer) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{0}
}

func (x *RespChooseServer) GetGSessionID() string {
	if x != nil {
		return x.GSessionID
	}
	return ""
}

func (x *RespChooseServer) GetUnknownInt() int32 {
	if x != nil {
		return x.UnknownInt
	}
	return 0
}

type RespCreateChannel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *WebChannelSessionData `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *RespCreateChannel) Reset() {
	*x = RespCreateChannel{}
	mi := &file_webchannel_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RespCreateChannel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RespCreateChannel) ProtoMessage() {}

func (x *RespCreateChannel) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RespCreateChannel.ProtoReflect.Descriptor instead.
func (*RespCreateChannel) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{1}
}

func (x *RespCreateChannel) GetData() *WebChannelSessionData {
	if x != nil {
		return x.Data
	}
	return nil
}

type WebChannelSessionData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field1  int32              `protobuf:"varint,1,opt,name=field1,proto3" json:"field1,omitempty"` // 0
	Session *WebChannelSession `protobuf:"bytes,2,opt,name=session,proto3" json:"session,omitempty"`
}

func (x *WebChannelSessionData) Reset() {
	*x = WebChannelSessionData{}
	mi := &file_webchannel_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WebChannelSessionData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebChannelSessionData) ProtoMessage() {}

func (x *WebChannelSessionData) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebChannelSessionData.ProtoReflect.Descriptor instead.
func (*WebChannelSessionData) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{2}
}

func (x *WebChannelSessionData) GetField1() int32 {
	if x != nil {
		return x.Field1
	}
	return 0
}

func (x *WebChannelSessionData) GetSession() *WebChannelSession {
	if x != nil {
		return x.Session
	}
	return nil
}

type WebChannelSession struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field1       string `protobuf:"bytes,1,opt,name=field1,proto3" json:"field1,omitempty"` // "c"?
	SessionID    string `protobuf:"bytes,2,opt,name=sessionID,proto3" json:"sessionID,omitempty"`
	Field3       string `protobuf:"bytes,3,opt,name=field3,proto3" json:"field3,omitempty"`              // empty string
	Field4       uint32 `protobuf:"varint,4,opt,name=field4,proto3" json:"field4,omitempty"`             // 8, version?
	Field5       uint32 `protobuf:"varint,5,opt,name=field5,proto3" json:"field5,omitempty"`             // 12, ???
	PingInterval uint32 `protobuf:"varint,6,opt,name=pingInterval,proto3" json:"pingInterval,omitempty"` // 30000
}

func (x *WebChannelSession) Reset() {
	*x = WebChannelSession{}
	mi := &file_webchannel_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WebChannelSession) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebChannelSession) ProtoMessage() {}

func (x *WebChannelSession) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebChannelSession.ProtoReflect.Descriptor instead.
func (*WebChannelSession) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{3}
}

func (x *WebChannelSession) GetField1() string {
	if x != nil {
		return x.Field1
	}
	return ""
}

func (x *WebChannelSession) GetSessionID() string {
	if x != nil {
		return x.SessionID
	}
	return ""
}

func (x *WebChannelSession) GetField3() string {
	if x != nil {
		return x.Field3
	}
	return ""
}

func (x *WebChannelSession) GetField4() uint32 {
	if x != nil {
		return x.Field4
	}
	return 0
}

func (x *WebChannelSession) GetField5() uint32 {
	if x != nil {
		return x.Field5
	}
	return 0
}

func (x *WebChannelSession) GetPingInterval() uint32 {
	if x != nil {
		return x.PingInterval
	}
	return 0
}

type WebChannelEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ArrayID     uint64                        `protobuf:"varint,1,opt,name=arrayID,proto3" json:"arrayID,omitempty"`
	DataWrapper []*WebChannelEventDataWrapper `protobuf:"bytes,2,rep,name=data_wrapper,json=dataWrapper,proto3" json:"data_wrapper,omitempty"`
}

func (x *WebChannelEvent) Reset() {
	*x = WebChannelEvent{}
	mi := &file_webchannel_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WebChannelEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebChannelEvent) ProtoMessage() {}

func (x *WebChannelEvent) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebChannelEvent.ProtoReflect.Descriptor instead.
func (*WebChannelEvent) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{4}
}

func (x *WebChannelEvent) GetArrayID() uint64 {
	if x != nil {
		return x.ArrayID
	}
	return 0
}

func (x *WebChannelEvent) GetDataWrapper() []*WebChannelEventDataWrapper {
	if x != nil {
		return x.DataWrapper
	}
	return nil
}

type WebChannelNoopEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ArrayID uint64   `protobuf:"varint,1,opt,name=arrayID,proto3" json:"arrayID,omitempty"`
	Noop    []string `protobuf:"bytes,2,rep,name=noop,proto3" json:"noop,omitempty"`
}

func (x *WebChannelNoopEvent) Reset() {
	*x = WebChannelNoopEvent{}
	mi := &file_webchannel_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WebChannelNoopEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebChannelNoopEvent) ProtoMessage() {}

func (x *WebChannelNoopEvent) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebChannelNoopEvent.ProtoReflect.Descriptor instead.
func (*WebChannelNoopEvent) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{5}
}

func (x *WebChannelNoopEvent) GetArrayID() uint64 {
	if x != nil {
		return x.ArrayID
	}
	return 0
}

func (x *WebChannelNoopEvent) GetNoop() []string {
	if x != nil {
		return x.Noop
	}
	return nil
}

type WebChannelEventDataWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// TODO this may also be the string "noop"
	Data    []*WebChannelEventData              `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	AltData *WebChannelEventDataWrapper_AltData `protobuf:"bytes,2,opt,name=altData,proto3" json:"altData,omitempty"`
}

func (x *WebChannelEventDataWrapper) Reset() {
	*x = WebChannelEventDataWrapper{}
	mi := &file_webchannel_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WebChannelEventDataWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebChannelEventDataWrapper) ProtoMessage() {}

func (x *WebChannelEventDataWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebChannelEventDataWrapper.ProtoReflect.Descriptor instead.
func (*WebChannelEventDataWrapper) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{6}
}

func (x *WebChannelEventDataWrapper) GetData() []*WebChannelEventData {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *WebChannelEventDataWrapper) GetAltData() *WebChannelEventDataWrapper_AltData {
	if x != nil {
		return x.AltData
	}
	return nil
}

type WebChannelEventData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventSource string `protobuf:"bytes,1,opt,name=eventSource,proto3" json:"eventSource,omitempty"` // Seems to be 1, 2, 3, 4, 5, 6 or 9. Same as the create channel reqX___data__'s
	Event       *Event `protobuf:"bytes,2,opt,name=event,proto3" json:"event,omitempty"`
}

func (x *WebChannelEventData) Reset() {
	*x = WebChannelEventData{}
	mi := &file_webchannel_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WebChannelEventData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebChannelEventData) ProtoMessage() {}

func (x *WebChannelEventData) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebChannelEventData.ProtoReflect.Descriptor instead.
func (*WebChannelEventData) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{7}
}

func (x *WebChannelEventData) GetEventSource() string {
	if x != nil {
		return x.EventSource
	}
	return ""
}

func (x *WebChannelEventData) GetEvent() *Event {
	if x != nil {
		return x.Event
	}
	return nil
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sub1 *EventSub1 `protobuf:"bytes,1,opt,name=sub1,proto3" json:"sub1,omitempty"`
	Sub2 *EventSub2 `protobuf:"bytes,2,opt,name=sub2,proto3" json:"sub2,omitempty"`
	Sub3 *EventSub3 `protobuf:"bytes,3,opt,name=sub3,proto3" json:"sub3,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	mi := &file_webchannel_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{8}
}

func (x *Event) GetSub1() *EventSub1 {
	if x != nil {
		return x.Sub1
	}
	return nil
}

func (x *Event) GetSub2() *EventSub2 {
	if x != nil {
		return x.Sub2
	}
	return nil
}

func (x *Event) GetSub3() *EventSub3 {
	if x != nil {
		return x.Sub3
	}
	return nil
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_webchannel_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{9}
}

type EventSub1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Unknown1 *Empty `protobuf:"bytes,1,opt,name=unknown1,proto3" json:"unknown1,omitempty"`
}

func (x *EventSub1) Reset() {
	*x = EventSub1{}
	mi := &file_webchannel_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventSub1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventSub1) ProtoMessage() {}

func (x *EventSub1) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventSub1.ProtoReflect.Descriptor instead.
func (*EventSub1) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{10}
}

func (x *EventSub1) GetUnknown1() *Empty {
	if x != nil {
		return x.Unknown1
	}
	return nil
}

type EventSub2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*EventSub2Data `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *EventSub2) Reset() {
	*x = EventSub2{}
	mi := &file_webchannel_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventSub2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventSub2) ProtoMessage() {}

func (x *EventSub2) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventSub2.ProtoReflect.Descriptor instead.
func (*EventSub2) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{11}
}

func (x *EventSub2) GetData() []*EventSub2Data {
	if x != nil {
		return x.Data
	}
	return nil
}

type EventSub2Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UnknownBytes           []byte                    `protobuf:"bytes,2,opt,name=unknownBytes,proto3" json:"unknownBytes,omitempty"`
	UnknownTimestampMillis string                    `protobuf:"bytes,3,opt,name=unknownTimestampMillis,proto3" json:"unknownTimestampMillis,omitempty"`
	UnknownTimestampMicros string                    `protobuf:"bytes,4,opt,name=unknownTimestampMicros,proto3" json:"unknownTimestampMicros,omitempty"`
	UnknownNestedData      *EventSub2Data_NestedData `protobuf:"bytes,5,opt,name=unknownNestedData,proto3" json:"unknownNestedData,omitempty"` // 5 has a nested value of some kind
}

func (x *EventSub2Data) Reset() {
	*x = EventSub2Data{}
	mi := &file_webchannel_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventSub2Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventSub2Data) ProtoMessage() {}

func (x *EventSub2Data) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventSub2Data.ProtoReflect.Descriptor instead.
func (*EventSub2Data) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{12}
}

func (x *EventSub2Data) GetUnknownBytes() []byte {
	if x != nil {
		return x.UnknownBytes
	}
	return nil
}

func (x *EventSub2Data) GetUnknownTimestampMillis() string {
	if x != nil {
		return x.UnknownTimestampMillis
	}
	return ""
}

func (x *EventSub2Data) GetUnknownTimestampMicros() string {
	if x != nil {
		return x.UnknownTimestampMicros
	}
	return ""
}

func (x *EventSub2Data) GetUnknownNestedData() *EventSub2Data_NestedData {
	if x != nil {
		return x.UnknownNestedData
	}
	return nil
}

type EventSub3 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UnknownTimestamp string `protobuf:"bytes,1,opt,name=unknownTimestamp,proto3" json:"unknownTimestamp,omitempty"`
}

func (x *EventSub3) Reset() {
	*x = EventSub3{}
	mi := &file_webchannel_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventSub3) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventSub3) ProtoMessage() {}

func (x *EventSub3) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventSub3.ProtoReflect.Descriptor instead.
func (*EventSub3) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{13}
}

func (x *EventSub3) GetUnknownTimestamp() string {
	if x != nil {
		return x.UnknownTimestamp
	}
	return ""
}

type WebChannelEventDataWrapper_AltData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reconnect bool `protobuf:"varint,1,opt,name=reconnect,proto3" json:"reconnect,omitempty"`
}

func (x *WebChannelEventDataWrapper_AltData) Reset() {
	*x = WebChannelEventDataWrapper_AltData{}
	mi := &file_webchannel_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WebChannelEventDataWrapper_AltData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebChannelEventDataWrapper_AltData) ProtoMessage() {}

func (x *WebChannelEventDataWrapper_AltData) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebChannelEventDataWrapper_AltData.ProtoReflect.Descriptor instead.
func (*WebChannelEventDataWrapper_AltData) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{6, 0}
}

func (x *WebChannelEventDataWrapper_AltData) GetReconnect() bool {
	if x != nil {
		return x.Reconnect
	}
	return false
}

type EventSub2Data_NestedData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ts1 string `protobuf:"bytes,1,opt,name=ts1,proto3" json:"ts1,omitempty"`
	Ts2 string `protobuf:"bytes,2,opt,name=ts2,proto3" json:"ts2,omitempty"`
	Ts3 string `protobuf:"bytes,3,opt,name=ts3,proto3" json:"ts3,omitempty"`
	Ts4 string `protobuf:"bytes,4,opt,name=ts4,proto3" json:"ts4,omitempty"`
	Ts5 string `protobuf:"bytes,5,opt,name=ts5,proto3" json:"ts5,omitempty"`
	Ts6 string `protobuf:"bytes,6,opt,name=ts6,proto3" json:"ts6,omitempty"`
	Ts7 string `protobuf:"bytes,7,opt,name=ts7,proto3" json:"ts7,omitempty"`
	Ts9 string `protobuf:"bytes,9,opt,name=ts9,proto3" json:"ts9,omitempty"`
}

func (x *EventSub2Data_NestedData) Reset() {
	*x = EventSub2Data_NestedData{}
	mi := &file_webchannel_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventSub2Data_NestedData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventSub2Data_NestedData) ProtoMessage() {}

func (x *EventSub2Data_NestedData) ProtoReflect() protoreflect.Message {
	mi := &file_webchannel_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventSub2Data_NestedData.ProtoReflect.Descriptor instead.
func (*EventSub2Data_NestedData) Descriptor() ([]byte, []int) {
	return file_webchannel_proto_rawDescGZIP(), []int{12, 0}
}

func (x *EventSub2Data_NestedData) GetTs1() string {
	if x != nil {
		return x.Ts1
	}
	return ""
}

func (x *EventSub2Data_NestedData) GetTs2() string {
	if x != nil {
		return x.Ts2
	}
	return ""
}

func (x *EventSub2Data_NestedData) GetTs3() string {
	if x != nil {
		return x.Ts3
	}
	return ""
}

func (x *EventSub2Data_NestedData) GetTs4() string {
	if x != nil {
		return x.Ts4
	}
	return ""
}

func (x *EventSub2Data_NestedData) GetTs5() string {
	if x != nil {
		return x.Ts5
	}
	return ""
}

func (x *EventSub2Data_NestedData) GetTs6() string {
	if x != nil {
		return x.Ts6
	}
	return ""
}

func (x *EventSub2Data_NestedData) GetTs7() string {
	if x != nil {
		return x.Ts7
	}
	return ""
}

func (x *EventSub2Data_NestedData) GetTs9() string {
	if x != nil {
		return x.Ts9
	}
	return ""
}

var File_webchannel_proto protoreflect.FileDescriptor

//go:embed webchannel.pb.raw
var file_webchannel_proto_rawDesc []byte

var (
	file_webchannel_proto_rawDescOnce sync.Once
	file_webchannel_proto_rawDescData = file_webchannel_proto_rawDesc
)

func file_webchannel_proto_rawDescGZIP() []byte {
	file_webchannel_proto_rawDescOnce.Do(func() {
		file_webchannel_proto_rawDescData = protoimpl.X.CompressGZIP(file_webchannel_proto_rawDescData)
	})
	return file_webchannel_proto_rawDescData
}

var file_webchannel_proto_msgTypes = make([]protoimpl.MessageInfo, 16)
var file_webchannel_proto_goTypes = []any{
	(*RespChooseServer)(nil),                   // 0: webchannel.RespChooseServer
	(*RespCreateChannel)(nil),                  // 1: webchannel.RespCreateChannel
	(*WebChannelSessionData)(nil),              // 2: webchannel.WebChannelSessionData
	(*WebChannelSession)(nil),                  // 3: webchannel.WebChannelSession
	(*WebChannelEvent)(nil),                    // 4: webchannel.WebChannelEvent
	(*WebChannelNoopEvent)(nil),                // 5: webchannel.WebChannelNoopEvent
	(*WebChannelEventDataWrapper)(nil),         // 6: webchannel.WebChannelEventDataWrapper
	(*WebChannelEventData)(nil),                // 7: webchannel.WebChannelEventData
	(*Event)(nil),                              // 8: webchannel.Event
	(*Empty)(nil),                              // 9: webchannel.Empty
	(*EventSub1)(nil),                          // 10: webchannel.EventSub1
	(*EventSub2)(nil),                          // 11: webchannel.EventSub2
	(*EventSub2Data)(nil),                      // 12: webchannel.EventSub2Data
	(*EventSub3)(nil),                          // 13: webchannel.EventSub3
	(*WebChannelEventDataWrapper_AltData)(nil), // 14: webchannel.WebChannelEventDataWrapper.AltData
	(*EventSub2Data_NestedData)(nil),           // 15: webchannel.EventSub2Data.NestedData
}
var file_webchannel_proto_depIdxs = []int32{
	2,  // 0: webchannel.RespCreateChannel.data:type_name -> webchannel.WebChannelSessionData
	3,  // 1: webchannel.WebChannelSessionData.session:type_name -> webchannel.WebChannelSession
	6,  // 2: webchannel.WebChannelEvent.data_wrapper:type_name -> webchannel.WebChannelEventDataWrapper
	7,  // 3: webchannel.WebChannelEventDataWrapper.data:type_name -> webchannel.WebChannelEventData
	14, // 4: webchannel.WebChannelEventDataWrapper.altData:type_name -> webchannel.WebChannelEventDataWrapper.AltData
	8,  // 5: webchannel.WebChannelEventData.event:type_name -> webchannel.Event
	10, // 6: webchannel.Event.sub1:type_name -> webchannel.EventSub1
	11, // 7: webchannel.Event.sub2:type_name -> webchannel.EventSub2
	13, // 8: webchannel.Event.sub3:type_name -> webchannel.EventSub3
	9,  // 9: webchannel.EventSub1.unknown1:type_name -> webchannel.Empty
	12, // 10: webchannel.EventSub2.data:type_name -> webchannel.EventSub2Data
	15, // 11: webchannel.EventSub2Data.unknownNestedData:type_name -> webchannel.EventSub2Data.NestedData
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_webchannel_proto_init() }
func file_webchannel_proto_init() {
	if File_webchannel_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_webchannel_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   16,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_webchannel_proto_goTypes,
		DependencyIndexes: file_webchannel_proto_depIdxs,
		MessageInfos:      file_webchannel_proto_msgTypes,
	}.Build()
	File_webchannel_proto = out.File
	file_webchannel_proto_rawDesc = nil
	file_webchannel_proto_goTypes = nil
	file_webchannel_proto_depIdxs = nil
}
