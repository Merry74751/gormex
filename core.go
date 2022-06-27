package gorm_expand

import (
	"fmt"
	"github.com/Merry74751/yutool/anyutil"
	"github.com/Merry74751/yutool/common"
	"github.com/Merry74751/yutool/str"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

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

func (m *Mapper[T]) SetDb(db *gorm.DB) {
	m.db = db
}

func (m Mapper[T]) Insert(t T) error {
	setInsertField(&t)
	sql, params := generateInsert(t)
	err := m.db.Exec(sql, params...).Error
	return err
}

func (m Mapper[T]) GetById(id any) (T, error) {
	value := new(T)
	columns := Columns(value)
	tx := m.db.Select(columns)

	err := tx.Where(columns[0]+"=?", id).Find(value).Error
	return *value, err
}

func (m Mapper[T]) DeleteById(id any) error {
	t := new(T)
	field := anyutil.StructField(t, 0)
	column := getColumn(field)
	return m.db.Where(column+"=?", id).Delete(t).Error
}

func (m Mapper[T]) UpdateById(t T) error {
	setUpdateField(&t)
	err := m.db.Updates(t).Error
	return err
}

func (m Mapper[T]) Get(t T) (T, error) {
	result := new(T)
	columns := Columns(result)
	tx := m.db.Select(columns)

	entityToCondition(t, tx)

	err := tx.Find(result).Error
	if anyutil.IsNil(result) {
		err = gorm.ErrRecordNotFound
	}
	return *result, err
}

func (m Mapper[T]) List() ([]T, error) {
	columns := Columns(new(T))
	tx := m.db.Select(columns)

	var result []T
	err := tx.Find(&result).Error
	return result, err
}

func (m Mapper[T]) ListByCondition(t T) ([]T, error) {
	columns := Columns(new(T))
	tx := m.db.Select(columns)

	var result []T
	entityToCondition(t, tx)

	err := tx.Find(&result).Error
	return result, err
}

func (m Mapper[T]) Page(page Page) ([]T, int64, error) {
	t := new(T)
	var total int64
	m.db.Model(t).Count(&total)

	columns := Columns(t)
	tx := m.db.Select(columns)

	var result []T
	current := (page.Current - 1) * page.PageSize
	err := tx.Offset(current).Limit(page.PageSize).Find(&result).Error
	return result, total, err
}

func (m Mapper[T]) PageByCondition(t T, page Page) ([]T, int64, error) {
	v := new(T)

	var total int64
	m.db.Model(v).Count(&total)

	columns := Columns(v)
	tx := m.db.Select(columns)

	entityToCondition(t, tx)

	var result []T
	current := (page.Current - 1) * page.PageSize
	err := tx.Offset(current).Limit(page.PageSize).Find(&result).Error
	return result, total, err
}

func entityToCondition(t any, db *gorm.DB) {
	if anyutil.IsNotNil(t) {
		v := anyutil.Value(t)
		valueToCondition(v, db)
	}
}

func valueToCondition(v reflect.Value, db *gorm.DB) {
	typ := v.Type()
	numField := v.NumField()
	for i := 0; i < numField; i++ {
		field := v.Field(i)
		if field.IsZero() {
			continue
		}
		if field.Kind() == reflect.Struct {
			valueToCondition(field, db)
		}
		structField := typ.Field(i)
		name := structField.Name
		tagColumn := parseTag(structField.Tag)
		if strings.HasSuffix(name, "GE") {
			chooseColumn(name, tagColumn, ">= ?", 2, field.Interface(), db)
		} else if strings.HasSuffix(name, "GT") {
			chooseColumn(name, tagColumn, "> ?", 2, field.Interface(), db)
		} else if strings.HasSuffix(name, "LE") {
			chooseColumn(name, tagColumn, "<= ?", 2, field.Interface(), db)
		} else if strings.HasSuffix(name, "LT") {
			chooseColumn(name, tagColumn, "< ?", 2, field.Interface(), db)
		} else if strings.HasSuffix(name, "NE") {
			chooseColumn(name, tagColumn, "<> ?", 2, field.Interface(), db)
		} else if strings.HasSuffix(name, "LikeLeft") {
			value := fmt.Sprint(field.Interface()) + "%"
			chooseColumn(name, tagColumn, "like ?", 8, value, db)
		} else {
			chooseColumn(name, tagColumn, "= ?", 0, field.Interface(), db)
		}
	}
}

func chooseColumn(name, tagColumn, expression string, index int, value any, db *gorm.DB) {
	if str.IsNotEmpty(tagColumn) {
		db = db.Where(tagColumn+expression, value)
	} else {
		name = toUnderLineCase(name[0 : len(name)-index])
		db = db.Where(name+expression, value)
	}
}

func setInsertField(t any) {
	v := anyutil.Value(t)
	id := v.FieldByName("Id")
	if (id != reflect.Value{}) {
		snowflakeID := uint(common.GetSnowflakeID())
		value := reflect.ValueOf(snowflakeID)
		id.Set(value)
	}

	createTime := v.FieldByName("CreateTime")
	if (createTime != reflect.Value{}) {
		createTime.Set(anyutil.Value(time.Now()))
	}
}

func setUpdateField(t any) {
	v := anyutil.Value(t)
	updateTime := v.FieldByName("UpdateTime")

	if (updateTime != reflect.Value{}) {
		updateTime.Set(anyutil.Value(time.Now()))
	}
}
