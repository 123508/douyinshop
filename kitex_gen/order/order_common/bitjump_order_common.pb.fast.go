// Code generated by Fastpb v0.0.2. DO NOT EDIT.

package order_common

import (
	fmt "fmt"
	fastpb "github.com/cloudwego/fastpb"
)

var (
	_ = fmt.Errorf
	_ = fastpb.Skip
)

func (x *Order) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 4:
		offset, err = x.fastReadField4(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 5:
		offset, err = x.fastReadField5(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 6:
		offset, err = x.fastReadField6(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 7:
		offset, err = x.fastReadField7(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 8:
		offset, err = x.fastReadField8(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 9:
		offset, err = x.fastReadField9(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 10:
		offset, err = x.fastReadField10(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 11:
		offset, err = x.fastReadField11(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 12:
		offset, err = x.fastReadField12(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 13:
		offset, err = x.fastReadField13(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 14:
		offset, err = x.fastReadField14(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 15:
		offset, err = x.fastReadField15(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 16:
		offset, err = x.fastReadField16(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_Order[number], err)
}

func (x *Order) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *Order) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.Number, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Order) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	x.Status, offset, err = fastpb.ReadInt32(buf, _type)
	return offset, err
}

func (x *Order) fastReadField4(buf []byte, _type int8) (offset int, err error) {
	x.AddressBookId, offset, err = fastpb.ReadUint64(buf, _type)
	return offset, err
}

func (x *Order) fastReadField5(buf []byte, _type int8) (offset int, err error) {
	x.PayMethod, offset, err = fastpb.ReadInt32(buf, _type)
	return offset, err
}

func (x *Order) fastReadField6(buf []byte, _type int8) (offset int, err error) {
	x.PayStatus, offset, err = fastpb.ReadInt32(buf, _type)
	return offset, err
}

func (x *Order) fastReadField7(buf []byte, _type int8) (offset int, err error) {
	x.Amount, offset, err = fastpb.ReadFloat(buf, _type)
	return offset, err
}

func (x *Order) fastReadField8(buf []byte, _type int8) (offset int, err error) {
	x.Remark, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Order) fastReadField9(buf []byte, _type int8) (offset int, err error) {
	x.Phone, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Order) fastReadField10(buf []byte, _type int8) (offset int, err error) {
	x.Address, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Order) fastReadField11(buf []byte, _type int8) (offset int, err error) {
	x.Username, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Order) fastReadField12(buf []byte, _type int8) (offset int, err error) {
	x.Consignee, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Order) fastReadField13(buf []byte, _type int8) (offset int, err error) {
	x.CancelReason, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Order) fastReadField14(buf []byte, _type int8) (offset int, err error) {
	x.RejectionReason, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Order) fastReadField15(buf []byte, _type int8) (offset int, err error) {
	x.ID, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *Order) fastReadField16(buf []byte, _type int8) (offset int, err error) {
	x.ShopId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *OrderDetail) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 4:
		offset, err = x.fastReadField4(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 5:
		offset, err = x.fastReadField5(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 6:
		offset, err = x.fastReadField6(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_OrderDetail[number], err)
}

func (x *OrderDetail) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.Name, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *OrderDetail) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.Image, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *OrderDetail) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	x.OrderId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *OrderDetail) fastReadField4(buf []byte, _type int8) (offset int, err error) {
	x.ProductId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *OrderDetail) fastReadField5(buf []byte, _type int8) (offset int, err error) {
	x.Number, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *OrderDetail) fastReadField6(buf []byte, _type int8) (offset int, err error) {
	x.Amount, offset, err = fastpb.ReadFloat(buf, _type)
	return offset, err
}

func (x *OrderResp) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_OrderResp[number], err)
}

func (x *OrderResp) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	var v Order
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.Order = &v
	return offset, nil
}

func (x *OrderResp) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	var v OrderDetail
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.OrderDetails = append(x.OrderDetails, &v)
	return offset, nil
}

func (x *Empty) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
}

func (x *OrderReq) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_OrderReq[number], err)
}

func (x *OrderReq) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.OrderId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *OrderReq) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	var v OrderDetail
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.List = append(x.List, &v)
	return offset, nil
}

