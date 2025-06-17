package sort

import (
	"fmt"
	"regexp"
	"strings"
)

//对添加排序字段的实现

type SortOrder string

const (
	ASC  SortOrder = "asc"
	DESC SortOrder = "desc"
)

// 仅处理末尾排序关键字（asc/desc），允许结尾有空格，不区分大小写
var orderBySuffixRegexp = regexp.MustCompile(`(?i)\s+(asc|desc)\s*$`)

// SortExpr 元样式
type SortExpr struct {
	FieldName string
	Order     SortOrder
}

func (s SortExpr) ToSortItem() string {

	if s.Order != ASC && s.Order != DESC {
		return fmt.Sprintf("%s %s", s.FieldName, ASC)
	}

	return fmt.Sprintf("%s %s", s.FieldName, s.Order)
}

func (s SortExpr) Reverse() string {
	if s.Order != ASC && s.Order != DESC {
		return fmt.Sprintf("%s %s", s.FieldName, DESC)
	}

	if s.Order == ASC {
		return fmt.Sprintf("%s %s", s.FieldName, DESC)
	}

	return fmt.Sprintf("%s %s", s.FieldName, ASC)
}

type RowExpr struct {
	Expr string
}

func (r RowExpr) ToSortItem() string {
	return r.Expr
}

func (r RowExpr) Reverse() string {
	// 倒序/升序互换
	return orderBySuffixRegexp.ReplaceAllStringFunc(r.Expr, func(s string) string {
		// s like " desc" or " asc" (可能带空格)
		if strings.Contains(strings.ToLower(s), "asc") {
			return strings.Replace(s, "asc", "desc", 1)
		}
		return strings.Replace(s, "desc", "asc", 1)
	})
}

type SortOnMySQL struct {
	Id        Sort
	CreatedAt Sort
	Sorts     []Sort
}

func NewSortOnMySQL() *SortOnMySQL {
	return &SortOnMySQL{
		Id:        SortExpr{FieldName: "id", Order: ASC},
		CreatedAt: SortExpr{FieldName: "created_at", Order: ASC},
		Sorts:     make([]Sort, 0),
	}
}

func (s *SortOnMySQL) SetIdName(id string) *SortOnMySQL {
	s.Id = SortExpr{FieldName: id}
	return s
}

func (s *SortOnMySQL) SetCreatedTimeName(createdTime string) *SortOnMySQL {
	s.CreatedAt = SortExpr{FieldName: createdTime}
	return s
}

func (s *SortOnMySQL) SetSortsItems(sorts []Sort) *SortOnMySQL {
	s.Sorts = make([]Sort, len(sorts))
	copy(s.Sorts, sorts)
	return s
}

func (s *SortOnMySQL) ToSorts(reverse bool) []string {
	var parts []string

	if s == nil {
		return parts
	}

	if s.Sorts == nil {
		s.Sorts = make([]Sort, 0)
	}

	if s.Id == nil {
		s.Id = SortExpr{}
	}

	if s.CreatedAt == nil {
		s.CreatedAt = SortExpr{}
	}

	//这里是正常字段的添加
	for _, v := range s.Sorts {
		if reverse {
			parts = append(parts, v.Reverse())
		} else {
			parts = append(parts, v.ToSortItem())
		}
	}

	//这里是id和createdAt字段的添加
	if reverse {
		parts = append(parts, s.Id.Reverse(), s.CreatedAt.Reverse())
	} else {
		parts = append(parts, s.Id.ToSortItem(), s.CreatedAt.ToSortItem())
	}

	return parts
}
