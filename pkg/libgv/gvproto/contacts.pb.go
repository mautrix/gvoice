// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.21.12
// source: contacts.proto

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

type ContactDisplayInfo_Photo_Type int32

const (
	ContactDisplayInfo_Photo_UNKNOWN_TYPE   ContactDisplayInfo_Photo_Type = 0
	ContactDisplayInfo_Photo_MONOGRAM       ContactDisplayInfo_Photo_Type = 1
	ContactDisplayInfo_Photo_USER_SPECIFIED ContactDisplayInfo_Photo_Type = 3
)

// Enum value maps for ContactDisplayInfo_Photo_Type.
var (
	ContactDisplayInfo_Photo_Type_name = map[int32]string{
		0: "UNKNOWN_TYPE",
		1: "MONOGRAM",
		3: "USER_SPECIFIED",
	}
	ContactDisplayInfo_Photo_Type_value = map[string]int32{
		"UNKNOWN_TYPE":   0,
		"MONOGRAM":       1,
		"USER_SPECIFIED": 3,
	}
)

func (x ContactDisplayInfo_Photo_Type) Enum() *ContactDisplayInfo_Photo_Type {
	p := new(ContactDisplayInfo_Photo_Type)
	*p = x
	return p
}

func (x ContactDisplayInfo_Photo_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ContactDisplayInfo_Photo_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_contacts_proto_enumTypes[0].Descriptor()
}

func (ContactDisplayInfo_Photo_Type) Type() protoreflect.EnumType {
	return &file_contacts_proto_enumTypes[0]
}

func (x ContactDisplayInfo_Photo_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ContactDisplayInfo_Photo_Type.Descriptor instead.
func (ContactDisplayInfo_Photo_Type) EnumDescriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{4, 0, 0}
}

type ContactDisplayInfo_Name_SourceType int32

const (
	ContactDisplayInfo_Name_UNKNOWN_SOURCE_TYPE ContactDisplayInfo_Name_SourceType = 0
	ContactDisplayInfo_Name_CONTACT             ContactDisplayInfo_Name_SourceType = 2
)

// Enum value maps for ContactDisplayInfo_Name_SourceType.
var (
	ContactDisplayInfo_Name_SourceType_name = map[int32]string{
		0: "UNKNOWN_SOURCE_TYPE",
		2: "CONTACT",
	}
	ContactDisplayInfo_Name_SourceType_value = map[string]int32{
		"UNKNOWN_SOURCE_TYPE": 0,
		"CONTACT":             2,
	}
)

func (x ContactDisplayInfo_Name_SourceType) Enum() *ContactDisplayInfo_Name_SourceType {
	p := new(ContactDisplayInfo_Name_SourceType)
	*p = x
	return p
}

func (x ContactDisplayInfo_Name_SourceType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ContactDisplayInfo_Name_SourceType) Descriptor() protoreflect.EnumDescriptor {
	return file_contacts_proto_enumTypes[1].Descriptor()
}

func (ContactDisplayInfo_Name_SourceType) Type() protoreflect.EnumType {
	return &file_contacts_proto_enumTypes[1]
}

func (x ContactDisplayInfo_Name_SourceType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ContactDisplayInfo_Name_SourceType.Descriptor instead.
func (ContactDisplayInfo_Name_SourceType) EnumDescriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{4, 1, 0}
}

type StructuredPhone_ValidationResult int32

const (
	StructuredPhone_IS_POSSIBLE StructuredPhone_ValidationResult = 0
)

// Enum value maps for StructuredPhone_ValidationResult.
var (
	StructuredPhone_ValidationResult_name = map[int32]string{
		0: "IS_POSSIBLE",
	}
	StructuredPhone_ValidationResult_value = map[string]int32{
		"IS_POSSIBLE": 0,
	}
)

func (x StructuredPhone_ValidationResult) Enum() *StructuredPhone_ValidationResult {
	p := new(StructuredPhone_ValidationResult)
	*p = x
	return p
}

func (x StructuredPhone_ValidationResult) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StructuredPhone_ValidationResult) Descriptor() protoreflect.EnumDescriptor {
	return file_contacts_proto_enumTypes[2].Descriptor()
}

func (StructuredPhone_ValidationResult) Type() protoreflect.EnumType {
	return &file_contacts_proto_enumTypes[2]
}

