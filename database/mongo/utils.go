package mongo

import (
	"encoding/json"

	"github.com/lordking/toolbox/common"
)

func UpdateJsonWith(obj interface{}) (map[string]interface{}, error) {

	var updateJson map[string]interface{}

	data, err := json.Marshal(obj)

	if err != nil {
		return nil, common.NewErrorWithOther(common.ErrCodeInternal, err)
	}
	err = json.Unmarshal(data, &updateJson)
	if err != nil {
		return nil, common.NewErrorWithOther(common.ErrCodeInternal, err)
	}

	for key, value := range updateJson {
		if value == "" || value == nil {
			delete(updateJson, key)
		}
	}

	return updateJson, err
}
