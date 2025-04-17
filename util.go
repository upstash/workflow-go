package workflow

import (
	"encoding/json"
	"reflect"
)

func Retry(r int) *int {
	return &r
}

func Parallelism(p int) *int {
	return &p
}

func Rate(r int) *int {
	return &r
}

func serializeToStr(data any) (string, bool, error) {
	t := reflect.ValueOf(data).Kind()
	if t == reflect.String {
		return data.(string), false, nil
	} else if t == reflect.Invalid {
		return "", false, nil
	} else {
		data, err := json.Marshal(data)
		if err != nil {
			return "", false, err
		}
		return string(data), true, nil
	}
}
