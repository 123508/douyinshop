package queryexec

import (
	"errors"
	"github.com/123508/douyinshop/pkg/component/condition"
	"github.com/123508/douyinshop/pkg/component/cursor"
	"github.com/123508/douyinshop/pkg/component/pub"
	"github.com/123508/douyinshop/pkg/component/serializer"
	"github.com/123508/douyinshop/pkg/component/sort"
	"github.com/123508/douyinshop/pkg/util"
	"gorm.io/gorm"
)

type QueryExecStructMySQL[Item pub.ItemType[Id], Id pub.IntegerNumber] struct {
	RawConditions condition.Condition
	PageSize      int
	Sort          *sort.SortOnMySQL
	Database      *gorm.DB
	Cur           string
	cur           *cursor.StandCursor
}

func NewQueryExecOnMySQL[Item pub.ItemType[Id], Id pub.IntegerNumber](rawConditions condition.Condition, PageSize int, Sort *sort.SortOnMySQL, db *gorm.DB, Cursor string) *QueryExecStructMySQL[Item, Id] {
	return &QueryExecStructMySQL[Item, Id]{
		RawConditions: rawConditions,
		PageSize:      PageSize,
		Sort:          Sort,
		Database:      db,
		Cur:           Cursor,
	}
}

// 页数校验
func (q *QueryExecStructMySQL[Item, Id]) validPageSize() bool {

	innerPageSize := q.PageSize

	if innerPageSize < 1 || innerPageSize > 500 {
		util.LogError("查询页数不被允许", "validPageSize", "", nil)
		return false
	}
	return true
}

// 添加过滤条件
func (q *QueryExecStructMySQL[Item, Id]) withFilter(db *gorm.DB) *gorm.DB {

	innerCond := q.RawConditions

	// 过滤条件
	if innerCond != nil {
		whereSQL, whereArgs := innerCond.ToSQL()
		db = db.Where(whereSQL, whereArgs...)
	}
	return db
}

// 添加游标部分
func (q *QueryExecStructMySQL[Item, Id]) withCursor(db *gorm.DB, reverse bool) *gorm.DB {

	innerCur := q.cur

	//游标查询
	if innerCur != nil {
		whereSQL, whereArgs := innerCur.ToCondition(reverse).ToSQL()
		db = db.Where(whereSQL, whereArgs...)
	}
	return db
}

// 添加排序
func (q *QueryExecStructMySQL[Item, Id]) withSort(db *gorm.DB, reverse bool) *gorm.DB {

	// 排序
	if q.Sort == nil {
		q.Sort = sort.NewSortOnMySQL()
	}

	innerSort := q.Sort

	Sorts := innerSort.ToSorts(reverse)

	if len(Sorts) != 0 {
		for _, v := range Sorts {
			db = db.Order(v)
		}
	}

	return db
}

// 边界处理函数
func (q *QueryExecStructMySQL[Item, Id]) isBoundary(cur *cursor.StandCursor, preQuery bool) (bool, string, error) {

	//构建基础数据
	EndCur := cursor.NewCursorStructBuilder().
		WithBaseId(cur.Id).
		WithBaseCreatedTime(cur.CreatedTime).
		Build()

	//左边界判定
	if cur.IsRightBound() && !preQuery {
		EndCur.SetRightBound()
	}

	//右边界判定
	if cur.IsLeftBound() && preQuery {
		EndCur.SetLeftBound()
	}

	//边界情况处理
	if (cur.IsRightBound() && !preQuery) || (cur.IsLeftBound() && preQuery) {
		encodeCursor, err := EndCur.EncodeCursor()
		if err != nil {
			util.LogError("封装游标错误", "QueryExec", "", err)
		}
		return true, encodeCursor, err
	}

	return false, "", nil
}

func (q *QueryExecStructMySQL[Item, Id]) processResults(QueryCur *cursor.StandCursor, list []Item, preQuery bool) ([]Item, string, int, error) {
	if len(list) == 0 {
		return nil, "", 0, nil
	} else if len(list) == 1 {

		EndCur := cursor.NewCursorStructBuilder().
			WithBaseId(QueryCur.Id).
			WithBaseCreatedTime(QueryCur.CreatedTime).
			Build()

		if preQuery {
			EndCur.SetLeftBound()
		} else {
			EndCur.SetRightBound()
		}

		encodeCursor, err := EndCur.EncodeCursor()
		if err != nil {
			util.LogError("封装游标错误", "QueryExec", "", err)
		}
		return list, encodeCursor, 0, err
	}

	curId := list[len(list)-1].GetID()
	curCreateTIme := list[len(list)-1].GetCreatedTime()

	Cur := cursor.NewCursorStructBuilder().
		WithId(
			QueryCur.Id.FieldName,
			curId,
			QueryCur.Id.IsDesc,
		).
		WithCreatedTime(
			QueryCur.CreatedTime.FieldName,
			curCreateTIme,
			QueryCur.CreatedTime.IsDesc,
		)

	if len(list) == q.PageSize+1 {
		//丢弃第一条数据
		list = list[:len(list)-1]
	} else {
		if preQuery {
			Cur.LeftBound()
		} else {
			Cur.RightBound()
		}
	}

	if QueryCur.Serializer != nil {
		Cur = Cur.WithSerializerWrapper(QueryCur.Serializer)
	} else {
		Cur = Cur.WithSerializer(serializer.JSON)
	}

	encodeCursor, err := Cur.Build().EncodeCursor()

	if err != nil {
		util.LogError("封装游标错误", "QueryExec", "", err)
		return nil, "", 0, err
	}

	if preQuery {
		for pre, aft := 0, len(list)-1; pre <= aft; pre, aft = pre+1, aft-1 {
			list[pre], list[aft] = list[aft], list[pre]
		}
	}

	return list, encodeCursor, len(list), nil
}

func (q *QueryExecStructMySQL[Item, Id]) QueryExec(preQuery bool) (List []Item, Cursor string, Length int, Err error) {

	q.cur = cursor.NewCursorStructBuilder().Build()

	if err := q.cur.DecodeCursor(q.Cur); err != nil {
		q.cur = cursor.NewCursorStructBuilder().Build()
	}

	QueryCur := q.cur

	isBoundary, ResCur, err := q.isBoundary(QueryCur, preQuery)

	if isBoundary {
		return nil, ResCur, 0, err
	}

	if !q.validPageSize() {
		return nil, "", 0, errors.New("查询不允许")
	}

	db := q.Database.Model(new(Item))
	//添加过滤条件
	db = q.withFilter(db)
	//添加游标
	db = q.withCursor(db, preQuery)
	//添加排序
	db = q.withSort(db, preQuery)
	// 分页
	db = db.Limit(q.PageSize + 1)

	// 查询
	var list []Item

	if err := db.Find(&list).Error; err != nil {
		util.LogError("查询出错", "QueryExec", "", err)
		return nil, "", 0, err
	}

	return q.processResults(QueryCur, list, preQuery)
}
