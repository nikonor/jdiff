package jdiff

import (
	"encoding/json"
	"fmt"
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
		return nil, err
	}

	for _, r := range res {
		println("\t\t" + r.String())
	}
	return res, nil
}

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
	fmt.Printf("\n\toldErr=%#v\n\tnew=%#v\n\n", oldErr, newErr)

	switch {

	case oldErr != nil && newErr != nil && oldErr.Value == "array" && newErr.Value == "array":
		// у нас два массива. будет сравнивать
		// TODO:
		if changed, err := cmdArray(
			[]byte(`{"array":`+string(old)+`}`),
			[]byte(`{"array":`+string(new)+`}`),
		); err != nil {
			return nil, err
		} else {
			if changed {
				ret = append(ret, DiffType{
					cmd:   "set",
					path:  path,
					value: new,
				})
			}
		}
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
		// println("тут:" + string(old) + "=>" + string(new))
		ret = append(ret, DiffType{
			cmd:   "set",
			path:  path,
			value: new,
		})

	case oldErr == nil && newErr != nil:
		ret = append(ret, DiffType{
			cmd:   "set",
			path:  path,
			value: new,
		})

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
			// newVj, _ := newV.MarshalJSON()
			if _, ok := oldMap[k]; !ok {
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

func cmdArray(old, new []byte) (bool, error) {
	var (
		oldA, newA aType
	)

	if err := json.Unmarshal(old, &oldA); err != nil {
		return false, err
	}
	if err := json.Unmarshal(new, &newA); err != nil {
		return false, err
	}

	// fmt.Printf("old=%#v\nnew=%#v\n", oldA.Array, newA.Array)
	if len(oldA.Array) != len(newA.Array) {
		return true, nil
	}

	return false, nil
}

type aType struct {
	Array []interface{} `json:"array"`
}