func (x *CancelReq) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_CancelReq[number], err)
}

func (x *CancelReq) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.OrderId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *CancelReq) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.CancelReason, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Order) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	offset += x.fastWriteField4(buf[offset:])
	offset += x.fastWriteField5(buf[offset:])
	offset += x.fastWriteField6(buf[offset:])
	offset += x.fastWriteField7(buf[offset:])
	offset += x.fastWriteField8(buf[offset:])
	offset += x.fastWriteField9(buf[offset:])
	offset += x.fastWriteField10(buf[offset:])
	offset += x.fastWriteField11(buf[offset:])
	offset += x.fastWriteField12(buf[offset:])
	offset += x.fastWriteField13(buf[offset:])
	offset += x.fastWriteField14(buf[offset:])
	offset += x.fastWriteField15(buf[offset:])
	offset += x.fastWriteField16(buf[offset:])
	return offset
}

func (x *Order) fastWriteField1(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.GetUserId())
	return offset
}

func (x *Order) fastWriteField2(buf []byte) (offset int) {
	if x.Number == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetNumber())
	return offset
}

func (x *Order) fastWriteField3(buf []byte) (offset int) {
	if x.Status == 0 {
		return offset
	}
	offset += fastpb.WriteInt32(buf[offset:], 3, x.GetStatus())
	return offset
}

func (x *Order) fastWriteField4(buf []byte) (offset int) {
	if x.AddressBookId == 0 {
		return offset
	}
	offset += fastpb.WriteUint64(buf[offset:], 4, x.GetAddressBookId())
	return offset
}

func (x *Order) fastWriteField5(buf []byte) (offset int) {
	if x.PayMethod == 0 {
		return offset
	}
	offset += fastpb.WriteInt32(buf[offset:], 5, x.GetPayMethod())
	return offset
}

func (x *Order) fastWriteField6(buf []byte) (offset int) {
	if x.PayStatus == 0 {
		return offset
	}
	offset += fastpb.WriteInt32(buf[offset:], 6, x.GetPayStatus())
	return offset
}

func (x *Order) fastWriteField7(buf []byte) (offset int) {
	if x.Amount == 0 {
		return offset
	}
	offset += fastpb.WriteFloat(buf[offset:], 7, x.GetAmount())
	return offset
}

func (x *Order) fastWriteField8(buf []byte) (offset int) {
	if x.Remark == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 8, x.GetRemark())
	return offset
}

func (x *Order) fastWriteField9(buf []byte) (offset int) {
	if x.Phone == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 9, x.GetPhone())
	return offset
}

func (x *Order) fastWriteField10(buf []byte) (offset int) {
	if x.Address == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 10, x.GetAddress())
	return offset
}

func (x *Order) fastWriteField11(buf []byte) (offset int) {
	if x.Username == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 11, x.GetUsername())
	return offset
}

func (x *Order) fastWriteField12(buf []byte) (offset int) {
	if x.Consignee == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 12, x.GetConsignee())
	return offset
}

func (x *Order) fastWriteField13(buf []byte) (offset int) {
	if x.CancelReason == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 13, x.GetCancelReason())
	return offset
}

func (x *Order) fastWriteField14(buf []byte) (offset int) {
	if x.RejectionReason == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 14, x.GetRejectionReason())
	return offset
}

func (x *Order) fastWriteField15(buf []byte) (offset int) {
	if x.ID == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 15, x.GetID())
	return offset
}

func (x *Order) fastWriteField16(buf []byte) (offset int) {
	if x.ShopId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 16, x.GetShopId())
	return offset
}

func (x *OrderDetail) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	offset += x.fastWriteField4(buf[offset:])
	offset += x.fastWriteField5(buf[offset:])
	offset += x.fastWriteField6(buf[offset:])
	return offset
}

func (x *OrderDetail) fastWriteField1(buf []byte) (offset int) {
	if x.Name == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 1, x.GetName())
	return offset
}

func (x *OrderDetail) fastWriteField2(buf []byte) (offset int) {
	if x.Image == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetImage())
	return offset
}

func (x *OrderDetail) fastWriteField3(buf []byte) (offset int) {
	if x.OrderId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 3, x.GetOrderId())
	return offset
}

