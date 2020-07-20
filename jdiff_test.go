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
			name: "Добавили параметр",
			old:  []byte(`{"one":1, "two":"TWO"}`),
			new:  []byte(`{"one":1, "two":"TWO","three":true}`),
			want: []DiffType{{
				isAdd: true,
				path:  "three",
				value: []byte("true"),
			}},
		},
		{
			name: "Удалили параметр",
			old:  []byte(`{"one":1, "two":"TWO","three":true}`),
			new:  []byte(`{"one":1, "two":"TWO"}`),
			want: []DiffType{{
				isAdd: false,
				path:  "three",
				value: nil,
			}},
		},
	}

	for _, c := range cases {
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
	}
}
