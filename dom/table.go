package dom

import (
	"fmt"
	"reflect"
)

func NewTable[T any](s []T, a ...HtmlAttr) HtmlElement {
	var fields []int

	e := NewElement("table", "", a...)

	r := NewElement("tr", "")
	t := reflect.TypeOf(s[0])
	for i := 0; i < t.NumField(); i++ {
		if f := t.Field(i).Tag.Get("html"); f != "" {
			fields = append(fields, i)
			c := NewElement("th", f)
			r.AppendChild(c)
		}
	}
	e.AppendChild(r)

	for i := 0; i < len(s); i++ {
		r := NewElement("tr", "")
		v := reflect.ValueOf(s[i])
		for _, f := range fields {
			vv := v.FieldByIndex([]int{f})
			text := fmt.Sprint(vv)
			c := NewElement("td", text)
			r.AppendChild(c)
		}
		e.AppendChild(r)
	}

	return e
}
