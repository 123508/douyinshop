package pub

import "time"

type IntegerNumber interface {
	~int | ~uint | ~int64 | ~uint64 | ~int32 | ~uint32 | ~int16 | ~uint16 | ~int8 | ~uint8
}

type ItemType[E IntegerNumber] interface {
	GetID() E
	GetCreatedTime() time.Time
}