func (x StructuredPhone_ValidationResult) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StructuredPhone_ValidationResult.Descriptor instead.
func (StructuredPhone_ValidationResult) EnumDescriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{5, 0}
}

type ContactID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Phone string `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty"`
}

func (x *ContactID) Reset() {
	*x = ContactID{}
	mi := &file_contacts_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContactID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContactID) ProtoMessage() {}

func (x *ContactID) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContactID.ProtoReflect.Descriptor instead.
func (*ContactID) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{0}
}

func (x *ContactID) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

type ContactSourceID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProfileID string `protobuf:"bytes,1,opt,name=profileID,proto3" json:"profileID,omitempty"`
	ContactID string `protobuf:"bytes,2,opt,name=contactID,proto3" json:"contactID,omitempty"`
}

func (x *ContactSourceID) Reset() {
	*x = ContactSourceID{}
	mi := &file_contacts_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContactSourceID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContactSourceID) ProtoMessage() {}

func (x *ContactSourceID) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContactSourceID.ProtoReflect.Descriptor instead.
func (*ContactSourceID) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{1}
}

func (x *ContactSourceID) GetProfileID() string {
	if x != nil {
		return x.ProfileID
	}
	return ""
}

func (x *ContactSourceID) GetContactID() string {
	if x != nil {
		return x.ContactID
	}
	return ""
}

type ContactPhone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DisplayValue   string `protobuf:"bytes,1,opt,name=displayValue,proto3" json:"displayValue,omitempty"`
	CanonicalValue string `protobuf:"bytes,2,opt,name=canonicalValue,proto3" json:"canonicalValue,omitempty"`
	OriginalValue  string `protobuf:"bytes,3,opt,name=originalValue,proto3" json:"originalValue,omitempty"`
}

func (x *ContactPhone) Reset() {
	*x = ContactPhone{}
	mi := &file_contacts_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContactPhone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContactPhone) ProtoMessage() {}

func (x *ContactPhone) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContactPhone.ProtoReflect.Descriptor instead.
func (*ContactPhone) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{2}
}

func (x *ContactPhone) GetDisplayValue() string {
	if x != nil {
		return x.DisplayValue
	}
	return ""
}

func (x *ContactPhone) GetCanonicalValue() string {
	if x != nil {
		return x.CanonicalValue
	}
	return ""
}

func (x *ContactPhone) GetOriginalValue() string {
	if x != nil {
		return x.OriginalValue
	}
	return ""
}

type ContactEmail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *ContactEmail) Reset() {
	*x = ContactEmail{}
	mi := &file_contacts_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContactEmail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContactEmail) ProtoMessage() {}

func (x *ContactEmail) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContactEmail.ProtoReflect.Descriptor instead.
func (*ContactEmail) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{3}
}

func (x *ContactEmail) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type ContactDisplayInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Photo    *ContactDisplayInfo_Photo `protobuf:"bytes,1,opt,name=photo,proto3" json:"photo,omitempty"`
	Name     *ContactDisplayInfo_Name  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Affinity *Affinity                 `protobuf:"bytes,3,opt,name=affinity,proto3" json:"affinity,omitempty"`
	Primary  bool                      `protobuf:"varint,4,opt,name=primary,proto3" json:"primary,omitempty"`
}

func (x *ContactDisplayInfo) Reset() {
	*x = ContactDisplayInfo{}
	mi := &file_contacts_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContactDisplayInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContactDisplayInfo) ProtoMessage() {}

func (x *ContactDisplayInfo) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContactDisplayInfo.ProtoReflect.Descriptor instead.
func (*ContactDisplayInfo) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{4}
}

func (x *ContactDisplayInfo) GetPhoto() *ContactDisplayInfo_Photo {
	if x != nil {
		return x.Photo
	}
	return nil
}

func (x *ContactDisplayInfo) GetName() *ContactDisplayInfo_Name {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *ContactDisplayInfo) GetAffinity() *Affinity {
	if x != nil {
		return x.Affinity
	}
	return nil
}

func (x *ContactDisplayInfo) GetPrimary() bool {
	if x != nil {
		return x.Primary
	}
	return false
}

