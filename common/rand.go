package common

import (
	"math/rand"
	"time"
)

//RandInt64 生成随机数字
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}

	rand.Seed(time.Now().Unix())
	return rand.Int63n(max-min) + min
}

//RandString 生成随机字符串
func RandString(length int) []byte {
	option := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	optionLen := len(option)

	ss := make([]byte, length)
	rand.Seed(time.Now().Unix())
	for i := 0; i < length; i++ {
		n := rand.Intn(optionLen)
		c := option[n]
		ss[i] = c
	}

	return ss
}
