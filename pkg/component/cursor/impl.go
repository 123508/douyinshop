package cursor

import (
	"encoding/base64"
	"github.com/123508/douyinshop/pkg/component/condition"
	"github.com/123508/douyinshop/pkg/component/serializer"
	"github.com/123508/douyinshop/pkg/util"
	"time"
)

type Base struct {
	FieldName string
	Value     interface{}
	IsDesc    bool
}

type StandCursor struct {
	Id          Base
	CreatedTime Base
	Values      []Base
	Serializer  *serializer.SerializerWrapper
	RightBound  bool
	LeftBound   bool
}

func NewBase(FieldName string, value interface{}, isDesc bool) Base {
	return Base{FieldName: FieldName, Value: value, IsDesc: isDesc}
}

func (c *StandCursor) SetRightBound() {
	c.RightBound = true
}

func (c *StandCursor) SetLeftBound() {
	c.LeftBound = true
}

func (c *StandCursor) IsRightBound() bool {
	return c.RightBound
}

func (c *StandCursor) IsLeftBound() bool {
	return c.LeftBound
}

// EncodeCursor 编码
func (c *StandCursor) EncodeCursor() (string, error) {
	serial := c.Serializer
	if serial == nil {
		serial = &serializer.SerializerWrapper{Strategy: "json"}
	}

	b, err := serial.Serialize(c)
	if err != nil {
		return "", err
	}
	// 2. json字节转base64字符串
	return base64.StdEncoding.EncodeToString(b), nil
}

// DecodeCursor 解码
func (c *StandCursor) DecodeCursor(s string) error {

	serial := c.Serializer

	if serial == nil {
		serial = &serializer.SerializerWrapper{Strategy: "json"}
	}

	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		util.LogError("转换错误", "DecodeCursor", "", err)
		return err
	}

	return serial.Deserialize(b, c)
}

// ToCondition 转换为查询条件
func (c *StandCursor) ToCondition(reverse bool) condition.Condition {

	cpy := *c // 浅拷贝一份

	if cpy.Id.FieldName == "" {
		cpy.Id.FieldName = "id"
	}

	if cpy.Id.Value == nil {
		cpy.Id.Value = 0
	}

	if cpy.CreatedTime.FieldName == "" {
		cpy.CreatedTime.FieldName = "created_at"
	}

	if cpy.CreatedTime.Value == nil {
		if reverse {
			cpy.CreatedTime.Value = time.Now()
		} else {
			cpy.CreatedTime.Value = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		}
	}

	var IdOp string
	var CreatedAtOp string

	if !reverse { //执行原查询逻辑
		IdOp = ">="
		if cpy.Id.IsDesc {
			IdOp = "<="
		}

		CreatedAtOp = ">="
		if cpy.CreatedTime.IsDesc {
			CreatedAtOp = "<="
		}
	} else { //执行翻转查询逻辑
		IdOp = "<="
		if cpy.Id.IsDesc {
			IdOp = ">="
		}

		CreatedAtOp = "<="
		if cpy.CreatedTime.IsDesc {
			CreatedAtOp = ">="
		}
	}

	// (time >= lastTime AND id >= lastId)
	cond := condition.And{
		Conditions: []condition.Condition{
			condition.Expr{Field: cpy.CreatedTime.FieldName, Operator: CreatedAtOp, Value: cpy.CreatedTime.Value},
			condition.Expr{Field: cpy.Id.FieldName, Operator: IdOp, Value: cpy.Id.Value},
		},
	}

	for _, v := range c.Values {
		if v.FieldName != "" {
			cond.Conditions = append(cond.Conditions, condition.Expr{Field: v.FieldName, Operator: v.FieldName, Value: v})
		}
	}

	return cond
}

func (c *StandCursor) Clone() *StandCursor {
	return &StandCursor{
		Id:          c.Id,
		CreatedTime: c.CreatedTime,
		Values:      c.Values,
		Serializer:  c.Serializer.Clone(),
		RightBound:  c.RightBound,
		LeftBound:   c.LeftBound,
	}
}
