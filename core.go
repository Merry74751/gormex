package gorm_expand

import (
	"github.com/Merry74751/yutool/anyutil"
	"gorm.io/gorm"
)

type condition interface {
	// Eq 等于
	Eq(column string, value any)
	// Ne 不等于
	Ne(column string, value any)
	// Lt 小于
	Lt(column string, value any)
	// Le 小于等于
	Le(column string, value any)
	// Gt 大于
	Gt(column string, value any)
	// Ge 大于等于
	Ge(column string, value any)
}

type queryCondition[T any] struct {
	params map[string]any
}

func NewQueryCondition[T any]() queryCondition[T] {
	q := queryCondition[T]{}
	q.params = make(map[string]any)
	return q
}

func (q queryCondition) Eq(column string, value any) {
	key := column + "=?"
	q.params[key] = value
}

func (q queryCondition) Ne(column string, value any) {
	key := column + "<> ? "
	q.params[key] = value
}

func (q queryCondition) Lt(column string, value any) {
	key := column + "< ?"
	q.params[key] = value
}

func (q queryCondition) Le(column string, value any) {
	key := column + "<= ?"
	q.params[key] = value
}

func (q queryCondition) Gt(column string, value any) {
	key := column + "> ?"
	q.params[key] = value
}

func (q queryCondition) Ge(column string, value any) {
	key := column + ">= ?"
	q.params[key] = value
}

type Mapper[T any] struct {
	Db *gorm.DB
}

func (m Mapper[T]) Insert(t T) {
	sql, params := generateInsert(t)
	m.Db.Exec(sql, params...)
}

func (m Mapper[T]) GetById(id any) T {
	value := new(T)
	field := anyutil.StructField(value, 0)
	column := getColumn(field)
	sql := generateSelect[T]()
	m.Db.Raw(sql).Where(column, id).Find(value)
	return *value
}

func (m Mapper[T]) Get(queryCondition queryCondition[T]) T {
	sql := generateSelect[T]()
	m.Db = m.Db.Raw(sql)
	for key, value := range queryCondition.params {
		m.Db = m.Db.Where(key, value)
	}
	t := new(T)
	m.Db.Find(t)
	return *t
}
