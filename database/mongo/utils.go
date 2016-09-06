package mongo

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/lordking/toolbox/common"
)

func UpdateJsonWith(obj interface{}) (map[string]interface{}, error) {

	var updateJson map[string]interface{}
	updateJson = make(map[string]interface{})

	data, err := json.Marshal(obj)

	if err != nil {
		return nil, common.NewError(common.ErrCodeInternal, err.Error())
	}
	err = json.Unmarshal(data, &updateJson)
	if err != nil {
		return nil, common.NewError(common.ErrCodeInternal, err.Error())
	}

	for key, value := range updateJson {
		if value != nil {

			typeName := reflect.TypeOf(value).Name()

			switch typeName {
			case "string":
				fmt.Printf("%s, %s, %s\n", key, typeName, value)
				if value.(string) == "" {
					delete(updateJson, key)
				}
				break
			}

		} else {
			delete(updateJson, key)
		}
	}

	return updateJson, err
}
