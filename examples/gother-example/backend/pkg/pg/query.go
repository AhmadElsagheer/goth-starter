package pg

import (
	"reflect"

	"github.com/bokwoon95/sq"
)

func AllFields(table sq.Table) []sq.Field {
	var fields []sq.Field
	t := reflect.TypeOf(table)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() != reflect.Struct {
			continue
		}
		if !field.Type.Implements(reflect.TypeOf((*sq.Field)(nil)).Elem()) {
			continue
		}
		tblField := reflect.ValueOf(table).Field(i).Interface().(sq.Field)
		fields = append(fields, tblField)
	}
	return fields
}
