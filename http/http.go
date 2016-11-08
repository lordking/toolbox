package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lordking/toolbox/common"
	"github.com/lordking/toolbox/log"
)

type (

	// Config http配置
	Config struct {
		Port    string `json:"port"`
		SSLPort string `json:"ssl_port"`
		SSLCert string `json:"ssl_cert"`
		SSLKey  string `json:"ssl_key"`
	}

	//Server http服务对象
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

//RunServOnSSL 带SSL运行http服务
func (h *Server) RunServOnSSL() {

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

//RunServ 运行http服务
func (h *Server) RunServ() {

	//设置WEB中间件
	h.Router.Use(gin.Recovery())
	h.Router.Use(gin.Logger())

	log.Info("HTTP  on %s", h.Config.Port)

	if err := http.ListenAndServe(h.Config.Port, h.Router); err != nil {
		log.Fatal("http serve failure: %s", err.Error())
	}

}

//NewServer 新建HTTP Server
func NewServer(config *Config) *Server {
	return &Server{
		Config: config,
		Router: gin.Default(),
	}
}

//CreateServer 创建http服务实例
func CreateServer(configPath string) *ClassicServer {

	data, err := common.ReadFileData(configPath)
	defer common.CheckFatal(err)

	config := &Config{}
	err = common.ReadJSON(config, data)
	defer common.CheckFatal(err)

	httpServer := NewServer(config)

	return &ClassicServer{httpServer.Router, httpServer}
}

//BasicAuth 提供http认证接口
func BasicAuth(authfn func(*gin.Context, string) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request

		authorization := req.Header.Get("Authorization")

		if authorization == "" {
			JSONResponse(c, http.StatusUnauthorized, "Not found authorization")
			c.Abort()
			return
		}

		if err := authfn(c, authorization); err != nil {
			JSONResponse(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}