type StructuredPhone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Phone *StructuredPhone_Phone `protobuf:"bytes,1,opt,name=phone,proto3" json:"phone,omitempty"`
	Type  string                 `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *StructuredPhone) Reset() {
	*x = StructuredPhone{}
	mi := &file_contacts_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StructuredPhone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StructuredPhone) ProtoMessage() {}

func (x *StructuredPhone) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StructuredPhone.ProtoReflect.Descriptor instead.
func (*StructuredPhone) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{5}
}

func (x *StructuredPhone) GetPhone() *StructuredPhone_Phone {
	if x != nil {
		return x.Phone
	}
	return nil
}

func (x *StructuredPhone) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type PhoneExtendedData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StructuredPhone *StructuredPhone `protobuf:"bytes,1,opt,name=structuredPhone,proto3" json:"structuredPhone,omitempty"`
}

func (x *PhoneExtendedData) Reset() {
	*x = PhoneExtendedData{}
	mi := &file_contacts_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PhoneExtendedData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PhoneExtendedData) ProtoMessage() {}

func (x *PhoneExtendedData) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PhoneExtendedData.ProtoReflect.Descriptor instead.
func (*PhoneExtendedData) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{6}
}

func (x *PhoneExtendedData) GetStructuredPhone() *StructuredPhone {
	if x != nil {
		return x.StructuredPhone
	}
	return nil
}

type PeopleStackFieldExtendedData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PhoneExtendedData *PhoneExtendedData `protobuf:"bytes,2,opt,name=phoneExtendedData,proto3" json:"phoneExtendedData,omitempty"`
}

func (x *PeopleStackFieldExtendedData) Reset() {
	*x = PeopleStackFieldExtendedData{}
	mi := &file_contacts_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PeopleStackFieldExtendedData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeopleStackFieldExtendedData) ProtoMessage() {}

func (x *PeopleStackFieldExtendedData) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PeopleStackFieldExtendedData.ProtoReflect.Descriptor instead.
func (*PeopleStackFieldExtendedData) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{7}
}

func (x *PeopleStackFieldExtendedData) GetPhoneExtendedData() *PhoneExtendedData {
	if x != nil {
		return x.PhoneExtendedData
	}
	return nil
}

type ContactMethod struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DisplayInfo  *ContactDisplayInfo           `protobuf:"bytes,1,opt,name=displayInfo,proto3" json:"displayInfo,omitempty"`
	Email        *ContactEmail                 `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Phone        *ContactPhone                 `protobuf:"bytes,3,opt,name=phone,proto3" json:"phone,omitempty"`
	SourceIDs    []*ContactSourceID            `protobuf:"bytes,5,rep,name=sourceIDs,proto3" json:"sourceIDs,omitempty"`
	IsPersonal   bool                          `protobuf:"varint,8,opt,name=isPersonal,proto3" json:"isPersonal,omitempty"`
	ExtendedData *PeopleStackFieldExtendedData `protobuf:"bytes,9,opt,name=extendedData,proto3" json:"extendedData,omitempty"`
	TypeLabel    string                        `protobuf:"bytes,10,opt,name=typeLabel,proto3" json:"typeLabel,omitempty"`
}

func (x *ContactMethod) Reset() {
	*x = ContactMethod{}
	mi := &file_contacts_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContactMethod) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContactMethod) ProtoMessage() {}

func (x *ContactMethod) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContactMethod.ProtoReflect.Descriptor instead.
func (*ContactMethod) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{8}
}

func (x *ContactMethod) GetDisplayInfo() *ContactDisplayInfo {
	if x != nil {
		return x.DisplayInfo
	}
	return nil
}

func (x *ContactMethod) GetEmail() *ContactEmail {
	if x != nil {
		return x.Email
	}
	return nil
}

func (x *ContactMethod) GetPhone() *ContactPhone {
	if x != nil {
		return x.Phone
	}
	return nil
}

func (x *ContactMethod) GetSourceIDs() []*ContactSourceID {
	if x != nil {
		return x.SourceIDs
	}
	return nil
}

func (x *ContactMethod) GetIsPersonal() bool {
	if x != nil {
		return x.IsPersonal
	}
	return false
}

func (x *ContactMethod) GetExtendedData() *PeopleStackFieldExtendedData {
	if x != nil {
		return x.ExtendedData
	}
	return nil
}

func (x *ContactMethod) GetTypeLabel() string {
	if x != nil {
		return x.TypeLabel
	}
	return ""
}

