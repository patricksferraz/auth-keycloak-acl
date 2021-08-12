package utils

import (
	"os"
	"reflect"
)

func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func StructToAttr(item interface{}) *map[string][]string {
	res := make(map[string][]string)
	if item == nil {
		return &res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("attr")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v, ok := field.(string); ok {
				res[tag] = []string{v}
			}
		}
	}
	return &res
}
