package goutils

import "github.com/gin-gonic/gin"

//JSONResponse 发送定制json内容
func JSONResponse(c *gin.Context, status int, obj interface{}) {

	if status == 200 {
		c.JSON(200, gin.H{"status": 200, "result": obj})
	} else {
		c.JSON(status, gin.H{"status": status, "error": obj})
	}
}
