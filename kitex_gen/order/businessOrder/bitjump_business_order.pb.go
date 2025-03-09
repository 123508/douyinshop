// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: order/bitjump_business_order.proto

package businessOrder

import (
	context "context"
	order_common "github.com/123508/douyinshop/kitex_gen/order/order_common"
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

type ConfirmReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId uint32 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	Status  int32  `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *ConfirmReq) Reset() {
	*x = ConfirmReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_bitjump_business_order_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfirmReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfirmReq) ProtoMessage() {}

func (x *ConfirmReq) ProtoReflect() protoreflect.Message {
	mi := &file_order_bitjump_business_order_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfirmReq.ProtoReflect.Descriptor instead.
func (*ConfirmReq) Descriptor() ([]byte, []int) {
	return file_order_bitjump_business_order_proto_rawDescGZIP(), []int{0}
}

func (x *ConfirmReq) GetOrderId() uint32 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

func (x *ConfirmReq) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

type RejectionReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId         uint32 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	RejectionReason string `protobuf:"bytes,2,opt,name=rejection_reason,json=rejectionReason,proto3" json:"rejection_reason,omitempty"`
}

func (x *RejectionReq) Reset() {
	*x = RejectionReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_bitjump_business_order_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RejectionReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RejectionReq) ProtoMessage() {}

func (x *RejectionReq) ProtoReflect() protoreflect.Message {
	mi := &file_order_bitjump_business_order_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RejectionReq.ProtoReflect.Descriptor instead.
func (*RejectionReq) Descriptor() ([]byte, []int) {
	return file_order_bitjump_business_order_proto_rawDescGZIP(), []int{1}
}

func (x *RejectionReq) GetOrderId() uint32 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

func (x *RejectionReq) GetRejectionReason() string {
	if x != nil {
		return x.RejectionReason
	}
	return ""
}

type DeliveryReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId uint32 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
}

func (x *DeliveryReq) Reset() {
	*x = DeliveryReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_bitjump_business_order_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeliveryReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeliveryReq) ProtoMessage() {}

func (x *DeliveryReq) ProtoReflect() protoreflect.Message {
	mi := &file_order_bitjump_business_order_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeliveryReq.ProtoReflect.Descriptor instead.
func (*DeliveryReq) Descriptor() ([]byte, []int) {
	return file_order_bitjump_business_order_proto_rawDescGZIP(), []int{2}
}

func (x *DeliveryReq) GetOrderId() uint32 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

type ReceiveReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId uint32 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
}

func (x *ReceiveReq) Reset() {
	*x = ReceiveReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_bitjump_business_order_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReceiveReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReceiveReq) ProtoMessage() {}

func (x *ReceiveReq) ProtoReflect() protoreflect.Message {
	mi := &file_order_bitjump_business_order_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReceiveReq.ProtoReflect.Descriptor instead.
func (*ReceiveReq) Descriptor() ([]byte, []int) {
	return file_order_bitjump_business_order_proto_rawDescGZIP(), []int{3}
}

func (x *ReceiveReq) GetOrderId() uint32 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

type GetOrderListReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShopId   uint32 `protobuf:"varint,1,opt,name=shop_id,json=shopId,proto3" json:"shop_id,omitempty"`
	Page     uint32 `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	PageSize uint32 `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
}

func (x *GetOrderListReq) Reset() {
	*x = GetOrderListReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_bitjump_business_order_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetOrderListReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrderListReq) ProtoMessage() {}

func (x *GetOrderListReq) ProtoReflect() protoreflect.Message {
	mi := &file_order_bitjump_business_order_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrderListReq.ProtoReflect.Descriptor instead.
func (*GetOrderListReq) Descriptor() ([]byte, []int) {
	return file_order_bitjump_business_order_proto_rawDescGZIP(), []int{4}
}

func (x *GetOrderListReq) GetShopId() uint32 {
	if x != nil {
		return x.ShopId
	}
	return 0
}

func (x *GetOrderListReq) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *GetOrderListReq) GetPageSize() uint32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type GetOrderListResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	List []*order_common.OrderResp `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
}

