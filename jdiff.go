package jdiff

import "encoding/json"

func Diff(old, new []byte) {
	diff("", old, new)
}

func diff(path string, old, new []byte) error {
	println("path=" + path)
	var oldMap, newMap map[string]json.RawMessage

	if err := json.Unmarshal(old, &oldMap); err != nil {
		return err
	}

	if err := json.Unmarshal(new, &newMap); err != nil {
		return err
	}

	// проверяем объекты, что есть в old
	for k, oldV := range oldMap {
		newV, ok := newMap[k]
		oldVj, _ := oldV.MarshalJSON()
		newVj, _ := newV.MarshalJSON()

		if !ok {
			// объект удален
			println("в new нет " + k + "=>" + string(oldVj))
			continue
		}

		// объект найден - проверяем его
		diff(path+"."+k, oldVj, newVj)
	}

	// проверяем объекты, что могли добавиться
	for k, newV := range newMap {
		newVj, _ := newV.MarshalJSON()
		if _, ok := oldMap[k]; !ok {
			println("в new есть новый объект " + k + "=>" + string(newVj))
		}
	}

	return nil
}
