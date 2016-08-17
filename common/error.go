package common

import (
	"fmt"

	"github.com/lordking/toolbox/log"
)

const (
	ErrCodedParams  = 403 //请求参数错误
	ErrCodeNotFound = 404 //没有发现
	ErrCodeInternal = 500 //内部错误
)

//自定义错误
type Error struct {
	Code    int    `json:"status"`
	Message string `json:"error"`
}

//继承错误输出接口
func (e *Error) Error() string {
	return fmt.Sprintf("\n code: %d \n error: %s", e.Code, e.Message)
}

//生成新的自定义错误对象
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

//使用其他的错误对象生成新的自定义错误对象
// func NewErrorWithOther(code int, err error) *Error {
// 	return &Error{
// 		Code:    code,
// 		Message: err.Error(),
// 	}
// }

//CheckFatal 打印失败类型错误，并因此停止程序运行
func CheckFatal(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

//CheckError 打印普通错误
func CheckError(err error) {
	if err != nil {
		log.Error(err.Error())
	}
}
