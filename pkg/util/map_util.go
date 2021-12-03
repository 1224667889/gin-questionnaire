package util

import (
	"reflect"
)

func StructToMapDemo(obj interface{}) map[string]interface{}{
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		if obj2.Field(i).CanInterface() {
			data[obj1.Field(i).Tag.Get("json")] = obj2.Field(i).Interface()
		}
	}
	return data
}
