package common

import (
	"fmt"
	"log"
)

const (
	//ErrCodedParams 请求参数错误
	ErrCodedParams = 400
	//ErrCodeNotFound 没有发现
	ErrCodeNotFound = 404
	//ErrCodeInternal 内部错误
	ErrCodeInternal = 500
)

//Error 自定义错误
type Error struct {
	Code    int    `json:"status"`
	Message string `json:"error"`
}

//Error 继承错误输出接口
func (e *Error) Error() string {
	return fmt.Sprintf("\n code: %d \n error: %s", e.Code, e.Message)
}

//NewError 生成新的自定义错误对象
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

//CheckFatal 打印失败类型错误，并因此停止程序运行
func CheckFatal(err error) {
	if err != nil {
		log.Fatal("fatal:", err.Error())
	}
}

//CheckError 打印普通错误
func CheckError(err error) {
	if err != nil {
		log.Print(err.Error())
	}
}
