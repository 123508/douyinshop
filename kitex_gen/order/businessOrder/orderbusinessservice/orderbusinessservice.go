// Code generated by Kitex v0.12.1. DO NOT EDIT.

package orderbusinessservice

import (
	"context"
	"errors"
	businessOrder "github.com/123508/douyinshop/kitex_gen/order/businessOrder"
	order_common "github.com/123508/douyinshop/kitex_gen/order/order_common"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
	proto "google.golang.org/protobuf/proto"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"GetOrderList": kitex.NewMethodInfo(
		getOrderListHandler,
		newGetOrderListArgs,
		newGetOrderListResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"Detail": kitex.NewMethodInfo(
		detailHandler,
		newDetailArgs,
		newDetailResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"Confirm": kitex.NewMethodInfo(
		confirmHandler,
		newConfirmArgs,
		newConfirmResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"Delivery": kitex.NewMethodInfo(
		deliveryHandler,
		newDeliveryArgs,
		newDeliveryResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"Receive": kitex.NewMethodInfo(
		receiveHandler,
		newReceiveArgs,
		newReceiveResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"Rejection": kitex.NewMethodInfo(
		rejectionHandler,
		newRejectionArgs,
		newRejectionResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"Cancel": kitex.NewMethodInfo(
		cancelHandler,
		newCancelArgs,
		newCancelResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
}

var (
	orderBusinessServiceServiceInfo                = NewServiceInfo()
	orderBusinessServiceServiceInfoForClient       = NewServiceInfoForClient()
	orderBusinessServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return orderBusinessServiceServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return orderBusinessServiceServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return orderBusinessServiceServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "OrderBusinessService"
	handlerType := (*businessOrder.OrderBusinessService)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "order.businessOrder",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Protobuf,
		KiteXGenVersion: "v0.12.1",
		Extra:           extra,
	}
	return svcInfo
}

func getOrderListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(businessOrder.GetOrderListReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(businessOrder.OrderBusinessService).GetOrderList(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *GetOrderListArgs:
		success, err := handler.(businessOrder.OrderBusinessService).GetOrderList(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*GetOrderListResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newGetOrderListArgs() interface{} {
	return &GetOrderListArgs{}
}

func newGetOrderListResult() interface{} {
	return &GetOrderListResult{}
}

type GetOrderListArgs struct {
	Req *businessOrder.GetOrderListReq
}

func (p *GetOrderListArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(businessOrder.GetOrderListReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *GetOrderListArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *GetOrderListArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *GetOrderListArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *GetOrderListArgs) Unmarshal(in []byte) error {
	msg := new(businessOrder.GetOrderListReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var GetOrderListArgs_Req_DEFAULT *businessOrder.GetOrderListReq

func (p *GetOrderListArgs) GetReq() *businessOrder.GetOrderListReq {
	if !p.IsSetReq() {
		return GetOrderListArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *GetOrderListArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *GetOrderListArgs) GetFirstArgument() interface{} {
	return p.Req
}

type GetOrderListResult struct {
	Success *businessOrder.GetOrderListResp
}

var GetOrderListResult_Success_DEFAULT *businessOrder.GetOrderListResp

func (p *GetOrderListResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(businessOrder.GetOrderListResp)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *GetOrderListResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *GetOrderListResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *GetOrderListResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *GetOrderListResult) Unmarshal(in []byte) error {
	msg := new(businessOrder.GetOrderListResp)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *GetOrderListResult) GetSuccess() *businessOrder.GetOrderListResp {
	if !p.IsSetSuccess() {
		return GetOrderListResult_Success_DEFAULT
	}
	return p.Success
}

func (p *GetOrderListResult) SetSuccess(x interface{}) {
	p.Success = x.(*businessOrder.GetOrderListResp)
}

func (p *GetOrderListResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *GetOrderListResult) GetResult() interface{} {
	return p.Success
}

func detailHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(order_common.OrderReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(businessOrder.OrderBusinessService).Detail(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *DetailArgs:
		success, err := handler.(businessOrder.OrderBusinessService).Detail(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*DetailResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newDetailArgs() interface{} {
	return &DetailArgs{}
}

func newDetailResult() interface{} {
	return &DetailResult{}
}

type DetailArgs struct {
	Req *order_common.OrderReq
}

func (p *DetailArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(order_common.OrderReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *DetailArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *DetailArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *DetailArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *DetailArgs) Unmarshal(in []byte) error {
	msg := new(order_common.OrderReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var DetailArgs_Req_DEFAULT *order_common.OrderReq

func (p *DetailArgs) GetReq() *order_common.OrderReq {
	if !p.IsSetReq() {
		return DetailArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *DetailArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *DetailArgs) GetFirstArgument() interface{} {
	return p.Req
}

type DetailResult struct {
	Success *order_common.OrderResp
}

var DetailResult_Success_DEFAULT *order_common.OrderResp

func (p *DetailResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(order_common.OrderResp)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *DetailResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *DetailResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *DetailResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *DetailResult) Unmarshal(in []byte) error {
	msg := new(order_common.OrderResp)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *DetailResult) GetSuccess() *order_common.OrderResp {
	if !p.IsSetSuccess() {
		return DetailResult_Success_DEFAULT
	}
	return p.Success
}

func (p *DetailResult) SetSuccess(x interface{}) {
	p.Success = x.(*order_common.OrderResp)
}

func (p *DetailResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *DetailResult) GetResult() interface{} {
	return p.Success
}

func confirmHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(businessOrder.ConfirmReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(businessOrder.OrderBusinessService).Confirm(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *ConfirmArgs:
		success, err := handler.(businessOrder.OrderBusinessService).Confirm(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*ConfirmResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newConfirmArgs() interface{} {
	return &ConfirmArgs{}
}

func newConfirmResult() interface{} {
	return &ConfirmResult{}
}

type ConfirmArgs struct {
	Req *businessOrder.ConfirmReq
}

func (p *ConfirmArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(businessOrder.ConfirmReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *ConfirmArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *ConfirmArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *ConfirmArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *ConfirmArgs) Unmarshal(in []byte) error {
	msg := new(businessOrder.ConfirmReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var ConfirmArgs_Req_DEFAULT *businessOrder.ConfirmReq

func (p *ConfirmArgs) GetReq() *businessOrder.ConfirmReq {
	if !p.IsSetReq() {
		return ConfirmArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *ConfirmArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *ConfirmArgs) GetFirstArgument() interface{} {
	return p.Req
}

type ConfirmResult struct {
	Success *order_common.Empty
}

var ConfirmResult_Success_DEFAULT *order_common.Empty

func (p *ConfirmResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(order_common.Empty)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *ConfirmResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *ConfirmResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *ConfirmResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *ConfirmResult) Unmarshal(in []byte) error {
	msg := new(order_common.Empty)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *ConfirmResult) GetSuccess() *order_common.Empty {
	if !p.IsSetSuccess() {
		return ConfirmResult_Success_DEFAULT
	}
	return p.Success
}

func (p *ConfirmResult) SetSuccess(x interface{}) {
	p.Success = x.(*order_common.Empty)
}

func (p *ConfirmResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *ConfirmResult) GetResult() interface{} {
	return p.Success
}

func deliveryHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(businessOrder.DeliveryReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(businessOrder.OrderBusinessService).Delivery(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *DeliveryArgs:
		success, err := handler.(businessOrder.OrderBusinessService).Delivery(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*DeliveryResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newDeliveryArgs() interface{} {
	return &DeliveryArgs{}
}

func newDeliveryResult() interface{} {
	return &DeliveryResult{}
}

type DeliveryArgs struct {
	Req *businessOrder.DeliveryReq
}

func (p *DeliveryArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(businessOrder.DeliveryReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *DeliveryArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *DeliveryArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *DeliveryArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *DeliveryArgs) Unmarshal(in []byte) error {
	msg := new(businessOrder.DeliveryReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var DeliveryArgs_Req_DEFAULT *businessOrder.DeliveryReq

func (p *DeliveryArgs) GetReq() *businessOrder.DeliveryReq {
	if !p.IsSetReq() {
		return DeliveryArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *DeliveryArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *DeliveryArgs) GetFirstArgument() interface{} {
	return p.Req
}

type DeliveryResult struct {
	Success *order_common.Empty
}

var DeliveryResult_Success_DEFAULT *order_common.Empty

func (p *DeliveryResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(order_common.Empty)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *DeliveryResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *DeliveryResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *DeliveryResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *DeliveryResult) Unmarshal(in []byte) error {
	msg := new(order_common.Empty)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *DeliveryResult) GetSuccess() *order_common.Empty {
	if !p.IsSetSuccess() {
		return DeliveryResult_Success_DEFAULT
	}
	return p.Success
}

func (p *DeliveryResult) SetSuccess(x interface{}) {
	p.Success = x.(*order_common.Empty)
}

func (p *DeliveryResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *DeliveryResult) GetResult() interface{} {
	return p.Success
}

func receiveHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(businessOrder.ReceiveReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(businessOrder.OrderBusinessService).Receive(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *ReceiveArgs:
		success, err := handler.(businessOrder.OrderBusinessService).Receive(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*ReceiveResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newReceiveArgs() interface{} {
	return &ReceiveArgs{}
}

func newReceiveResult() interface{} {
	return &ReceiveResult{}
}

type ReceiveArgs struct {
	Req *businessOrder.ReceiveReq
}

func (p *ReceiveArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(businessOrder.ReceiveReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *ReceiveArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *ReceiveArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *ReceiveArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *ReceiveArgs) Unmarshal(in []byte) error {
	msg := new(businessOrder.ReceiveReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var ReceiveArgs_Req_DEFAULT *businessOrder.ReceiveReq

func (p *ReceiveArgs) GetReq() *businessOrder.ReceiveReq {
	if !p.IsSetReq() {
		return ReceiveArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *ReceiveArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *ReceiveArgs) GetFirstArgument() interface{} {
	return p.Req
}

type ReceiveResult struct {
	Success *order_common.Empty
}

var ReceiveResult_Success_DEFAULT *order_common.Empty

func (p *ReceiveResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(order_common.Empty)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *ReceiveResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *ReceiveResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *ReceiveResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *ReceiveResult) Unmarshal(in []byte) error {
	msg := new(order_common.Empty)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *ReceiveResult) GetSuccess() *order_common.Empty {
	if !p.IsSetSuccess() {
		return ReceiveResult_Success_DEFAULT
	}
	return p.Success
}

func (p *ReceiveResult) SetSuccess(x interface{}) {
	p.Success = x.(*order_common.Empty)
}

func (p *ReceiveResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *ReceiveResult) GetResult() interface{} {
	return p.Success
}

func rejectionHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(businessOrder.ReceiveReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(businessOrder.OrderBusinessService).Rejection(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *RejectionArgs:
		success, err := handler.(businessOrder.OrderBusinessService).Rejection(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*RejectionResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newRejectionArgs() interface{} {
	return &RejectionArgs{}
}

func newRejectionResult() interface{} {
	return &RejectionResult{}
}

type RejectionArgs struct {
	Req *businessOrder.ReceiveReq
}

func (p *RejectionArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(businessOrder.ReceiveReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *RejectionArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *RejectionArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *RejectionArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *RejectionArgs) Unmarshal(in []byte) error {
	msg := new(businessOrder.ReceiveReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var RejectionArgs_Req_DEFAULT *businessOrder.ReceiveReq

func (p *RejectionArgs) GetReq() *businessOrder.ReceiveReq {
	if !p.IsSetReq() {
		return RejectionArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *RejectionArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *RejectionArgs) GetFirstArgument() interface{} {
	return p.Req
}

type RejectionResult struct {
	Success *order_common.Empty
}

var RejectionResult_Success_DEFAULT *order_common.Empty

func (p *RejectionResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(order_common.Empty)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *RejectionResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *RejectionResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *RejectionResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *RejectionResult) Unmarshal(in []byte) error {
	msg := new(order_common.Empty)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *RejectionResult) GetSuccess() *order_common.Empty {
	if !p.IsSetSuccess() {
		return RejectionResult_Success_DEFAULT
	}
	return p.Success
}

func (p *RejectionResult) SetSuccess(x interface{}) {
	p.Success = x.(*order_common.Empty)
}

func (p *RejectionResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *RejectionResult) GetResult() interface{} {
	return p.Success
}

func cancelHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(order_common.CancelReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(businessOrder.OrderBusinessService).Cancel(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *CancelArgs:
		success, err := handler.(businessOrder.OrderBusinessService).Cancel(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*CancelResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newCancelArgs() interface{} {
	return &CancelArgs{}
}

func newCancelResult() interface{} {
	return &CancelResult{}
}

type CancelArgs struct {
	Req *order_common.CancelReq
}

func (p *CancelArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(order_common.CancelReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *CancelArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *CancelArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *CancelArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *CancelArgs) Unmarshal(in []byte) error {
	msg := new(order_common.CancelReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var CancelArgs_Req_DEFAULT *order_common.CancelReq

func (p *CancelArgs) GetReq() *order_common.CancelReq {
	if !p.IsSetReq() {
		return CancelArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *CancelArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *CancelArgs) GetFirstArgument() interface{} {
	return p.Req
}

type CancelResult struct {
	Success *order_common.Empty
}

var CancelResult_Success_DEFAULT *order_common.Empty

func (p *CancelResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(order_common.Empty)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *CancelResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *CancelResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *CancelResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *CancelResult) Unmarshal(in []byte) error {
	msg := new(order_common.Empty)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *CancelResult) GetSuccess() *order_common.Empty {
	if !p.IsSetSuccess() {
		return CancelResult_Success_DEFAULT
	}
	return p.Success
}

func (p *CancelResult) SetSuccess(x interface{}) {
	p.Success = x.(*order_common.Empty)
}

func (p *CancelResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *CancelResult) GetResult() interface{} {
	return p.Success
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) GetOrderList(ctx context.Context, Req *businessOrder.GetOrderListReq) (r *businessOrder.GetOrderListResp, err error) {
	var _args GetOrderListArgs
	_args.Req = Req
	var _result GetOrderListResult
	if err = p.c.Call(ctx, "GetOrderList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) Detail(ctx context.Context, Req *order_common.OrderReq) (r *order_common.OrderResp, err error) {
	var _args DetailArgs
	_args.Req = Req
	var _result DetailResult
	if err = p.c.Call(ctx, "Detail", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) Confirm(ctx context.Context, Req *businessOrder.ConfirmReq) (r *order_common.Empty, err error) {
	var _args ConfirmArgs
	_args.Req = Req
	var _result ConfirmResult
	if err = p.c.Call(ctx, "Confirm", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) Delivery(ctx context.Context, Req *businessOrder.DeliveryReq) (r *order_common.Empty, err error) {
	var _args DeliveryArgs
	_args.Req = Req
	var _result DeliveryResult
	if err = p.c.Call(ctx, "Delivery", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) Receive(ctx context.Context, Req *businessOrder.ReceiveReq) (r *order_common.Empty, err error) {
	var _args ReceiveArgs
	_args.Req = Req
	var _result ReceiveResult
	if err = p.c.Call(ctx, "Receive", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) Rejection(ctx context.Context, Req *businessOrder.ReceiveReq) (r *order_common.Empty, err error) {
	var _args RejectionArgs
	_args.Req = Req
	var _result RejectionResult
	if err = p.c.Call(ctx, "Rejection", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) Cancel(ctx context.Context, Req *order_common.CancelReq) (r *order_common.Empty, err error) {
	var _args CancelArgs
	_args.Req = Req
	var _result CancelResult
	if err = p.c.Call(ctx, "Cancel", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
