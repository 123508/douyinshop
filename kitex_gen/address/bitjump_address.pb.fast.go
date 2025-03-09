// Code generated by Fastpb v0.0.2. DO NOT EDIT.

package address

import (
	fmt "fmt"
	fastpb "github.com/cloudwego/fastpb"
)

var (
	_ = fmt.Errorf
	_ = fastpb.Skip
)

func (x *Address) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_Address[number], err)
}

func (x *Address) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.StreetAddress, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Address) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.City, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Address) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	x.State, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Address) fastReadField4(buf []byte, _type int8) (offset int, err error) {
	x.Country, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Address) fastReadField5(buf []byte, _type int8) (offset int, err error) {
	x.ZipCode, offset, err = fastpb.ReadInt32(buf, _type)
	return offset, err
}

func (x *Address) fastReadField6(buf []byte, _type int8) (offset int, err error) {
	x.Consignee, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Address) fastReadField7(buf []byte, _type int8) (offset int, err error) {
	x.Gender, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *Address) fastReadField8(buf []byte, _type int8) (offset int, err error) {
	x.Phone, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Address) fastReadField9(buf []byte, _type int8) (offset int, err error) {
	x.Label, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *Address) fastReadField10(buf []byte, _type int8) (offset int, err error) {
	x.IsDefault, offset, err = fastpb.ReadBool(buf, _type)
	return offset, err
}

func (x *Address) fastReadField11(buf []byte, _type int8) (offset int, err error) {
	x.AddressId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *AddressItem) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_AddressItem[number], err)
}

func (x *AddressItem) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.AddrId, offset, err = fastpb.ReadUint64(buf, _type)
	return offset, err
}

func (x *AddressItem) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	var v Address
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.Address = &v
	return offset, nil
}

func (x *AddAddressReq) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_AddAddressReq[number], err)
}

func (x *AddAddressReq) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *AddAddressReq) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	var v Address
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.Address = &v
	return offset, nil
}

func (x *AddAddressResp) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_AddAddressResp[number], err)
}

func (x *AddAddressResp) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.AddrId, offset, err = fastpb.ReadUint64(buf, _type)
	return offset, err
}

func (x *GetAddressListReq) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_GetAddressListReq[number], err)
}

func (x *GetAddressListReq) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *GetAddressListResp) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_GetAddressListResp[number], err)
}

func (x *GetAddressListResp) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	var v AddressItem
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.Address = append(x.Address, &v)
	return offset, nil
}

func (x *DeleteAddressReq) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_DeleteAddressReq[number], err)
}

func (x *DeleteAddressReq) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *DeleteAddressReq) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.AddrId, offset, err = fastpb.ReadUint64(buf, _type)
	return offset, err
}

func (x *DeleteAddressResp) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_DeleteAddressResp[number], err)
}

func (x *DeleteAddressResp) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.Res, offset, err = fastpb.ReadBool(buf, _type)
	return offset, err
}

func (x *UpdateAddressReq) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_UpdateAddressReq[number], err)
}

func (x *UpdateAddressReq) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *UpdateAddressReq) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.AddrId, offset, err = fastpb.ReadUint64(buf, _type)
	return offset, err
}

func (x *UpdateAddressReq) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	var v Address
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.Address = &v
	return offset, nil
}

func (x *UpdateAddressResp) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_UpdateAddressResp[number], err)
}

func (x *UpdateAddressResp) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.Res, offset, err = fastpb.ReadBool(buf, _type)
	return offset, err
}

func (x *SetDefaultAddressReq) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_SetDefaultAddressReq[number], err)
}

func (x *SetDefaultAddressReq) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *SetDefaultAddressReq) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.AddrId, offset, err = fastpb.ReadUint64(buf, _type)
	return offset, err
}

func (x *SetDefaultAddressResp) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_SetDefaultAddressResp[number], err)
}

func (x *SetDefaultAddressResp) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.Res, offset, err = fastpb.ReadBool(buf, _type)
	return offset, err
}

func (x *GetAddressInfoReq) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_GetAddressInfoReq[number], err)
}

func (x *GetAddressInfoReq) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *GetAddressInfoReq) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.AddrId, offset, err = fastpb.ReadUint64(buf, _type)
	return offset, err
}

func (x *GetAddressInfoResp) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
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
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_GetAddressInfoResp[number], err)
}

func (x *GetAddressInfoResp) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	var v Address
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.Addr = &v
	return offset, nil
}

func (x *Address) FastWrite(buf []byte) (offset int) {
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
	return offset
}

