package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/lordking/toolbox/common"
)

//Header 定义一个header
type Header struct {
	Key   string
	Value string
}

//JSONResponse 发送定制json内容
func JSONResponse(c *gin.Context, status int, obj interface{}) {

	if status == 200 {

		switch obj.(type) {
		case string:
			c.JSON(status, gin.H{"result": obj})
			break
		default:
			c.JSON(200, obj)
		}

	} else {
		c.JSON(status, gin.H{"error": obj})
	}
}

//RequestJSON 以json字符串方式向服务器发起请求
func RequestJSON(method, url string, data []byte, headers ...interface{}) ([]byte, error) {

	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, common.NewError(400, err.Error())
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
		return nil, common.NewError(400, err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, common.NewError(400, err.Error())
	}

	if resp.StatusCode == 200 {
		return body, nil
	}

	return nil, common.NewError(resp.StatusCode, string(body))
}

//以form方式向服务器发起请求
func PostForm(url string, form url.Values) ([]byte, error) {

	client := &http.Client{}
	resp, err := client.PostForm(url, form)
	if err != nil {
		msg := fmt.Sprintf("Request error: %s", err)
		return nil, common.NewError(400, msg)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("Request error: %s", err)
		return nil, common.NewError(400, msg)
	}

	if resp.StatusCode != 200 {
		return nil, common.NewError(resp.StatusCode, string(body))
	}

	return body, nil
}

//以get方式向服务器发起请求
func GetForm(url string) ([]byte, error) {

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		msg := fmt.Sprintf("Request error: %s", err)
		return nil, common.NewError(400, msg)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("Request error: %s", err)
		return nil, common.NewError(400, msg)
	}

	if resp.StatusCode != 200 {
		return nil, common.NewError(resp.StatusCode, string(body))
	}

	return body, nil
}