type Affinity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LoggingID string  `protobuf:"bytes,1,opt,name=loggingID,proto3" json:"loggingID,omitempty"`
	Value     float64 `protobuf:"fixed64,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Affinity) Reset() {
	*x = Affinity{}
	mi := &file_contacts_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Affinity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Affinity) ProtoMessage() {}

func (x *Affinity) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Affinity.ProtoReflect.Descriptor instead.
func (*Affinity) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{9}
}

func (x *Affinity) GetLoggingID() string {
	if x != nil {
		return x.LoggingID
	}
	return ""
}

func (x *Affinity) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type Person struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContactMethods []*ContactMethod `protobuf:"bytes,1,rep,name=contactMethods,proto3" json:"contactMethods,omitempty"`
	Affinity       *Affinity        `protobuf:"bytes,2,opt,name=affinity,proto3" json:"affinity,omitempty"`
}

func (x *Person) Reset() {
	*x = Person{}
	mi := &file_contacts_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Person) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Person) ProtoMessage() {}

func (x *Person) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Person.ProtoReflect.Descriptor instead.
func (*Person) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{10}
}

func (x *Person) GetContactMethods() []*ContactMethod {
	if x != nil {
		return x.ContactMethods
	}
	return nil
}

func (x *Person) GetAffinity() *Affinity {
	if x != nil {
		return x.Affinity
	}
	return nil
}

type PersonWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Person *Person `protobuf:"bytes,1,opt,name=person,proto3" json:"person,omitempty"`
}

func (x *PersonWrapper) Reset() {
	*x = PersonWrapper{}
	mi := &file_contacts_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PersonWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PersonWrapper) ProtoMessage() {}

func (x *PersonWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PersonWrapper.ProtoReflect.Descriptor instead.
func (*PersonWrapper) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{11}
}

func (x *PersonWrapper) GetPerson() *Person {
	if x != nil {
		return x.Person
	}
	return nil
}

type ContactDisplayInfo_Photo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	URL  string                        `protobuf:"bytes,1,opt,name=URL,proto3" json:"URL,omitempty"`
	Type ContactDisplayInfo_Photo_Type `protobuf:"varint,2,opt,name=type,proto3,enum=contacts.ContactDisplayInfo_Photo_Type" json:"type,omitempty"`
}

func (x *ContactDisplayInfo_Photo) Reset() {
	*x = ContactDisplayInfo_Photo{}
	mi := &file_contacts_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContactDisplayInfo_Photo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContactDisplayInfo_Photo) ProtoMessage() {}

func (x *ContactDisplayInfo_Photo) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContactDisplayInfo_Photo.ProtoReflect.Descriptor instead.
func (*ContactDisplayInfo_Photo) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{4, 0}
}

func (x *ContactDisplayInfo_Photo) GetURL() string {
	if x != nil {
		return x.URL
	}
	return ""
}

func (x *ContactDisplayInfo_Photo) GetType() ContactDisplayInfo_Photo_Type {
	if x != nil {
		return x.Type
	}
	return ContactDisplayInfo_Photo_UNKNOWN_TYPE
}

type ContactDisplayInfo_Name struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value      string                             `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	GivenName  string                             `protobuf:"bytes,2,opt,name=givenName,proto3" json:"givenName,omitempty"`
	SourceType ContactDisplayInfo_Name_SourceType `protobuf:"varint,4,opt,name=sourceType,proto3,enum=contacts.ContactDisplayInfo_Name_SourceType" json:"sourceType,omitempty"`
}

func (x *ContactDisplayInfo_Name) Reset() {
	*x = ContactDisplayInfo_Name{}
	mi := &file_contacts_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContactDisplayInfo_Name) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContactDisplayInfo_Name) ProtoMessage() {}

func (x *ContactDisplayInfo_Name) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContactDisplayInfo_Name.ProtoReflect.Descriptor instead.
func (*ContactDisplayInfo_Name) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{4, 1}
}

func (x *ContactDisplayInfo_Name) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *ContactDisplayInfo_Name) GetGivenName() string {
	if x != nil {
		return x.GivenName
	}
	return ""
}

func (x *ContactDisplayInfo_Name) GetSourceType() ContactDisplayInfo_Name_SourceType {
	if x != nil {
		return x.SourceType
	}
	return ContactDisplayInfo_Name_UNKNOWN_SOURCE_TYPE
}