func (x *Address) fastWriteField1(buf []byte) (offset int) {
	if x.StreetAddress == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 1, x.GetStreetAddress())
	return offset
}

func (x *Address) fastWriteField2(buf []byte) (offset int) {
	if x.City == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetCity())
	return offset
}

func (x *Address) fastWriteField3(buf []byte) (offset int) {
	if x.State == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 3, x.GetState())
	return offset
}

func (x *Address) fastWriteField4(buf []byte) (offset int) {
	if x.Country == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 4, x.GetCountry())
	return offset
}

func (x *Address) fastWriteField5(buf []byte) (offset int) {
	if x.ZipCode == 0 {
		return offset
	}
	offset += fastpb.WriteInt32(buf[offset:], 5, x.GetZipCode())
	return offset
}

func (x *Address) fastWriteField6(buf []byte) (offset int) {
	if x.Consignee == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 6, x.GetConsignee())
	return offset
}

func (x *Address) fastWriteField7(buf []byte) (offset int) {
	if x.Gender == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 7, x.GetGender())
	return offset
}

func (x *Address) fastWriteField8(buf []byte) (offset int) {
	if x.Phone == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 8, x.GetPhone())
	return offset
}

func (x *Address) fastWriteField9(buf []byte) (offset int) {
	if x.Label == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 9, x.GetLabel())
	return offset
}

func (x *Address) fastWriteField10(buf []byte) (offset int) {
	if !x.IsDefault {
		return offset
	}
	offset += fastpb.WriteBool(buf[offset:], 10, x.GetIsDefault())
	return offset
}

func (x *Address) fastWriteField11(buf []byte) (offset int) {
	if x.AddressId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 11, x.GetAddressId())
	return offset
}

func (x *AddressItem) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *AddressItem) fastWriteField1(buf []byte) (offset int) {
	if x.AddrId == 0 {
		return offset
	}
	offset += fastpb.WriteUint64(buf[offset:], 1, x.GetAddrId())
	return offset
}

func (x *AddressItem) fastWriteField2(buf []byte) (offset int) {
	if x.Address == nil {
		return offset
	}
	offset += fastpb.WriteMessage(buf[offset:], 2, x.GetAddress())
	return offset
}

func (x *AddAddressReq) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *AddAddressReq) fastWriteField1(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.GetUserId())
	return offset
}

func (x *AddAddressReq) fastWriteField2(buf []byte) (offset int) {
	if x.Address == nil {
		return offset
	}
	offset += fastpb.WriteMessage(buf[offset:], 2, x.GetAddress())
	return offset
}

func (x *AddAddressResp) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	return offset
}

func (x *AddAddressResp) fastWriteField1(buf []byte) (offset int) {
	if x.AddrId == 0 {
		return offset
	}
	offset += fastpb.WriteUint64(buf[offset:], 1, x.GetAddrId())
	return offset
}

func (x *GetAddressListReq) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	return offset
}

func (x *GetAddressListReq) fastWriteField1(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.GetUserId())
	return offset
}

func (x *GetAddressListResp) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	return offset
}

func (x *GetAddressListResp) fastWriteField1(buf []byte) (offset int) {
	if x.Address == nil {
		return offset
	}
	for i := range x.GetAddress() {
		offset += fastpb.WriteMessage(buf[offset:], 1, x.GetAddress()[i])
	}
	return offset
}

func (x *DeleteAddressReq) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *DeleteAddressReq) fastWriteField1(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.GetUserId())
	return offset
}

func (x *DeleteAddressReq) fastWriteField2(buf []byte) (offset int) {
	if x.AddrId == 0 {
		return offset
	}
	offset += fastpb.WriteUint64(buf[offset:], 2, x.GetAddrId())
	return offset
}

func (x *DeleteAddressResp) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	return offset
}

func (x *DeleteAddressResp) fastWriteField1(buf []byte) (offset int) {
	if !x.Res {
		return offset
	}
	offset += fastpb.WriteBool(buf[offset:], 1, x.GetRes())
	return offset
}

func (x *UpdateAddressReq) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	return offset
}

func (x *UpdateAddressReq) fastWriteField1(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.GetUserId())
	return offset
}

func (x *UpdateAddressReq) fastWriteField2(buf []byte) (offset int) {
	if x.AddrId == 0 {
		return offset
	}
	offset += fastpb.WriteUint64(buf[offset:], 2, x.GetAddrId())
	return offset
}

func (x *UpdateAddressReq) fastWriteField3(buf []byte) (offset int) {
	if x.Address == nil {
		return offset
	}
	offset += fastpb.WriteMessage(buf[offset:], 3, x.GetAddress())
	return offset
}

