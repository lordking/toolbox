package goutils

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func init() {
	LogDebug("Init http.go")
}

// HTTPConfig http配置
type HTTPConfig struct {
	Port    string `json:"port"`
	SSLPort string `json:"ssl_port"`
	SSLCert string `json:"ssl_cert"`
	SSLKey  string `json:"ssl_key"`
}

//HTTPServer http服务对象
type HTTPServer struct {
	Config *HTTPConfig
	Router *gin.Engine
}

//RunServ 运行http服务
func (h *HTTPServer) RunServ() {

	//设置WEB中间件
	h.Router.Use(gin.Recovery())
	h.Router.Use(gin.Logger())

	go func() {
		LogInfo("HTTP  on %s", h.Config.Port)

		if err := http.ListenAndServe(h.Config.Port, h.Router); err != nil {
			LogFatal("start failure:%q", err)
		}
	}()

	LogInfo("HTTPS on %s", h.Config.SSLPort)

	if err := http.ListenAndServeTLS(h.Config.SSLPort, h.Config.SSLCert, h.Config.SSLKey, h.Router); err != nil {
		LogFatal("start failure:%q", err)
	}

}

//NewHTTPServer 新建
func NewHTTPServer(config *HTTPConfig) *HTTPServer {
	return &HTTPServer{
		Config: config,
		Router: gin.Default(),
	}
}

//ClassicHTTPServer 实例化的http服务对象
type ClassicHTTPServer struct {
	*gin.Engine
	*HTTPServer
}

//LoadHTTPConfig 导入http配置文件
func LoadHTTPConfig(configPath string) *HTTPConfig {
	//导入配置文件
	data, err := GetFileData(configPath)
	CheckFatal(err)

	config := &HTTPConfig{}
	err = ReadJSON(config, data)
	CheckFatal(err)

	return config
}

//CreateHTTPServer 创建http服务实例
func CreateHTTPServer(configPath string) *ClassicHTTPServer {

	config := LoadHTTPConfig(configPath)
	httpServer := NewHTTPServer(config)

	return &ClassicHTTPServer{httpServer.Router, httpServer}
}

//BasicAuth 提供http认证接口
func BasicAuth(authfn func(string, string) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request

		auth := req.Header.Get("Authorization")
		if len(auth) < 6 || auth[:6] != "Basic " {
			JSONResponse(c, http.StatusUnauthorized, "not found token")
			c.Abort()
			return
		}
		b, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			JSONResponse(c, http.StatusUnauthorized, "wrong token format")
			c.Abort()
			return
		}
		tokens := strings.SplitN(string(b), ":", 2)
		if len(tokens) != 2 {
			JSONResponse(c, http.StatusUnauthorized, "wrong token formart")
			c.Abort()
			return
		}

		err = authfn(tokens[0], tokens[1])
		if err != nil {
			JSONResponse(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}

		c.Next()
	}
}
