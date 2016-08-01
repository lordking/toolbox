package common

// GetFileData 读取文件内容
import (
	"io/ioutil"
)

func ReadFileData(path string) ([]byte, error) {
	var (
		data []byte
		err  error
	)

	if data, err = ioutil.ReadFile(path); err != nil {
		return nil, NewErrorWithOther(ErrCodeInternal, err)
	}

	return data, nil
}

func ReadConfig(config interface{}, path string) error {

	var (
		data []byte
		err  error
	)

	if path == "" {
		return NewError(ErrCodedParams, "Not found config file path")
	}

	if data, err = ReadFileData(path); err != nil {
		return err
	}

	if err = ReadJSON(config, data); err != nil {
		return err
	}

	if err = ReadEnv(config); err != nil {
		return err
	}

	return nil
}
