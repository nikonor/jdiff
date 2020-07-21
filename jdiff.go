package jdiff

import (
	"encoding/json"
	// "github.com/tidwall/sjson"
)

type DiffType struct {
	cmd   string
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
	println("\n\n" + path + "\n\told=" + string(old) + "\n\tnew=" + string(new))
	var (
		ret            []DiffType
		oldMap, newMap map[string]json.RawMessage
		oldErr, newErr *json.UnmarshalTypeError
		ok             bool
	)
	// sj   son.SetBytes()

	err := json.Unmarshal(old, &oldMap)
	if err != nil {
		if oldErr, ok = err.(*json.UnmarshalTypeError); !ok {
			return nil, err
		}
	}

	err = json.Unmarshal(new, &newMap)
	if err != nil {
		if newErr, ok = err.(*json.UnmarshalTypeError); !ok {
			return nil, err
		}
	}

	//  если было что-то, а стало другого типа, то записываем заменить
	//  если было значений тоже проверяем тип
	//  если чего-то не было, а мы его нашли, то записываем добавить
	//  если что-то было, а его не стало, то записываем удалить

	// fmt.Printf("%s =>\n\told=%s\n\tnew=%s\n", path, string(old), string(new))
	// fmt.Printf("\n\toldErr=%#v\n\tnew=%#v\n\n", oldErr, newErr)

	switch {

	case oldErr != nil && newErr != nil:
		// у два не объекта
		if oldErr.Value != newErr.Value {
			ret = append(ret, DiffType{
				cmd:   "set",
				path:  path,
				value: new,
			})
		}

	case oldErr != nil && newErr == nil:
		// было значение,а стало нет
		println("тут:" + string(old) + "=>" + string(new))
		ret = append(ret, DiffType{
			cmd:   "set",
			path:  path,
			value: new,
		})

	case oldErr == nil && newErr != nil:

	case oldErr == nil && newErr == nil:
		// у нас два объекта

		// проверяем объекты, что есть в old
		for k, oldV := range oldMap {
			newV, ok := newMap[k]
			oldVj, _ := oldV.MarshalJSON()
			newVj, _ := newV.MarshalJSON()

			if !ok {
				// объект удален
				// println("в new нет " + k + "=>" + string(oldVj))
				ret = append(ret, DiffType{
					cmd:   "delete",
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
					cmd:   "add",
					path:  appendPath(path, k),
					value: newV,
				})
			}
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
