// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.29.2
// source: order/bitjump_order.proto

package userOrder

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

type OrderSubmitReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId        uint32  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	AddressBookId int32   `protobuf:"varint,2,opt,name=address_book_id,json=addressBookId,proto3" json:"address_book_id,omitempty"`
	PayMethod     int32   `protobuf:"varint,3,opt,name=pay_method,json=payMethod,proto3" json:"pay_method,omitempty"`
	Remark        string  `protobuf:"bytes,4,opt,name=remark,proto3" json:"remark,omitempty"`
	Amount        float32 `protobuf:"fixed32,5,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *OrderSubmitReq) Reset() {
	*x = OrderSubmitReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_bitjump_order_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderSubmitReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderSubmitReq) ProtoMessage() {}

func (x *OrderSubmitReq) ProtoReflect() protoreflect.Message {
	mi := &file_order_bitjump_order_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderSubmitReq.ProtoReflect.Descriptor instead.
func (*OrderSubmitReq) Descriptor() ([]byte, []int) {
	return file_order_bitjump_order_proto_rawDescGZIP(), []int{0}
}

func (x *OrderSubmitReq) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *OrderSubmitReq) GetAddressBookId() int32 {
	if x != nil {
		return x.AddressBookId
	}
	return 0
}

func (x *OrderSubmitReq) GetPayMethod() int32 {
	if x != nil {
		return x.PayMethod
	}
	return 0
}

func (x *OrderSubmitReq) GetRemark() string {
	if x != nil {
		return x.Remark
	}
	return ""
}

func (x *OrderSubmitReq) GetAmount() float32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type OrderSubmitResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId     uint32                 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	Number      string                 `protobuf:"bytes,2,opt,name=number,proto3" json:"number,omitempty"`
	OrderAmount float32                `protobuf:"fixed32,3,opt,name=order_amount,json=orderAmount,proto3" json:"order_amount,omitempty"`
	Order       *order_common.OrderReq `protobuf:"bytes,4,opt,name=order,proto3" json:"order,omitempty"`
}

func (x *OrderSubmitResp) Reset() {
	*x = OrderSubmitResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_bitjump_order_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderSubmitResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderSubmitResp) ProtoMessage() {}

func (x *OrderSubmitResp) ProtoReflect() protoreflect.Message {
	mi := &file_order_bitjump_order_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderSubmitResp.ProtoReflect.Descriptor instead.
func (*OrderSubmitResp) Descriptor() ([]byte, []int) {
	return file_order_bitjump_order_proto_rawDescGZIP(), []int{1}
}

func (x *OrderSubmitResp) GetOrderId() uint32 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

func (x *OrderSubmitResp) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *OrderSubmitResp) GetOrderAmount() float32 {
	if x != nil {
		return x.OrderAmount
	}
	return 0
}

func (x *OrderSubmitResp) GetOrder() *order_common.OrderReq {
	if x != nil {
		return x.Order
	}
	return nil
}

type HistoryReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   uint32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Page     uint32 `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	PageSize uint32 `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Status   int32  `protobuf:"varint,4,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *HistoryReq) Reset() {
	*x = HistoryReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_bitjump_order_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HistoryReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HistoryReq) ProtoMessage() {}

func (x *HistoryReq) ProtoReflect() protoreflect.Message {
	mi := &file_order_bitjump_order_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HistoryReq.ProtoReflect.Descriptor instead.
func (*HistoryReq) Descriptor() ([]byte, []int) {
	return file_order_bitjump_order_proto_rawDescGZIP(), []int{2}
}

func (x *HistoryReq) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *HistoryReq) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *HistoryReq) GetPageSize() uint32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *HistoryReq) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

type HistoryResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total    uint32                    `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	List     []*order_common.OrderResp `protobuf:"bytes,2,rep,name=list,proto3" json:"list,omitempty"`
	Page     uint32                    `protobuf:"varint,3,opt,name=page,proto3" json:"page,omitempty"`
	PageSize uint32                    `protobuf:"varint,4,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
}

func (x *HistoryResp) Reset() {
	*x = HistoryResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_bitjump_order_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HistoryResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HistoryResp) ProtoMessage() {}

func (x *HistoryResp) ProtoReflect() protoreflect.Message {
	mi := &file_order_bitjump_order_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HistoryResp.ProtoReflect.Descriptor instead.
func (*HistoryResp) Descriptor() ([]byte, []int) {
	return file_order_bitjump_order_proto_rawDescGZIP(), []int{3}
}

func (x *HistoryResp) GetTotal() uint32 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *HistoryResp) GetList() []*order_common.OrderResp {
	if x != nil {
		return x.List
	}
	return nil
}

