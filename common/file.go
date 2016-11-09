package common

import (
	"io/ioutil"
)

func ReadFileData(path string) ([]byte, error) {
	var (
		data []byte
		err  error
	)

	if data, err = ioutil.ReadFile(path); err != nil {
		return nil, NewError(ErrCodeInternal, err.Error())
	}

	return data, nil
}
