package gorm_expand

import "strings"

func generateSelect[T any]() string {
	v := new(T)
	builder := strings.Builder{}
	builder.WriteString("select ")
	columns := Columns(v)
	length := len(columns)

	for i, column := range columns {
		builder.WriteString(column)

		if i == length-1 {
			builder.WriteString(" from ")
			builder.WriteString(TableName(v))
			break
		}

		builder.WriteString(",")
	}

	return builder.String()
}

func generateCreate(v any) (string, int, []any) {
	columns, length, values := ColumnsNotNil(v)

	tableName := TableName(v)

	builder := strings.Builder{}
	builder.WriteString("insert into ")
	builder.WriteString(tableName)
	builder.WriteString("(")

	for i, column := range columns {
		builder.WriteString(column)
		if i == length-1 {
			builder.WriteString(")")
			break
		}
		builder.WriteString(",")
	}

	builder.WriteString(" ")
	builder.WriteString("values(")
	for i := 0; i < length; i++ {
		builder.WriteString("?")
		if i == length-1 {
			builder.WriteString(")")
			break
		}
		builder.WriteString(",")
	}

	return builder.String(), length, values
}
