package cursor

import (
	"github.com/123508/douyinshop/pkg/component/serializer"
	"time"
)

//游标构建器

type CursorStructBuilder struct {
	cur *StandCursor
}

func NewCursorStructBuilder() *CursorStructBuilder {
	return &CursorStructBuilder{
		cur: &StandCursor{
			Id:          NewBase("id", 0, false),
			CreatedTime: NewBase("created_at", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), false),
			Serializer:  &serializer.SerializerWrapper{Strategy: "json"},
			RightBound:  false,
		},
	}
}

func (c *CursorStructBuilder) RightBound() *CursorStructBuilder {
	c.cur.SetRightBound()
	return c
}

func (c *CursorStructBuilder) LeftBound() *CursorStructBuilder {
	c.cur.SetLeftBound()
	return c
}

func (c *CursorStructBuilder) WithBaseId(id Base) *CursorStructBuilder {
	c.cur.Id = id
	return c
}

func (c *CursorStructBuilder) WithId(field string, value interface{}, isDesc bool) *CursorStructBuilder {
	c.cur.Id = NewBase(field, value, isDesc)
	return c
}

func (c *CursorStructBuilder) WithBaseCreatedTime(createdTime Base) *CursorStructBuilder {
	c.cur.CreatedTime = createdTime
	return c
}

func (c *CursorStructBuilder) WithCreatedTime(field string, value interface{}, isDesc bool) *CursorStructBuilder {
	c.cur.CreatedTime = NewBase(field, value, isDesc)
	return c
}

func (c *CursorStructBuilder) WithSerializer(s serializer.SerializerType) *CursorStructBuilder {
	c.cur.Serializer = serializer.NewSerializerWrapper(s)
	return c
}

func (c *CursorStructBuilder) WithSerializerWrapper(s *serializer.SerializerWrapper) *CursorStructBuilder {

	c.cur.Serializer = s.Clone()
	return c
}

func (c *CursorStructBuilder) Build() *StandCursor {
	return c.cur.Clone()
}

func (c *CursorStructBuilder) Clone() *CursorStructBuilder {
	return &CursorStructBuilder{
		cur: c.cur.Clone(),
	}
}
