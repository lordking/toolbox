package common

//ReadOSEnv 读取struct内tag=osenv的标签内容，转为读取环境变量
import (
	"bytes"
	"encoding/gob"
	"os"
	"reflect"
)

//复制对象
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

//根据对象的标签定义的key，读取环境变量的值到对象
func ReadEnv(obj interface{}) error {

	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	} else {
		NewError(ErrCodeNotFound, "The object is null")
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		key := field.Tag.Get("env")
		value := os.Getenv(key)

		if key != "" && value != "" {
			val.Field(i).SetString(value)
		}

	}

	return nil
}
