package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/lordking/toolbox/common"
)

type (
	Config struct {
		Host        string `json:"host" env:"REDIS_HOST"`
		Port        string `json:"port" env:"REDIS_PORT"`
		Password    string `json:"password" env:"REDIS_PASSWORD"`
		MaxIdle     string `json:"maxIdle" env:"REDIS_MAX_IDLE"`
		IdleTimeout string `json:"idleTimeout" env:"REDIS_IDLE_TIMEOUT"`
	}

	Redis struct {
		Config     *Config
		Connection *redis.Pool
	}
)

func (m *Redis) NewConfig() interface{} {
	m.Config = &Config{}
	return m.Config
}

func (m *Redis) ValidateBefore() error {

	if m.Config.Host == "" {
		return common.NewError(common.ErrCodeInternal, "Not found `host` in config file and `REDIS_HOST` in env")
	}

	if m.Config.Port == "" {
		return common.NewError(common.ErrCodeInternal, "Not found `port` in config file and `REDIS_PORT` in env")
	}

	if m.Config.Password == "" {
		return common.NewError(common.ErrCodeInternal, "Not found `password` in config file and `REDIS_PASSWORD` in env")
	}

	if m.Config.MaxIdle == "" {
		return common.NewError(common.ErrCodeInternal, "Not found `maxIdle` in config file and `REDIS_MAX_IDLE` in env")
	}

	if m.Config.IdleTimeout == "" {
		return common.NewError(common.ErrCodeInternal, "Not found `idleTimeout` in config file and `REDIS_IDLE_TIMEOUT` in env")
	}

	return nil
}

func (m *Redis) Connect() error {

	server := m.Config.Host + m.Config.Port
	m.Connection = newPool(server, m.Config.Password)

	return nil
}

func (m *Redis) GetConnection() interface{} {
	return m.Connection
}

func (m *Redis) Close() error {
	if err := m.Connection.Close(); err != nil {
		return common.NewErrorWithOther(common.ErrCodeInternal, err)
	}

	return nil
}

func New() *Redis {
	return &Redis{}
}

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
