package common

// GetFileData 读取文件内容
import "io/ioutil"

func GetFileData(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, NewErrorWithOther(ErrCodeInternal, err)
	}

	return data, nil
}