func (x *UpdateAddressResp) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	return offset
}

func (x *UpdateAddressResp) fastWriteField1(buf []byte) (offset int) {
	if !x.Res {
		return offset
	}
	offset += fastpb.WriteBool(buf[offset:], 1, x.GetRes())
	return offset
}

func (x *SetDefaultAddressReq) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *SetDefaultAddressReq) fastWriteField1(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.GetUserId())
	return offset
}

func (x *SetDefaultAddressReq) fastWriteField2(buf []byte) (offset int) {
	if x.AddrId == 0 {
		return offset
	}
	offset += fastpb.WriteUint64(buf[offset:], 2, x.GetAddrId())
	return offset
}

func (x *SetDefaultAddressResp) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	return offset
}

func (x *SetDefaultAddressResp) fastWriteField1(buf []byte) (offset int) {
	if !x.Res {
		return offset
	}
	offset += fastpb.WriteBool(buf[offset:], 1, x.GetRes())
	return offset
}

func (x *GetAddressInfoReq) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *GetAddressInfoReq) fastWriteField1(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 1, x.GetUserId())
	return offset
}

func (x *GetAddressInfoReq) fastWriteField2(buf []byte) (offset int) {
	if x.AddrId == 0 {
		return offset
	}
	offset += fastpb.WriteUint64(buf[offset:], 2, x.GetAddrId())
	return offset
}

func (x *GetAddressInfoResp) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	return offset
}

func (x *GetAddressInfoResp) fastWriteField1(buf []byte) (offset int) {
	if x.Addr == nil {
		return offset
	}
	offset += fastpb.WriteMessage(buf[offset:], 1, x.GetAddr())
	return offset
}

func (x *Address) Size() (n int) {
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
	return n
}

func (x *Address) sizeField1() (n int) {
	if x.StreetAddress == "" {
		return n
	}
	n += fastpb.SizeString(1, x.GetStreetAddress())
	return n
}

func (x *Address) sizeField2() (n int) {
	if x.City == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetCity())
	return n
}

func (x *Address) sizeField3() (n int) {
	if x.State == "" {
		return n
	}
	n += fastpb.SizeString(3, x.GetState())
	return n
}

func (x *Address) sizeField4() (n int) {
	if x.Country == "" {
		return n
	}
	n += fastpb.SizeString(4, x.GetCountry())
	return n
}

func (x *Address) sizeField5() (n int) {
	if x.ZipCode == 0 {
		return n
	}
	n += fastpb.SizeInt32(5, x.GetZipCode())
	return n
}

func (x *Address) sizeField6() (n int) {
	if x.Consignee == "" {
		return n
	}
	n += fastpb.SizeString(6, x.GetConsignee())
	return n
}

func (x *Address) sizeField7() (n int) {
	if x.Gender == 0 {
		return n
	}
	n += fastpb.SizeUint32(7, x.GetGender())
	return n
}

func (x *Address) sizeField8() (n int) {
	if x.Phone == "" {
		return n
	}
	n += fastpb.SizeString(8, x.GetPhone())
	return n
}

func (x *Address) sizeField9() (n int) {
	if x.Label == "" {
		return n
	}
	n += fastpb.SizeString(9, x.GetLabel())
	return n
}

func (x *Address) sizeField10() (n int) {
	if !x.IsDefault {
		return n
	}
	n += fastpb.SizeBool(10, x.GetIsDefault())
	return n
}

func (x *Address) sizeField11() (n int) {
	if x.AddressId == 0 {
		return n
	}
	n += fastpb.SizeUint32(11, x.GetAddressId())
	return n
}

func (x *AddressItem) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *AddressItem) sizeField1() (n int) {
	if x.AddrId == 0 {
		return n
	}
	n += fastpb.SizeUint64(1, x.GetAddrId())
	return n
}

func (x *AddressItem) sizeField2() (n int) {
	if x.Address == nil {
		return n
	}
	n += fastpb.SizeMessage(2, x.GetAddress())
	return n
}

func (x *AddAddressReq) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *AddAddressReq) sizeField1() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.GetUserId())
	return n
}

func (x *AddAddressReq) sizeField2() (n int) {
	if x.Address == nil {
		return n
	}
	n += fastpb.SizeMessage(2, x.GetAddress())
	return n
}

func (x *AddAddressResp) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	return n
}

func (x *AddAddressResp) sizeField1() (n int) {
	if x.AddrId == 0 {
		return n
	}
	n += fastpb.SizeUint64(1, x.GetAddrId())
	return n
}

func (x *GetAddressListReq) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	return n
}