type StructuredPhone_I18NData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NationalNumber      string                           `protobuf:"bytes,1,opt,name=nationalNumber,proto3" json:"nationalNumber,omitempty"`
	InternationalNumber string                           `protobuf:"bytes,2,opt,name=internationalNumber,proto3" json:"internationalNumber,omitempty"`
	CountryCode         int32                            `protobuf:"varint,3,opt,name=countryCode,proto3" json:"countryCode,omitempty"`
	RegionCode          string                           `protobuf:"bytes,4,opt,name=regionCode,proto3" json:"regionCode,omitempty"`
	IsValid             bool                             `protobuf:"varint,5,opt,name=isValid,proto3" json:"isValid,omitempty"`
	ValidationResult    StructuredPhone_ValidationResult `protobuf:"varint,6,opt,name=validationResult,proto3,enum=contacts.StructuredPhone_ValidationResult" json:"validationResult,omitempty"`
}

func (x *StructuredPhone_I18NData) Reset() {
	*x = StructuredPhone_I18NData{}
	mi := &file_contacts_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StructuredPhone_I18NData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StructuredPhone_I18NData) ProtoMessage() {}

func (x *StructuredPhone_I18NData) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StructuredPhone_I18NData.ProtoReflect.Descriptor instead.
func (*StructuredPhone_I18NData) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{5, 0}
}

func (x *StructuredPhone_I18NData) GetNationalNumber() string {
	if x != nil {
		return x.NationalNumber
	}
	return ""
}

func (x *StructuredPhone_I18NData) GetInternationalNumber() string {
	if x != nil {
		return x.InternationalNumber
	}
	return ""
}

func (x *StructuredPhone_I18NData) GetCountryCode() int32 {
	if x != nil {
		return x.CountryCode
	}
	return 0
}

func (x *StructuredPhone_I18NData) GetRegionCode() string {
	if x != nil {
		return x.RegionCode
	}
	return ""
}

func (x *StructuredPhone_I18NData) GetIsValid() bool {
	if x != nil {
		return x.IsValid
	}
	return false
}

func (x *StructuredPhone_I18NData) GetValidationResult() StructuredPhone_ValidationResult {
	if x != nil {
		return x.ValidationResult
	}
	return StructuredPhone_IS_POSSIBLE
}

type StructuredPhone_Phone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	E164     string                    `protobuf:"bytes,1,opt,name=e164,proto3" json:"e164,omitempty"`
	I18NData *StructuredPhone_I18NData `protobuf:"bytes,2,opt,name=i18nData,proto3" json:"i18nData,omitempty"`
}

func (x *StructuredPhone_Phone) Reset() {
	*x = StructuredPhone_Phone{}
	mi := &file_contacts_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StructuredPhone_Phone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StructuredPhone_Phone) ProtoMessage() {}

func (x *StructuredPhone_Phone) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StructuredPhone_Phone.ProtoReflect.Descriptor instead.
func (*StructuredPhone_Phone) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{5, 1}
}

func (x *StructuredPhone_Phone) GetE164() string {
	if x != nil {
		return x.E164
	}
	return ""
}

func (x *StructuredPhone_Phone) GetI18NData() *StructuredPhone_I18NData {
	if x != nil {
		return x.I18NData
	}
	return nil
}

var File_contacts_proto protoreflect.FileDescriptor

//go:embed contacts.pb.raw
var file_contacts_proto_rawDesc []byte

var (
	file_contacts_proto_rawDescOnce sync.Once
	file_contacts_proto_rawDescData = file_contacts_proto_rawDesc
)

func file_contacts_proto_rawDescGZIP() []byte {
	file_contacts_proto_rawDescOnce.Do(func() {
		file_contacts_proto_rawDescData = protoimpl.X.CompressGZIP(file_contacts_proto_rawDescData)
	})
	return file_contacts_proto_rawDescData
}

