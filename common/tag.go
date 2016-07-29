package common

//ReadOSEnv 读取struct内tag=osenv的标签内容，转为读取环境变量
import (
	"os"
	"reflect"
)

func ReadOSEnv(obj interface{}) error {
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		osenvKey := field.Tag.Get("osenv")

		if osenvKey != "" {
			osenvValue := os.Getenv(osenvKey)

			if osenvValue != "" {
				val.Field(i).SetString(osenvValue)
			}

		}

	}

	return nil
}