func (x *HistoryResp) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *HistoryResp) GetPageSize() uint32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type ReminderReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  uint32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	OrderId uint32 `protobuf:"varint,2,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
}

func (x *ReminderReq) Reset() {
	*x = ReminderReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_bitjump_order_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReminderReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReminderReq) ProtoMessage() {}

func (x *ReminderReq) ProtoReflect() protoreflect.Message {
	mi := &file_order_bitjump_order_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReminderReq.ProtoReflect.Descriptor instead.
func (*ReminderReq) Descriptor() ([]byte, []int) {
	return file_order_bitjump_order_proto_rawDescGZIP(), []int{4}
}

func (x *ReminderReq) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ReminderReq) GetOrderId() uint32 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

type CompleteReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId uint32 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
}

func (x *CompleteReq) Reset() {
	*x = CompleteReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_bitjump_order_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompleteReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompleteReq) ProtoMessage() {}

func (x *CompleteReq) ProtoReflect() protoreflect.Message {
	mi := &file_order_bitjump_order_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompleteReq.ProtoReflect.Descriptor instead.
func (*CompleteReq) Descriptor() ([]byte, []int) {
	return file_order_bitjump_order_proto_rawDescGZIP(), []int{5}
}

func (x *CompleteReq) GetOrderId() uint32 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

var File_order_bitjump_order_proto protoreflect.FileDescriptor

var file_order_bitjump_order_proto_rawDesc = []byte{
	0x0a, 0x19, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2f, 0x62, 0x69, 0x74, 0x6a, 0x75, 0x6d, 0x70, 0x5f,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x1a, 0x20, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x2f, 0x62, 0x69, 0x74, 0x6a, 0x75, 0x6d, 0x70, 0x5f, 0x6f, 0x72, 0x64, 0x65,
	0x72, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa0,
	0x01, 0x0a, 0x0e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x52, 0x65,
	0x71, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x5f, 0x62, 0x6f, 0x6f, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0d, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x42, 0x6f, 0x6f, 0x6b,
	0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x79, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x70, 0x61, 0x79, 0x4d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x72, 0x65, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x02, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x22, 0x95, 0x01, 0x0a, 0x0f, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x53, 0x75, 0x62, 0x6d, 0x69,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x6f, 0x72, 0x64, 0x65,
	0x72, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0b,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2c, 0x0a, 0x05, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x52, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x22, 0x6e, 0x0a, 0x0a, 0x48, 0x69, 0x73,
	0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x81, 0x01, 0x0a, 0x0b, 0x48, 0x69,
	0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12,
	0x2b, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x41, 0x0a,
	0x0b, 0x52, 0x65, 0x6d, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64,
	0x22, 0x28, 0x0a, 0x0b, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x12,
	0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x32, 0x96, 0x03, 0x0a, 0x10, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x4b, 0x0a, 0x06, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x12, 0x1f, 0x2e, 0x6f, 0x72, 0x64, 0x65,
	0x72, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x20, 0x2e, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x12, 0x44, 0x0a, 0x07,
	0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x1b, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e,
	0x75, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72,
	0x79, 0x52, 0x65, 0x71, 0x1a, 0x1c, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65,
	0x73, 0x70, 0x12, 0x39, 0x0a, 0x06, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x16, 0x2e, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x12, 0x36, 0x0a,
	0x06, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x12, 0x17, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x52, 0x65, 0x71,
	0x1a, 0x13, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x3d, 0x0a, 0x08, 0x52, 0x65, 0x6d, 0x69, 0x6e, 0x64, 0x65,
	0x72, 0x12, 0x1c, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x6d, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x1a,
	0x13, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x12, 0x3d, 0x0a, 0x08, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65,
	0x12, 0x1c, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x13,
	0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x31, 0x32, 0x33, 0x35, 0x30, 0x38, 0x2f, 0x64, 0x6f, 0x75, 0x79, 0x69, 0x6e, 0x73,
	0x68, 0x6f, 0x70, 0x2f, 0x6b, 0x69, 0x74, 0x65, 0x78, 0x5f, 0x67, 0x65, 0x6e, 0x2f, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_order_bitjump_order_proto_rawDescOnce sync.Once
	file_order_bitjump_order_proto_rawDescData = file_order_bitjump_order_proto_rawDesc
)

func file_order_bitjump_order_proto_rawDescGZIP() []byte {
	file_order_bitjump_order_proto_rawDescOnce.Do(func() {
		file_order_bitjump_order_proto_rawDescData = protoimpl.X.CompressGZIP(file_order_bitjump_order_proto_rawDescData)
	})
	return file_order_bitjump_order_proto_rawDescData
}

var file_order_bitjump_order_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_order_bitjump_order_proto_goTypes = []interface{}{
	(*OrderSubmitReq)(nil),         // 0: order.userOrder.OrderSubmitReq
	(*OrderSubmitResp)(nil),        // 1: order.userOrder.OrderSubmitResp
	(*HistoryReq)(nil),             // 2: order.userOrder.HistoryReq
	(*HistoryResp)(nil),            // 3: order.userOrder.HistoryResp
	(*ReminderReq)(nil),            // 4: order.userOrder.ReminderReq
	(*CompleteReq)(nil),            // 5: order.userOrder.CompleteReq
	(*order_common.OrderReq)(nil),  // 6: order.common.OrderReq
	(*order_common.OrderResp)(nil), // 7: order.common.OrderResp
	(*order_common.CancelReq)(nil), // 8: order.common.CancelReq
	(*order_common.Empty)(nil),     // 9: order.common.Empty
}
var file_order_bitjump_order_proto_depIdxs = []int32{
	6, // 0: order.userOrder.OrderSubmitResp.order:type_name -> order.common.OrderReq
	7, // 1: order.userOrder.HistoryResp.list:type_name -> order.common.OrderResp
	0, // 2: order.userOrder.OrderUserService.Submit:input_type -> order.userOrder.OrderSubmitReq
	2, // 3: order.userOrder.OrderUserService.History:input_type -> order.userOrder.HistoryReq
	6, // 4: order.userOrder.OrderUserService.Detail:input_type -> order.common.OrderReq
	8, // 5: order.userOrder.OrderUserService.Cancel:input_type -> order.common.CancelReq
	4, // 6: order.userOrder.OrderUserService.Reminder:input_type -> order.userOrder.ReminderReq
	5, // 7: order.userOrder.OrderUserService.Complete:input_type -> order.userOrder.CompleteReq
	1, // 8: order.userOrder.OrderUserService.Submit:output_type -> order.userOrder.OrderSubmitResp
	3, // 9: order.userOrder.OrderUserService.History:output_type -> order.userOrder.HistoryResp
	7, // 10: order.userOrder.OrderUserService.Detail:output_type -> order.common.OrderResp
	9, // 11: order.userOrder.OrderUserService.Cancel:output_type -> order.common.Empty
	9, // 12: order.userOrder.OrderUserService.Reminder:output_type -> order.common.Empty
	9, // 13: order.userOrder.OrderUserService.Complete:output_type -> order.common.Empty
	8, // [8:14] is the sub-list for method output_type
	2, // [2:8] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_order_bitjump_order_proto_init() }
func file_order_bitjump_order_proto_init() {
	if File_order_bitjump_order_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_order_bitjump_order_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OrderSubmitReq); i {
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
		file_order_bitjump_order_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OrderSubmitResp); i {
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
		file_order_bitjump_order_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HistoryReq); i {
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
		file_order_bitjump_order_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HistoryResp); i {
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
		file_order_bitjump_order_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReminderReq); i {
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
		file_order_bitjump_order_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CompleteReq); i {
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
			RawDescriptor: file_order_bitjump_order_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_order_bitjump_order_proto_goTypes,
		DependencyIndexes: file_order_bitjump_order_proto_depIdxs,
		MessageInfos:      file_order_bitjump_order_proto_msgTypes,
	}.Build()
	File_order_bitjump_order_proto = out.File
	file_order_bitjump_order_proto_rawDesc = nil
	file_order_bitjump_order_proto_goTypes = nil
	file_order_bitjump_order_proto_depIdxs = nil
}

var _ context.Context

// Code generated by Kitex v0.12.1. DO NOT EDIT.

type OrderUserService interface {
	Submit(ctx context.Context, req *OrderSubmitReq) (res *OrderSubmitResp, err error)
	History(ctx context.Context, req *HistoryReq) (res *HistoryResp, err error)
	Detail(ctx context.Context, req *order_common.OrderReq) (res *order_common.OrderResp, err error)
	Cancel(ctx context.Context, req *order_common.CancelReq) (res *order_common.Empty, err error)
	Reminder(ctx context.Context, req *ReminderReq) (res *order_common.Empty, err error)
	Complete(ctx context.Context, req *CompleteReq) (res *order_common.Empty, err error)
}