var file_contacts_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_contacts_proto_msgTypes = make([]protoimpl.MessageInfo, 16)
var file_contacts_proto_goTypes = []any{
	(ContactDisplayInfo_Photo_Type)(0),      // 0: contacts.ContactDisplayInfo.Photo.Type
	(ContactDisplayInfo_Name_SourceType)(0), // 1: contacts.ContactDisplayInfo.Name.SourceType
	(StructuredPhone_ValidationResult)(0),   // 2: contacts.StructuredPhone.ValidationResult
	(*ContactID)(nil),                       // 3: contacts.ContactID
	(*ContactSourceID)(nil),                 // 4: contacts.ContactSourceID
	(*ContactPhone)(nil),                    // 5: contacts.ContactPhone
	(*ContactEmail)(nil),                    // 6: contacts.ContactEmail
	(*ContactDisplayInfo)(nil),              // 7: contacts.ContactDisplayInfo
	(*StructuredPhone)(nil),                 // 8: contacts.StructuredPhone
	(*PhoneExtendedData)(nil),               // 9: contacts.PhoneExtendedData
	(*PeopleStackFieldExtendedData)(nil),    // 10: contacts.PeopleStackFieldExtendedData
	(*ContactMethod)(nil),                   // 11: contacts.ContactMethod
	(*Affinity)(nil),                        // 12: contacts.Affinity
	(*Person)(nil),                          // 13: contacts.Person
	(*PersonWrapper)(nil),                   // 14: contacts.PersonWrapper
	(*ContactDisplayInfo_Photo)(nil),        // 15: contacts.ContactDisplayInfo.Photo
	(*ContactDisplayInfo_Name)(nil),         // 16: contacts.ContactDisplayInfo.Name
	(*StructuredPhone_I18NData)(nil),        // 17: contacts.StructuredPhone.I18nData
	(*StructuredPhone_Phone)(nil),           // 18: contacts.StructuredPhone.Phone
}
var file_contacts_proto_depIdxs = []int32{
	15, // 0: contacts.ContactDisplayInfo.photo:type_name -> contacts.ContactDisplayInfo.Photo
	16, // 1: contacts.ContactDisplayInfo.name:type_name -> contacts.ContactDisplayInfo.Name
	12, // 2: contacts.ContactDisplayInfo.affinity:type_name -> contacts.Affinity
	18, // 3: contacts.StructuredPhone.phone:type_name -> contacts.StructuredPhone.Phone
	8,  // 4: contacts.PhoneExtendedData.structuredPhone:type_name -> contacts.StructuredPhone
	9,  // 5: contacts.PeopleStackFieldExtendedData.phoneExtendedData:type_name -> contacts.PhoneExtendedData
	7,  // 6: contacts.ContactMethod.displayInfo:type_name -> contacts.ContactDisplayInfo
	6,  // 7: contacts.ContactMethod.email:type_name -> contacts.ContactEmail
	5,  // 8: contacts.ContactMethod.phone:type_name -> contacts.ContactPhone
	4,  // 9: contacts.ContactMethod.sourceIDs:type_name -> contacts.ContactSourceID
	10, // 10: contacts.ContactMethod.extendedData:type_name -> contacts.PeopleStackFieldExtendedData
	11, // 11: contacts.Person.contactMethods:type_name -> contacts.ContactMethod
	12, // 12: contacts.Person.affinity:type_name -> contacts.Affinity
	13, // 13: contacts.PersonWrapper.person:type_name -> contacts.Person
	0,  // 14: contacts.ContactDisplayInfo.Photo.type:type_name -> contacts.ContactDisplayInfo.Photo.Type
	1,  // 15: contacts.ContactDisplayInfo.Name.sourceType:type_name -> contacts.ContactDisplayInfo.Name.SourceType
	2,  // 16: contacts.StructuredPhone.I18nData.validationResult:type_name -> contacts.StructuredPhone.ValidationResult
	17, // 17: contacts.StructuredPhone.Phone.i18nData:type_name -> contacts.StructuredPhone.I18nData
	18, // [18:18] is the sub-list for method output_type
	18, // [18:18] is the sub-list for method input_type
	18, // [18:18] is the sub-list for extension type_name
	18, // [18:18] is the sub-list for extension extendee
	0,  // [0:18] is the sub-list for field type_name
}

func init() { file_contacts_proto_init() }
func file_contacts_proto_init() {
	if File_contacts_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_contacts_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   16,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_contacts_proto_goTypes,
		DependencyIndexes: file_contacts_proto_depIdxs,
		EnumInfos:         file_contacts_proto_enumTypes,
		MessageInfos:      file_contacts_proto_msgTypes,
	}.Build()
	File_contacts_proto = out.File
	file_contacts_proto_rawDesc = nil
	file_contacts_proto_goTypes = nil
	file_contacts_proto_depIdxs = nil
}