func (x *GetOrderListResp) Reset() {
	*x = GetOrderListResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_bitjump_business_order_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetOrderListResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrderListResp) ProtoMessage() {}

func (x *GetOrderListResp) ProtoReflect() protoreflect.Message {
	mi := &file_order_bitjump_business_order_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrderListResp.ProtoReflect.Descriptor instead.
func (*GetOrderListResp) Descriptor() ([]byte, []int) {
	return file_order_bitjump_business_order_proto_rawDescGZIP(), []int{5}
}

func (x *GetOrderListResp) GetList() []*order_common.OrderResp {
	if x != nil {
		return x.List
	}
	return nil
}

// 提醒消息
type Notify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId uint32 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	Content string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *Notify) Reset() {
	*x = Notify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_bitjump_business_order_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Notify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notify) ProtoMessage() {}

func (x *Notify) ProtoReflect() protoreflect.Message {
	mi := &file_order_bitjump_business_order_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Notify.ProtoReflect.Descriptor instead.
func (*Notify) Descriptor() ([]byte, []int) {
	return file_order_bitjump_business_order_proto_rawDescGZIP(), []int{6}
}

func (x *Notify) GetOrderId() uint32 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

func (x *Notify) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type GetNotifyReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId uint32 `protobuf:"varint,2,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
}

func (x *GetNotifyReq) Reset() {
	*x = GetNotifyReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_bitjump_business_order_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetNotifyReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNotifyReq) ProtoMessage() {}

func (x *GetNotifyReq) ProtoReflect() protoreflect.Message {
	mi := &file_order_bitjump_business_order_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetNotifyReq.ProtoReflect.Descriptor instead.
func (*GetNotifyReq) Descriptor() ([]byte, []int) {
	return file_order_bitjump_business_order_proto_rawDescGZIP(), []int{7}
}

func (x *GetNotifyReq) GetOrderId() uint32 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

type GetNotifyResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Notify *Notify `protobuf:"bytes,1,opt,name=notify,proto3" json:"notify,omitempty"`
}

func (x *GetNotifyResp) Reset() {
	*x = GetNotifyResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_bitjump_business_order_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetNotifyResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNotifyResp) ProtoMessage() {}

func (x *GetNotifyResp) ProtoReflect() protoreflect.Message {
	mi := &file_order_bitjump_business_order_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetNotifyResp.ProtoReflect.Descriptor instead.
func (*GetNotifyResp) Descriptor() ([]byte, []int) {
	return file_order_bitjump_business_order_proto_rawDescGZIP(), []int{8}
}

func (x *GetNotifyResp) GetNotify() *Notify {
	if x != nil {
		return x.Notify
	}
	return nil
}

var File_order_bitjump_business_order_proto protoreflect.FileDescriptor

var file_order_bitjump_business_order_proto_rawDesc = []byte{
	0x0a, 0x22, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2f, 0x62, 0x69, 0x74, 0x6a, 0x75, 0x6d, 0x70, 0x5f,
	0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x62, 0x75, 0x73, 0x69,
	0x6e, 0x65, 0x73, 0x73, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x1a, 0x20, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x2f, 0x62, 0x69, 0x74, 0x6a, 0x75, 0x6d, 0x70, 0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3f, 0x0a, 0x0a, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x52, 0x65, 0x71, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x54, 0x0a, 0x0c,
	0x52, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x19, 0x0a, 0x08,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x10, 0x72, 0x65, 0x6a, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0f, 0x72, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x22, 0x28, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x52, 0x65,
	0x71, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x27, 0x0a, 0x0a,
	0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x5b, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x68, 0x6f, 0x70,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x73, 0x68, 0x6f, 0x70, 0x49,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69,
	0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69,
	0x7a, 0x65, 0x22, 0x3f, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x12, 0x2b, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x52, 0x04, 0x6c,
	0x69, 0x73, 0x74, 0x22, 0x3d, 0x0a, 0x06, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x19, 0x0a,
	0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x22, 0x29, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52,
	0x65, 0x71, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x44, 0x0a,
	0x0d, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70, 0x12, 0x33,
	0x0a, 0x06, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b,
	0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x06, 0x6e, 0x6f, 0x74,
	0x69, 0x66, 0x79, 0x32, 0xc4, 0x04, 0x0a, 0x14, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x75, 0x73,
	0x69, 0x6e, 0x65, 0x73, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5b, 0x0a, 0x0c,
	0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x24, 0x2e, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x2e, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x65, 0x71, 0x1a, 0x25, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x62, 0x75, 0x73, 0x69, 0x6e,
	0x65, 0x73, 0x73, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x12, 0x39, 0x0a, 0x06, 0x44, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x12, 0x16, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x12, 0x3f, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x12,
	0x1f, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x52, 0x65, 0x71,
	0x1a, 0x13, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x41, 0x0a, 0x08, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72,
	0x79, 0x12, 0x20, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65,
	0x73, 0x73, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79,
	0x52, 0x65, 0x71, 0x1a, 0x13, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x3f, 0x0a, 0x07, 0x52, 0x65, 0x63, 0x65,
	0x69, 0x76, 0x65, 0x12, 0x1f, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x62, 0x75, 0x73, 0x69,
	0x6e, 0x65, 0x73, 0x73, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x52, 0x65, 0x71, 0x1a, 0x13, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x43, 0x0a, 0x09, 0x52, 0x65, 0x6a,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x62,
	0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x6a,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x13, 0x2e, 0x6f, 0x72, 0x64, 0x65,
	0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x36,
	0x0a, 0x06, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x12, 0x17, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x52, 0x65,
	0x71, 0x1a, 0x13, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x52, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74,
	0x69, 0x66, 0x79, 0x12, 0x21, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x62, 0x75, 0x73, 0x69,
	0x6e, 0x65, 0x73, 0x73, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74,
	0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x1a, 0x22, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x62,
	0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70, 0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x31, 0x32, 0x33, 0x35, 0x30, 0x38, 0x2f,
	0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x73, 0x68, 0x6f, 0x70, 0x2f, 0x6b, 0x69, 0x74, 0x65, 0x78,
	0x5f, 0x67, 0x65, 0x6e, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2f, 0x62, 0x75, 0x73, 0x69, 0x6e,
	0x65, 0x73, 0x73, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_order_bitjump_business_order_proto_rawDescOnce sync.Once
	file_order_bitjump_business_order_proto_rawDescData = file_order_bitjump_business_order_proto_rawDesc
)

func file_order_bitjump_business_order_proto_rawDescGZIP() []byte {
	file_order_bitjump_business_order_proto_rawDescOnce.Do(func() {
		file_order_bitjump_business_order_proto_rawDescData = protoimpl.X.CompressGZIP(file_order_bitjump_business_order_proto_rawDescData)
	})
	return file_order_bitjump_business_order_proto_rawDescData
}

var file_order_bitjump_business_order_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_order_bitjump_business_order_proto_goTypes = []interface{}{
	(*ConfirmReq)(nil),             // 0: order.businessOrder.ConfirmReq
	(*RejectionReq)(nil),           // 1: order.businessOrder.RejectionReq
	(*DeliveryReq)(nil),            // 2: order.businessOrder.DeliveryReq
	(*ReceiveReq)(nil),             // 3: order.businessOrder.ReceiveReq
	(*GetOrderListReq)(nil),        // 4: order.businessOrder.GetOrderListReq
	(*GetOrderListResp)(nil),       // 5: order.businessOrder.GetOrderListResp
	(*Notify)(nil),                 // 6: order.businessOrder.Notify
	(*GetNotifyReq)(nil),           // 7: order.businessOrder.GetNotifyReq
	(*GetNotifyResp)(nil),          // 8: order.businessOrder.GetNotifyResp
	(*order_common.OrderResp)(nil), // 9: order.common.OrderResp
	(*order_common.OrderReq)(nil),  // 10: order.common.OrderReq
	(*order_common.CancelReq)(nil), // 11: order.common.CancelReq
	(*order_common.Empty)(nil),     // 12: order.common.Empty
}
var file_order_bitjump_business_order_proto_depIdxs = []int32{
	9,  // 0: order.businessOrder.GetOrderListResp.list:type_name -> order.common.OrderResp
	6,  // 1: order.businessOrder.GetNotifyResp.notify:type_name -> order.businessOrder.Notify
	4,  // 2: order.businessOrder.OrderBusinessService.GetOrderList:input_type -> order.businessOrder.GetOrderListReq
	10, // 3: order.businessOrder.OrderBusinessService.Detail:input_type -> order.common.OrderReq
	0,  // 4: order.businessOrder.OrderBusinessService.Confirm:input_type -> order.businessOrder.ConfirmReq
	2,  // 5: order.businessOrder.OrderBusinessService.Delivery:input_type -> order.businessOrder.DeliveryReq
	3,  // 6: order.businessOrder.OrderBusinessService.Receive:input_type -> order.businessOrder.ReceiveReq
	1,  // 7: order.businessOrder.OrderBusinessService.Rejection:input_type -> order.businessOrder.RejectionReq
	11, // 8: order.businessOrder.OrderBusinessService.Cancel:input_type -> order.common.CancelReq
	7,  // 9: order.businessOrder.OrderBusinessService.GetNotify:input_type -> order.businessOrder.GetNotifyReq
	5,  // 10: order.businessOrder.OrderBusinessService.GetOrderList:output_type -> order.businessOrder.GetOrderListResp
	9,  // 11: order.businessOrder.OrderBusinessService.Detail:output_type -> order.common.OrderResp
	12, // 12: order.businessOrder.OrderBusinessService.Confirm:output_type -> order.common.Empty
	12, // 13: order.businessOrder.OrderBusinessService.Delivery:output_type -> order.common.Empty
	12, // 14: order.businessOrder.OrderBusinessService.Receive:output_type -> order.common.Empty
	12, // 15: order.businessOrder.OrderBusinessService.Rejection:output_type -> order.common.Empty
	12, // 16: order.businessOrder.OrderBusinessService.Cancel:output_type -> order.common.Empty
	8,  // 17: order.businessOrder.OrderBusinessService.GetNotify:output_type -> order.businessOrder.GetNotifyResp
	10, // [10:18] is the sub-list for method output_type
	2,  // [2:10] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_order_bitjump_business_order_proto_init() }
func file_order_bitjump_business_order_proto_init() {
	if File_order_bitjump_business_order_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_order_bitjump_business_order_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfirmReq); i {
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
		file_order_bitjump_business_order_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RejectionReq); i {
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
		file_order_bitjump_business_order_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeliveryReq); i {
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
		file_order_bitjump_business_order_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReceiveReq); i {
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
		file_order_bitjump_business_order_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetOrderListReq); i {
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
		file_order_bitjump_business_order_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetOrderListResp); i {
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
		file_order_bitjump_business_order_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Notify); i {
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
		file_order_bitjump_business_order_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetNotifyReq); i {
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
		file_order_bitjump_business_order_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetNotifyResp); i {
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
			RawDescriptor: file_order_bitjump_business_order_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_order_bitjump_business_order_proto_goTypes,
		DependencyIndexes: file_order_bitjump_business_order_proto_depIdxs,
		MessageInfos:      file_order_bitjump_business_order_proto_msgTypes,
	}.Build()
	File_order_bitjump_business_order_proto = out.File
	file_order_bitjump_business_order_proto_rawDesc = nil
	file_order_bitjump_business_order_proto_goTypes = nil
	file_order_bitjump_business_order_proto_depIdxs = nil
}

var _ context.Context

// Code generated by Kitex v0.12.1. DO NOT EDIT.

type OrderBusinessService interface {
	GetOrderList(ctx context.Context, req *GetOrderListReq) (res *GetOrderListResp, err error)
	Detail(ctx context.Context, req *order_common.OrderReq) (res *order_common.OrderResp, err error)
	Confirm(ctx context.Context, req *ConfirmReq) (res *order_common.Empty, err error)
	Delivery(ctx context.Context, req *DeliveryReq) (res *order_common.Empty, err error)
	Receive(ctx context.Context, req *ReceiveReq) (res *order_common.Empty, err error)
	Rejection(ctx context.Context, req *RejectionReq) (res *order_common.Empty, err error)
	Cancel(ctx context.Context, req *order_common.CancelReq) (res *order_common.Empty, err error)
	GetNotify(ctx context.Context, req *GetNotifyReq) (res *GetNotifyResp, err error)
}
