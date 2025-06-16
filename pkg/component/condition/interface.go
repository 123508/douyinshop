package condition

type Condition interface {
	ToSQL() (string, []interface{})
}
