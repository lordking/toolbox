package common

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
)

//EncryptMAC 用hmac加密字符串。每次加密后的密文长度都是32位。可用于字符串完整性校验
func EncryptMAC(s, key string) []byte {

	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(s))

	return mac.Sum(nil)
}

//EncryptMD5 用md5加密字符串。每次加密后的密文长度都是16位。可用于字符串完整性校验
func EncryptMD5(s string) []byte {

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(s))

	return md5Ctx.Sum(nil)
}
