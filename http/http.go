package http

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lordking/toolbox/common"
	"github.com/lordking/toolbox/log"
)

type (

	// HTTPConfig http配置
	Config struct {
		Port    string `json:"port"`
		SSLPort string `json:"ssl_port"`
		SSLCert string `json:"ssl_cert"`
		SSLKey  string `json:"ssl_key"`
	}

	//HTTPServer http服务对象
	Server struct {
		Config *Config
		Router *gin.Engine
	}

	//ClassicHTTPServer 实例化的http服务对象
	ClassicServer struct {
		*gin.Engine
		*Server
	}
)

//RunServ 运行http服务
func (h *Server) RunServ() {

	//设置WEB中间件
	h.Router.Use(gin.Recovery())
	h.Router.Use(gin.Logger())

	go func() {
		log.Info("HTTP  on %s", h.Config.Port)

		if err := http.ListenAndServe(h.Config.Port, h.Router); err != nil {
			log.Fatal("http serve failure: %s", err.Error())
		}
	}()

	log.Info("HTTPS on %s", h.Config.SSLPort)

	if err := http.ListenAndServeTLS(h.Config.SSLPort, h.Config.SSLCert, h.Config.SSLKey, h.Router); err != nil {
		log.Fatal("https serve failure: %s", err.Error())
	}

}

//新建HTTP Server
func NewServer(config *Config) *Server {
	return &Server{
		Config: config,
		Router: gin.Default(),
	}
}

//CreateHTTPServer 创建http服务实例
func CreateServer(path string) *ClassicServer {

	data, err := common.ReadFileData(path)
	defer common.CheckFatal(err)

	config := &Config{}
	err = common.ReadJSON(config, data)
	defer common.CheckFatal(err)

	httpServer := NewServer(config)

	return &ClassicServer{httpServer.Router, httpServer}
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
