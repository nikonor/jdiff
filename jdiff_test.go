package jsonconfigdiff

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
					Cmd:   SetAction,
					Path:  "two",
					Value: []byte("22"),
				},
			},
		},
		{
			name: "Cмена типа значения #2",
			old:  []byte(`{"one":1, "two":"TWO","three":{"four":"FOUR"}}`),
			new:  []byte(`{"one":1, "two":"TWO","three":{"four":44}}`),
			want: []DiffType{
				{
					Cmd:   SetAction,
					Path:  "three.four",
					Value: []byte("44"),
				},
			},
		},
		{
			name: "Было значение, а стал объект",
			old:  []byte(`{"one":1, "two":"TWO"}`),
			new:  []byte(`{"one":1, "two": {"four":"FOUR"}}`),
			want: []DiffType{
				{
					Cmd:   SetAction,
					Path:  "two",
					Value: []byte(`{"four":"FOUR"}`),
				},
			},
		},
		{
			name: "Был объект, а стало значение",
			old:  []byte(`{"one":1, "two": {"four":"FOUR"}}`),
			new:  []byte(`{"one":1, "two":"TWO"}`),
			want: []DiffType{
				{
					Cmd:   SetAction,
					Path:  "two",
					Value: []byte(`"TWO"`),
				},
			},
		},
		{
			name: "Добавили параметр #1",
			old:  []byte(`{"one":1, "two":"TWO","three":{"four":"FOUR"}}`),
			new:  []byte(`{"one":1, "two":"TWO","three":{"four":"FOUR"},"four":44}`),
			want: []DiffType{
				{
					Cmd:   SetAction,
					Path:  "four",
					Value: []byte("44"),
				},
			},
		},
		{
			name: "Добавили параметр #2",
			old:  []byte(`{"one":1, "two":"TWO","three":{"four":"FOUR"}}`),
			new:  []byte(`{"one":1, "two":"TWO","three":{"four":"FOUR","five":false}}`),
			want: []DiffType{
				{
					Cmd:   SetAction,
					Path:  "three.five",
					Value: []byte("false"),
				},
			},
		},
		{
			name: "Удалили параметр",
			old:  []byte(`{"one":1, "two":"TWO","three":true}`),
			new:  []byte(`{"one":1, "two":"TWO"}`),
			want: []DiffType{{
				Cmd:   DelectAction,
				Path:  "three",
				Value: nil,
			}},
		},
		{
			name: "Удалили целый объект",
			old:  []byte(`{"one":1, "two":"TWO","three":{"four":"FOUR"}}`),
			new:  []byte(`{"one":1, "two":"TWO"}`),
			want: []DiffType{{
				Cmd:   DelectAction,
				Path:  "three",
				Value: nil,
			}},
		},
		{
			name: "Удалили целый массив",
			old:  []byte(`{"one":1, "two":"TWO","three":[1,2,3,128]}`),
			new:  []byte(`{"one":1, "two":"TWO"}`),
			want: []DiffType{{
				Cmd:   DelectAction,
				Path:  "three",
				Value: nil,
			}},
		},
		{
			name: "Добавили массив",
			old:  []byte(`{"one":1, "two":"TWO"}`),
			new:  []byte(`{"one":1, "two":"TWO","three":[{"key":1,"val":"1"},{"key":2, "val":"2"}]}`),
			want: []DiffType{{
				Cmd:   SetAction,
				Path:  "three",
				Value: []byte(`[{"key":1,"val":"1"},{"key":2, "val":"2"}]`),
			}},
		},
		{
			name: "Добавили в массив элемент",
			old:  []byte(`{"one":1, "two":"TWO","three":[{"key":1,"val":"1"},{"key":2, "val":"2"}]}`),
			new:  []byte(`{"one":1, "two":"TWO","three":[{"key":1,"val":"1"},{"key":2, "val":"2"},{"key":3,"val":"3"}]}`),
			want: []DiffType{{
				Cmd:   SetAction,
				Path:  "three",
				Value: []byte(`[{"key":1,"val":"1"},{"key":2, "val":"2"},{"key":3,"val":"3"}]`),
			}},
		},
		{
			name: "Изменили элемент массива",
			old:  []byte(`{"one":1, "two":"TWO","three":[{"key":1,"val":"1"},{"key":2, "val":"2"},{"key":3,"val":"3333"}]}`),
			new:  []byte(`{"one":1, "two":"TWO","three":[{"key":1,"val":"1"},{"key":2, "val":"2"},{"key":3,"val":"3"}]}`),
			want: nil,
		},
	}

	for _, c := range cases {
		// println("begin::" + c.name)
		// println("\told=" + string(c.old))
		// println("\tnew=" + string(c.new))
		// println("\t===")
		t.Run(c.name, func(t *testing.T) {
			got, err := JDiff(c.old, c.new)
			if err != nil {
				t.Error(err.Error())
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Error("Error on " + c.name + ":\nold=" + string(c.old) +
					"\nnew=" + string(c.new) +
					"\nwant=" + fmt.Sprintf("%#v", c.want) +
					"\ngot =" + fmt.Sprintf("%#v", got))
			}
		})
		// println("end::" + c.name + "\n")
	}
}
