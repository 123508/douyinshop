package condition

//条件构建器,使用时默认会线性向后拼接,一般不用做复杂场景下

type ConditionBuilder struct {
	cond Condition
}

func NewConditionBuilder() *ConditionBuilder {
	return &ConditionBuilder{}
}

func (b *ConditionBuilder) And(Field string, Operator string, Value interface{}) *ConditionBuilder {

	c := Expr{
		Field:    Field,
		Operator: Operator,
		Value:    Value,
	}

	if b.cond == nil {
		b.cond = c
	} else {
		// 如果b.cond本身是And，直接追加，避免嵌套
		if and, ok := b.cond.(And); ok {
			and.Conditions = append(and.Conditions, c)
			b.cond = and
		} else {
			b.cond = And{Conditions: []Condition{b.cond, c}}
		}
	}
	return b
}

func (b *ConditionBuilder) Or(Field string, Operator string, Value interface{}) *ConditionBuilder {

	c := Expr{
		Field:    Field,
		Operator: Operator,
		Value:    Value,
	}

	if b.cond == nil {
		b.cond = c
	} else {
		if or, ok := b.cond.(Or); ok {
			or.Conditions = append(or.Conditions, c)
			b.cond = or
		} else {
			b.cond = Or{Conditions: []Condition{b.cond, c}}
		}
	}
	return b
}

func (b *ConditionBuilder) Build() Condition {

	if b.cond == nil {
		return Expr{"1", "=", "1"}
	}

	return b.cond
}
