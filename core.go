package gorm_expand

import (
	"github.com/Merry74751/yutool/anyutil"
	"gorm.io/gorm"
)

type Mapper[T any] struct {
	Db *gorm.DB
}

func (m Mapper[T]) GetById(id any) T {
	value := new(T)
	field := anyutil.StructField(value, 0)
	column := getColumn(field)
	sql := generateSelect[T]()
	m.Db.Exec(sql).Where(column, id).Find(value)
	return *value
}