func (x *OrderDetail) fastWriteField4(buf []byte) (offset int) {
	if x.ProductId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 4, x.GetProductId())
	return offset
}

func (x *OrderDetail) fastWriteField5(buf []byte) (offset int) {
	if x.Number == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 5, x.GetNumber())
	return offset
}

func (x *OrderDetail) fastWriteField6(buf []byte) (offset int) {
	if x.Amount == 0 {
		return offset
	}
	offset += fastpb.WriteFloat(buf[offset:], 6, x.GetAmount())
	return offset
}

func (x *OrderResp) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *OrderResp) fastWriteField1(buf []byte) (offset int) {
	if x.Order == nil {
		return offset
	}
	offset += fastpb.WriteMessage(buf[offset:], 1, x.GetOrder())
	return offset
}

func (x *OrderResp) fastWriteField2(buf []byte) (offset int) {
	if x.OrderDetails == nil {
		return offset
	}
	for i := range x.GetOrderDetails() {
		offset += fastpb.WriteMessage(buf[offset:], 2, x.GetOrderDetails()[i])
	}
	return offset
}

func (x *Empty) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	return offset
}

func (x *OrderReq) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *OrderReq) fastWriteField1(buf []byte) (offset int) {
	if x.OrderId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.GetOrderId())
	return offset
}

func (x *OrderReq) fastWriteField2(buf []byte) (offset int) {
	if x.List == nil {
		return offset
	}
	for i := range x.GetList() {
		offset += fastpb.WriteMessage(buf[offset:], 2, x.GetList()[i])
	}
	return offset
}

func (x *CancelReq) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *CancelReq) fastWriteField1(buf []byte) (offset int) {
	if x.OrderId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.GetOrderId())
	return offset
}

func (x *CancelReq) fastWriteField2(buf []byte) (offset int) {
	if x.CancelReason == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetCancelReason())
	return offset
}

func (x *Order) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	n += x.sizeField4()
	n += x.sizeField5()
	n += x.sizeField6()
	n += x.sizeField7()
	n += x.sizeField8()
	n += x.sizeField9()
	n += x.sizeField10()
	n += x.sizeField11()
	n += x.sizeField12()
	n += x.sizeField13()
	n += x.sizeField14()
	n += x.sizeField15()
	n += x.sizeField16()
	return n
}

func (x *Order) sizeField1() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.GetUserId())
	return n
}

func (x *Order) sizeField2() (n int) {
	if x.Number == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetNumber())
	return n
}

func (x *Order) sizeField3() (n int) {
	if x.Status == 0 {
		return n
	}
	n += fastpb.SizeInt32(3, x.GetStatus())
	return n
}

func (x *Order) sizeField4() (n int) {
	if x.AddressBookId == 0 {
		return n
	}
	n += fastpb.SizeUint64(4, x.GetAddressBookId())
	return n
}

func (x *Order) sizeField5() (n int) {
	if x.PayMethod == 0 {
		return n
	}
	n += fastpb.SizeInt32(5, x.GetPayMethod())
	return n
}

func (x *Order) sizeField6() (n int) {
	if x.PayStatus == 0 {
		return n
	}
	n += fastpb.SizeInt32(6, x.GetPayStatus())
	return n
}

func (x *Order) sizeField7() (n int) {
	if x.Amount == 0 {
		return n
	}
	n += fastpb.SizeFloat(7, x.GetAmount())
	return n
}

func (x *Order) sizeField8() (n int) {
	if x.Remark == "" {
		return n
	}
	n += fastpb.SizeString(8, x.GetRemark())
	return n
}

func (x *Order) sizeField9() (n int) {
	if x.Phone == "" {
		return n
	}
	n += fastpb.SizeString(9, x.GetPhone())
	return n
}

func (x *Order) sizeField10() (n int) {
	if x.Address == "" {
		return n
	}
	n += fastpb.SizeString(10, x.GetAddress())
	return n
}

func (x *Order) sizeField11() (n int) {
	if x.Username == "" {
		return n
	}
	n += fastpb.SizeString(11, x.GetUsername())
	return n
}

func (x *Order) sizeField12() (n int) {
	if x.Consignee == "" {
		return n
	}
	n += fastpb.SizeString(12, x.GetConsignee())
	return n
}

