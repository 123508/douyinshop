package cursor

import "github.com/123508/douyinshop/pkg/component/condition"

type Cursor interface {
	EncodeCursor() (string, error)
	DecodeCursor(string) error
	ToCondition(bool) condition.Condition
}
