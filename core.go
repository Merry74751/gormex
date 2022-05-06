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

type queryCondition struct {
	params map[string]any
}

func NewQueryCondition() queryCondition {
	q := queryCondition{}
	q.params = make(map[string]any)
	return q
}

func EmptyCondition() queryCondition {
	return queryCondition{}
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

type Page struct {
	Current  int
	PageSize int
}

type Mapper[T any] struct {
	db *gorm.DB
}

func (m Mapper[T]) GetDb() *gorm.DB {
	return m.db
}

func (m Mapper[T]) SetDb(db *gorm.DB) {
	m.db = db
}

func (m Mapper[T]) Insert(t T) error {
	sql, params := generateInsert(t)
	err := m.db.Exec(sql, params...).Error
	return err
}

func (m Mapper[T]) GetById(id any) (T, error) {
	value := new(T)
	columns := Columns(value)
	m.db = m.db.Select(columns)

	err := m.db.Where(columns[0]+"=?", id).Find(value).Error
	return *value, err
}

func (m Mapper[T]) DeleteById(id any) error {
	t := new(T)
	field := anyutil.StructField(t, 0)
	column := getColumn(field)
	return m.db.Where(column+"=?", id).Delete(field).Error
}

func (m Mapper[T]) UpdateById(t T) error {
	err := m.db.Updates(t).Error
	return err
}

func (m Mapper[T]) Get(t T, queryCondition queryCondition) (T, error) {
	result := new(T)
	columns := Columns(result)
	m.db = m.db.Select(columns)

	entityToCondition(t, m.db)
	parseCondition(queryCondition, m.db)

	err := m.db.Find(result).Error
	return *result, err
}

func (m Mapper[T]) List() ([]T, error) {
	columns := Columns(new(T))
	m.db = m.db.Select(columns)

	var result []T
	err := m.db.Find(&result).Error
	return result, err
}

func (m Mapper[T]) ListByCondition(t T, condition queryCondition) ([]T, error) {
	columns := Columns(new(T))
	m.db = m.db.Select(columns)

	var result []T
	entityToCondition(t, m.db)
	parseCondition(condition, m.db)

	err := m.db.Find(&result).Error
	return result, err
}

func (m Mapper[T]) Page(page Page) ([]T, int64, error) {
	columns := Columns(new(T))
	m.db = m.db.Select(columns)

	var total int64
	m.db.Count(&total)

	var result []T
	current := (page.Current - 1) * page.PageSize
	err := m.db.Offset(current).Limit(page.PageSize).Find(&result).Error
	return result, total, err
}

func entityToCondition(t any, db *gorm.DB) {
	if anyutil.IsNotNil(t) {
		fields := anyutil.Fields(t)
		for i, field := range fields {
			if field.IsZero() {
				continue
			}
			structField := anyutil.StructField(t, i)
			column := getColumn(structField)
			db = db.Where(column+"=?", field.Interface())
		}
	}
}

func parseCondition(condition queryCondition, db *gorm.DB) {
	if anyutil.IsNotNil(condition) {
		for key, value := range condition.params {
			db = db.Where(key, value)
		}
	}
}