func (x *GetAddressListReq) sizeField1() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.GetUserId())
	return n
}

func (x *GetAddressListResp) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	return n
}

func (x *GetAddressListResp) sizeField1() (n int) {
	if x.Address == nil {
		return n
	}
	for i := range x.GetAddress() {
		n += fastpb.SizeMessage(1, x.GetAddress()[i])
	}
	return n
}

func (x *DeleteAddressReq) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *DeleteAddressReq) sizeField1() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.GetUserId())
	return n
}

func (x *DeleteAddressReq) sizeField2() (n int) {
	if x.AddrId == 0 {
		return n
	}
	n += fastpb.SizeUint64(2, x.GetAddrId())
	return n
}

func (x *DeleteAddressResp) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	return n
}

func (x *DeleteAddressResp) sizeField1() (n int) {
	if !x.Res {
		return n
	}
	n += fastpb.SizeBool(1, x.GetRes())
	return n
}

func (x *UpdateAddressReq) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	return n
}

func (x *UpdateAddressReq) sizeField1() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.GetUserId())
	return n
}

func (x *UpdateAddressReq) sizeField2() (n int) {
	if x.AddrId == 0 {
		return n
	}
	n += fastpb.SizeUint64(2, x.GetAddrId())
	return n
}

func (x *UpdateAddressReq) sizeField3() (n int) {
	if x.Address == nil {
		return n
	}
	n += fastpb.SizeMessage(3, x.GetAddress())
	return n
}

func (x *UpdateAddressResp) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	return n
}

func (x *UpdateAddressResp) sizeField1() (n int) {
	if !x.Res {
		return n
	}
	n += fastpb.SizeBool(1, x.GetRes())
	return n
}

func (x *SetDefaultAddressReq) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *SetDefaultAddressReq) sizeField1() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.GetUserId())
	return n
}

func (x *SetDefaultAddressReq) sizeField2() (n int) {
	if x.AddrId == 0 {
		return n
	}
	n += fastpb.SizeUint64(2, x.GetAddrId())
	return n
}

func (x *SetDefaultAddressResp) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	return n
}

func (x *SetDefaultAddressResp) sizeField1() (n int) {
	if !x.Res {
		return n
	}
	n += fastpb.SizeBool(1, x.GetRes())
	return n
}

func (x *GetAddressInfoReq) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *GetAddressInfoReq) sizeField1() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeUint32(1, x.GetUserId())
	return n
}

func (x *GetAddressInfoReq) sizeField2() (n int) {
	if x.AddrId == 0 {
		return n
	}
	n += fastpb.SizeUint64(2, x.GetAddrId())
	return n
}

func (x *GetAddressInfoResp) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	return n
}

func (x *GetAddressInfoResp) sizeField1() (n int) {
	if x.Addr == nil {
		return n
	}
	n += fastpb.SizeMessage(1, x.GetAddr())
	return n
}

var fieldIDToName_Address = map[int32]string{
	1:  "StreetAddress",
	2:  "City",
	3:  "State",
	4:  "Country",
	5:  "ZipCode",
	6:  "Consignee",
	7:  "Gender",
	8:  "Phone",
	9:  "Label",
	10: "IsDefault",
	11: "AddressId",
}

var fieldIDToName_AddressItem = map[int32]string{
	1: "AddrId",
	2: "Address",
}

var fieldIDToName_AddAddressReq = map[int32]string{
	1: "UserId",
	2: "Address",
}

var fieldIDToName_AddAddressResp = map[int32]string{
	1: "AddrId",
}

var fieldIDToName_GetAddressListReq = map[int32]string{
	1: "UserId",
}

var fieldIDToName_GetAddressListResp = map[int32]string{
	1: "Address",
}

var fieldIDToName_DeleteAddressReq = map[int32]string{
	1: "UserId",
	2: "AddrId",
}

var fieldIDToName_DeleteAddressResp = map[int32]string{
	1: "Res",
}

var fieldIDToName_UpdateAddressReq = map[int32]string{
	1: "UserId",
	2: "AddrId",
	3: "Address",
}

var fieldIDToName_UpdateAddressResp = map[int32]string{
	1: "Res",
}

var fieldIDToName_SetDefaultAddressReq = map[int32]string{
	1: "UserId",
	2: "AddrId",
}

var fieldIDToName_SetDefaultAddressResp = map[int32]string{
	1: "Res",
}

var fieldIDToName_GetAddressInfoReq = map[int32]string{
	1: "UserId",
	2: "AddrId",
}

var fieldIDToName_GetAddressInfoResp = map[int32]string{
	1: "Addr",
}
