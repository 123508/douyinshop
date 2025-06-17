package condition

// Condition 查询条件父类
type Condition interface {
	ToSQL() (string, []interface{})
}
