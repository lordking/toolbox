package mongo

import (
	"encoding/json"
	"reflect"

	"github.com/lordking/toolbox/common"
)

//UpdateJsonWith 制作update内容
func UpdateJsonWith(obj interface{}) (map[string]interface{}, error) {

	var updateJSON map[string]interface{}

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, common.NewError(common.ErrCodeInternal, err.Error())
	}

	json.Unmarshal(data, &updateJSON)

	for key, value := range updateJSON {
		if value != nil {

			typeName := reflect.TypeOf(value).Name()
			switch typeName {
			case "string":
				if value.(string) == "" {
					delete(updateJSON, key)
				}
				break
			}

		} else {
			delete(updateJSON, key)
		}
	}

	return updateJSON, nil
}
