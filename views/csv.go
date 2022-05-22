package views

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
)

func WriteCSVTo(o, e io.Writer, v interface{}) {
	top := reflect.TypeOf(v)
	if k := top.Kind(); k != reflect.Slice {
		panic(`called WriteCSVTo with a non-slice parameter`)
	}
	vv := reflect.ValueOf(v)

	t := reflect.TypeOf(v).Elem()
	if k := t.Kind(); k != reflect.Struct {
		panic(`called WriteCSVTo with a non-slice-of-struct parameter`)
	}

	records := [][]string{}
	fields := []reflect.StructField{}
	var field reflect.StructField
	fieldLabels := []string{}

	// reflect all the fields and build the header row
	for i := 0; i < t.NumField(); i++ {
		field = t.Field(i)
		fields = append(fields, field)
		fieldLabels = append(fieldLabels, field.Tag.Get(`csv`))
	}
	records = append(records, fieldLabels)
	for i := 0; i < vv.Len(); i++ {
		// retrieve the record as a Value
		r := vv.Index(i)

		row := []string{}
		for _, f := range fields {
			row = append(row, fmt.Sprintf("%v", r.FieldByName(f.Name)))
		}
		records = append(records, row)
	}

	// write it all
	w := csv.NewWriter(o)
	err := w.WriteAll(records)
	if err != nil {
		fmt.Fprintln(e, err)
	}
	if err = w.Error(); err != nil {
		fmt.Fprintln(e, err)
	}
}
