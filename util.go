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
		column = toUnderLineCase(field.Name)
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
	columns := make([]string, 0)
	values := make([]any, 0)
	types := anyutil.Type(v)

	for i := 0; i < numField; i++ {
		field := value.Field(i)
		if field.IsZero() {
			continue
		}
		structField := types.Field(i)
		column := getColumn(structField)
		columns = append(columns, column)
		values = append(values, field.Interface())
		length += 1
	}

	return columns, length, values
}

func endWith(src, s string) bool {
	srcLen := len(src)
	sLen := len(s)

	if srcLen == 0 || sLen == 0 {
		return false
	}
	if srcLen < sLen {
		return false
	}

	for i := 1; i <= sLen; i++ {
		if src[srcLen-i] != s[sLen-i] {
			return false
		}
	}
	return true
}

func toUnderLineCase(s string) string {
	if str.IsEmpty(s) {
		return ""
	}
	sb := strings.Builder{}
	sb.Grow(len(s))
	for i, item := range s {
		if i == 0 && 65 <= item && item <= 90 {
			sb.WriteByte(byte(item + 32))
			continue
		}
		if 65 <= item && item <= 90 {
			sb.WriteByte(95)
			sb.WriteByte(byte(item + 32))
			continue
		}
		sb.WriteByte(byte(item))
	}
	return sb.String()
}
