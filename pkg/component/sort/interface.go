package sort

//这个接口用来添加排序字段

type Sort interface {
	ToSortItem() string
	Reverse() string
}
