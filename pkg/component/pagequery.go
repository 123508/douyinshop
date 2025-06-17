package component

import (
	"context"
	"github.com/123508/douyinshop/pkg/component/pub"
	"github.com/123508/douyinshop/pkg/component/queryexec"
	"gorm.io/gorm"
)

type PageQuery[Item pub.ItemType[Id], Id pub.IntegerNumber] struct {
	Database  *gorm.DB
	QueryExec *queryexec.QueryExecStructMySQL[Item, Id]
	Cursor    string
}

type PageResult[Item pub.ItemType[Id], Id pub.IntegerNumber] struct {
	List       []Item
	HasMore    bool
	NextCursor string
}

// CursorPaginate 需要改造
func (p *PageQuery[Item, Id]) CursorPaginate(ctx context.Context, cur string) (res PageResult[Item, Id], err error) {

	if p.QueryExec == nil {
		p.QueryExec = &queryexec.QueryExecStructMySQL[Item, Id]{
			RawConditions: nil,
			PageSize:      100,
			Sort:          nil,
			Database:      nil,
			Cur:           cur,
		}
	}

	p.QueryExec.Cur = p.Cursor

	exec, nextCur, _, err := p.QueryExec.QueryExec(false)

	if err != nil {
		return PageResult[Item, Id]{}, err
	}

	return PageResult[Item, Id]{
		List:       exec,
		HasMore:    len(exec) == p.QueryExec.PageSize,
		NextCursor: nextCur,
	}, nil
}