func (x *Order) sizeField13() (n int) {
	if x.CancelReason == "" {
		return n
	}
	n += fastpb.SizeString(13, x.GetCancelReason())
	return n
}

func (x *Order) sizeField14() (n int) {
	if x.RejectionReason == "" {
		return n
	}
	n += fastpb.SizeString(14, x.GetRejectionReason())
	return n
}

func (x *Order) sizeField15() (n int) {
	if x.ID == 0 {
		return n
	}
	n += fastpb.SizeUint32(15, x.GetID())
	return n
}

func (x *Order) sizeField16() (n int) {
	if x.ShopId == 0 {
		return n
	}
	n += fastpb.SizeUint32(16, x.GetShopId())
	return n
}

func (x *OrderDetail) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	n += x.sizeField4()
	n += x.sizeField5()
	n += x.sizeField6()
	return n
}

func (x *OrderDetail) sizeField1() (n int) {
	if x.Name == "" {
		return n
	}
	n += fastpb.SizeString(1, x.GetName())
	return n
}

func (x *OrderDetail) sizeField2() (n int) {
	if x.Image == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetImage())
	return n
}

func (x *OrderDetail) sizeField3() (n int) {
	if x.OrderId == 0 {
		return n
	}
	n += fastpb.SizeUint32(3, x.GetOrderId())
	return n
}

func (x *OrderDetail) sizeField4() (n int) {
	if x.ProductId == 0 {
		return n
	}
	n += fastpb.SizeUint32(4, x.GetProductId())
	return n
}

func (x *OrderDetail) sizeField5() (n int) {
	if x.Number == 0 {
		return n
	}
	n += fastpb.SizeUint32(5, x.GetNumber())
	return n
}

func (x *OrderDetail) sizeField6() (n int) {
	if x.Amount == 0 {
		return n
	}
	n += fastpb.SizeFloat(6, x.GetAmount())
	return n
}

func (x *OrderResp) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *OrderResp) sizeField1() (n int) {
	if x.Order == nil {
		return n
	}
	n += fastpb.SizeMessage(1, x.GetOrder())
	return n
}

func (x *OrderResp) sizeField2() (n int) {
	if x.OrderDetails == nil {
		return n
	}
	for i := range x.GetOrderDetails() {
		n += fastpb.SizeMessage(2, x.GetOrderDetails()[i])
	}
	return n
}

func (x *Empty) Size() (n int) {
	if x == nil {
		return n
	}
	return n
}

func (x *OrderReq) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *OrderReq) sizeField1() (n int) {
	if x.OrderId == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.GetOrderId())
	return n
}

func (x *OrderReq) sizeField2() (n int) {
	if x.List == nil {
		return n
	}
	for i := range x.GetList() {
		n += fastpb.SizeMessage(2, x.GetList()[i])
	}
	return n
}

func (x *CancelReq) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *CancelReq) sizeField1() (n int) {
	if x.OrderId == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.GetOrderId())
	return n
}

func (x *CancelReq) sizeField2() (n int) {
	if x.CancelReason == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetCancelReason())
	return n
}

var fieldIDToName_Order = map[int32]string{
	1:  "UserId",
	2:  "Number",
	3:  "Status",
	4:  "AddressBookId",
	5:  "PayMethod",
	6:  "PayStatus",
	7:  "Amount",
	8:  "Remark",
	9:  "Phone",
	10: "Address",
	11: "Username",
	12: "Consignee",
	13: "CancelReason",
	14: "RejectionReason",
	15: "ID",
	16: "ShopId",
}

var fieldIDToName_OrderDetail = map[int32]string{
	1: "Name",
	2: "Image",
	3: "OrderId",
	4: "ProductId",
	5: "Number",
	6: "Amount",
}

var fieldIDToName_OrderResp = map[int32]string{
	1: "Order",
	2: "OrderDetails",
}

var fieldIDToName_Empty = map[int32]string{}

var fieldIDToName_OrderReq = map[int32]string{
	1: "OrderId",
	2: "List",
}

var fieldIDToName_CancelReq = map[int32]string{
	1: "OrderId",
	2: "CancelReason",
}
