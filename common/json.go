package common

import (
	"bytes"
	"encoding/json"
)

//ReadJSON 将json格式的byte内容转换成对象
func ReadJSON(obj interface{}, data []byte) error {

	if data == nil {
		return nil
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		return NewError(ErrCodeInternal, err.Error())
	}

	if obj == nil {
		return NewError(ErrCodeInternal, "config empty")
	}

	return nil
}

//JSONObjectConvert 将一个对象转换成序列化后的另一个对象
func JSONObjectConvert(obj1 interface{}, obj2 interface{}) error {

	data, err := json.Marshal(obj1)
	if err != nil {
		return NewError(ErrCodeInternal, err.Error())
	}

	if err := json.Unmarshal(data, obj2); err != nil {
		return NewError(ErrCodeInternal, err.Error())
	}

	return nil
}

//PrettyJSON 打印缩进后的json内容
func PrettyJSON(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "    ")
	return out.Bytes(), err
}

func PrettyObject(obj interface{}) []byte {

	data, _ := json.Marshal(obj)
	result, _ := PrettyJSON(data)

	return result
}
