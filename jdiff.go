package jdiff

import (
	"encoding/json"
	// "github.com/tidwall/sjson"
)

type DiffType struct {
	isAdd bool
	path  string
	value json.RawMessage
}

func (d DiffType) String() string {
	j, _ := d.value.MarshalJSON()
	return d.path + "=>" + string(j)
}

func JDiff(old, new []byte) ([]DiffType, error) {
	res, err := jdiff("", old, new)
	if err != nil {
		println("ERROR=" + err.Error())
	}

	for _, r := range res {
		println(r.String())
	}
	return res, nil
}

// TODO: надо возвращать что и куда надо добавить/удалить
//      1) формат???
func jdiff(path string, old, new []byte) ([]DiffType, error) {
	println("path=" + path)
	var ret []DiffType
	// sjson.SetBytes()

	// TODO: анализ того, что пришло. Нужно ли разбирать на словарь?

	var oldMap, newMap map[string]json.RawMessage

	if err := json.Unmarshal(old, &oldMap); err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(new, &newMap); err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, nil
		}
		return nil, err
	}

	// проверяем объекты, что есть в old
	for k, oldV := range oldMap {
		newV, ok := newMap[k]
		oldVj, _ := oldV.MarshalJSON()
		newVj, _ := newV.MarshalJSON()

		if !ok {
			// объект удален
			println("в new нет " + k + "=>" + string(oldVj))
			ret = append(ret, DiffType{
				isAdd: false,
				path:  appendPath(path, k),
				value: newV,
			})
			continue
		}

		// объект найден - проверяем его
		r, err := jdiff(appendPath(path, k), oldVj, newVj)
		if err != nil {
			return nil, err
		}
		ret = append(ret, r...)
	}

	// проверяем объекты, что могли добавиться
	for k, newV := range newMap {
		newVj, _ := newV.MarshalJSON()
		if _, ok := oldMap[k]; !ok {
			println("в new есть новый объект " + k + "=>" + string(newVj))
			ret = append(ret, DiffType{
				isAdd: true,
				path:  appendPath(path, k),
				value: newV,
			})
		}
	}

	return ret, nil
}

func appendPath(path, k string) string {
	ret := ""
	if len(path) > 0 {
		ret = path + "."
	}
	return ret + k
}
