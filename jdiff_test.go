package jdiff

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	cases := []struct {
		name string
		old  []byte
		new  []byte
		want []DiffType
	}{
		{
			name: "Нет изменений",
			old:  []byte(`{"one":1, "two":"TWO"}`),
			new:  []byte(`{"one":1, "two":"TWO"}`),
			want: nil,
		},
		{
			name: "Cмена типа значения #1",
			old:  []byte(`{"one":1, "two":"TWO","three":{"four":"FOUR"}}`),
			new:  []byte(`{"one":1, "two":22,"three":{"four":"FOUR"}}`),
			want: []DiffType{
				{
					cmd:   "set",
					path:  "two",
					value: []byte("22"),
				},
			},
		},
		{
			name: "Cмена типа значения #2",
			old:  []byte(`{"one":1, "two":"TWO","three":{"four":"FOUR"}}`),
			new:  []byte(`{"one":1, "two":"TWO","three":{"four":44}}`),
			want: []DiffType{
				{
					cmd:   "set",
					path:  "three.four",
					value: []byte("44"),
				},
			},
		},
		{
			name: "Было значение, а стал объект",
			old:  []byte(`{"one":1, "two":"TWO"}`),
			new:  []byte(`{"one":1, "two": {"four":"FOUR"}}`),
			want: []DiffType{
				{
					cmd:   "set",
					path:  "two",
					value: []byte(`{"four":"FOUR"}`),
				},
			},
		},
		{
			name: "Был объект, а стало значение",
			old:  []byte(`{"one":1, "two": {"four":"FOUR"}}`),
			new:  []byte(`{"one":1, "two":"TWO"}`),
			want: []DiffType{
				{
					cmd:   "set",
					path:  "two",
					value: []byte(`"TWO"`),
				},
			},
		},
		{
			name: "Добавили параметр #1",
			old:  []byte(`{"one":1, "two":"TWO","three":{"four":"FOUR"}}`),
			new:  []byte(`{"one":1, "two":"TWO","three":{"four":"FOUR"},"four":44}`),
			want: []DiffType{
				{
					cmd:   "add",
					path:  "four",
					value: []byte("44"),
				},
			},
		},
		{
			name: "Добавили параметр #2",
			old:  []byte(`{"one":1, "two":"TWO","three":{"four":"FOUR"}}`),
			new:  []byte(`{"one":1, "two":"TWO","three":{"four":"FOUR","five":false}}`),
			want: []DiffType{
				{
					cmd:   "add",
					path:  "three.five",
					value: []byte("false"),
				},
			},
		},
		{
			name: "Удалили параметр",
			old:  []byte(`{"one":1, "two":"TWO","three":true}`),
			new:  []byte(`{"one":1, "two":"TWO"}`),
			want: []DiffType{{
				cmd:   "delete",
				path:  "three",
				value: nil,
			}},
		},
		{
			name: "Удалили целый объект",
			old:  []byte(`{"one":1, "two":"TWO","three":{"four":"FOUR"}}`),
			new:  []byte(`{"one":1, "two":"TWO"}`),
			want: []DiffType{{
				cmd:   "delete",
				path:  "three",
				value: nil,
			}},
		},
	}

	for _, c := range cases {
		println("begin::" + c.name)
		println("\told=" + string(c.old))
		println("\tnew=" + string(c.new))
		println("\t===")
		got, err := JDiff(c.old, c.new)
		println("\t===")
		if err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Error("Error on " + c.name + ":\nold=" + string(c.old) +
				"\nnew=" + string(c.new) +
				"\nwant=" + fmt.Sprintf("%#v", c.want) +
				"\ngot =" + fmt.Sprintf("%#v", got))
		}
		println("end::" + c.name + "\n")
	}
}
