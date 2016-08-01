package common

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

//Header 定义一个header
type Header struct {
	Key   string
	Value string
}

//RequestJSON 以json方式向服务器发起请求
func RequestJSON(method, url string, data []byte, headers ...interface{}) ([]byte, error) {

	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, NewErrorWithOther(400, err)
	}
	req.Header.Set("Content-Type", "application/json")

	count := len(headers)
	for i := 0; i < count; i++ {
		header := headers[i].(Header)
		req.Header.Add(header.Key, header.Value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewErrorWithOther(400, err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, NewErrorWithOther(400, err)
	}

	if resp.StatusCode == 200 {
		return body, nil
	}

	return nil, NewError(resp.StatusCode, string(body))
}
