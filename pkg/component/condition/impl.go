package condition

import (
	"fmt"
	"reflect"
	"strings"
)

// Expr 查询基础表达式,将 id = 1 分解为了  id , = , 1 三个部分
type Expr struct {
	Field    string
	Operator string
	Value    interface{}
}

//将表达式转换为sql查询

func (e Expr) ToSQL() (string, []interface{}) {
	switch strings.ToLower(e.Operator) {
	//处理 in
	case "in":
		val := reflect.ValueOf(e.Value)
		if val.Kind() == reflect.Slice {
			if val.Len() == 0 {
				return "1=0", nil // 防止in空数组导致SQL错误
			}
			ps := make([]string, val.Len())
			params := make([]interface{}, val.Len())
			for i := 0; i < val.Len(); i++ {
				ps[i] = "?"
				params[i] = val.Index(i).Interface()
			}
			return fmt.Sprintf("%s IN (%s)", e.Field, strings.Join(ps, ",")), params
		}
		return fmt.Sprintf("%s IN (?)", e.Field), []interface{}{e.Value}
	//处理 between
	case "between":
		val := reflect.ValueOf(e.Value)
		if val.Kind() == reflect.Slice && val.Len() == 2 {
			return fmt.Sprintf("%s BETWEEN ? AND ?", e.Field), []interface{}{val.Index(0).Interface(), val.Index(1).Interface()}
		}
		return "1=0", nil // 不合法的between
	//处理 like
	case "like":
		return fmt.Sprintf("%s LIKE ?", e.Field), []interface{}{"%" + fmt.Sprintf("%v", e.Value) + "%"}
	//处理 is null和is not null
	case "is null", "is not null":
		return fmt.Sprintf("%s %s", e.Field, strings.ToUpper(e.Operator)), nil
	}
	//operator白名单校验
	allowedOps := map[string]struct{}{"=": {}, "!=": {}, ">": {}, ">=": {}, "<": {}, "<=": {}}
	if _, ok := allowedOps[e.Operator]; !ok {
		return "1=0", nil
	}
	return fmt.Sprintf("%s %s ?", e.Field, e.Operator), []interface{}{e.Value}
}

type And struct {
	Conditions []Condition
}

func (a And) ToSQL() (string, []interface{}) {
	var parts []string
	var params []interface{}
	for _, c := range a.Conditions {
		sql, p := c.ToSQL()
		parts = append(parts, sql)
		params = append(params, p...)
	}
	return "(" + strings.Join(parts, " AND ") + ")", params
}

type Or struct {
	Conditions []Condition
}

func (o Or) ToSQL() (string, []interface{}) {
	var parts []string
	var params []interface{}
	for _, c := range o.Conditions {
		sql, p := c.ToSQL()
		parts = append(parts, sql)
		params = append(params, p...)
	}
	return "(" + strings.Join(parts, " OR ") + ")", params
}
