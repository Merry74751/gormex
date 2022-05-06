package gorm_expand

import (
	"github.com/Merry74751/yutool/anyutil"
	"github.com/Merry74751/yutool/str"
	"reflect"
	"strings"
)

func getColumn(field reflect.StructField) string {
	column := parseTag(field.Tag)
	if str.IsEmpty(column) {
		column = strings.ToLower(field.Name)
	}
	return column
}

func parseTag(tag reflect.StructTag) string {
	v := tag.Get("gorm")
	if str.IsNotEmpty(v) {
		s := strings.Split(v, ";")
		for _, item := range s {
			s1 := strings.Split(item, ":")
			if s1[0] == "column" {
				return s1[1]
			}
		}
	}
	return ""
}

func Columns(v any) []string {
	t := anyutil.Type(v)
	numField := t.NumField()
	columns := make([]string, numField)
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		column := getColumn(field)
		columns[i] = column
	}
	return columns
}

func TableName(v any) string {
	value := anyutil.Value(v)
	method := value.MethodByName("TableName")

	if method.IsValid() {
		call := method.Call(nil)
		return call[0].Interface().(string)
	}

	structName := reflect.TypeOf(v).String()
	index := strings.Index(structName, ".")
	tableName := str.SubString(structName, index+1, -1)
	return str.ConvertUnderline(tableName)
}

func ColumnsNotNil(v any) ([]string, int, []any) {
	value := anyutil.Value(v)
	numField := value.NumField()

	length := 0
	columns := make([]string, numField)
	values := make([]any, numField)
	types := anyutil.Type(v)

	for i := 0; i < numField; i++ {
		field := value.Field(i)
		if field.IsZero() {
			continue
		}
		structField := types.Field(i)
		column := getColumn(structField)
		columns[length] = column
		values[length] = field.Interface()
		length += 1
	}

	return columns, length, values
}
