package redis

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"

	"github.com/lordking/toolbox/common"
)

type (
	Config struct {
		Host        string `json:"host" env:"REDIS_HOST"`
		Port        string `json:"port" env:"REDIS_PORT"`
		Password    string `json:"password" env:"REDIS_PASSWORD"`
		MaxIdle     int    `json:"maxIdle" env:"REDIS_MAX_IDLE"`
		IdleTimeout int64  `json:"idleTimeout" env:"REDIS_IDLE_TIMEOUT"`
	}

	Connection struct {
		redis.Conn
	}

	Redis struct {
		Config     *Config
		Connection redis.Conn
		Pool       *redis.Pool
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

	if m.Config.MaxIdle == 0 {
		return common.NewError(common.ErrCodeInternal, "Not found `maxIdle` in config file and `REDIS_MAX_IDLE` in env")
	}

	if m.Config.IdleTimeout == 0 {
		return common.NewError(common.ErrCodeInternal, "Not found `idleTimeout` in config file and `REDIS_IDLE_TIMEOUT` in env")
	}

	return nil
}

func (m *Redis) Connect() error {

	address := m.Config.Host + m.Config.Port

	c, err := redis.Dial("tcp", address)
	if err != nil {
		return common.NewErrorWithOther(common.ErrCodeInternal, err)
	}

	if m.Config.Password != "" {
		_, err = c.Do("AUTH", m.Config.Password)

		if err != nil {
			c.Close()
			return common.NewErrorWithOther(common.ErrCodeInternal, err)
		}
	}

	m.Connection = c

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

func (m *Redis) SetObject(key string, value interface{}) error {

	data, _ := json.Marshal(value)
	str := string(data)

	if err := m.Connection.Send("SET", key, str); err != nil {
		return common.NewErrorWithOther(common.ErrCodeInternal, err)
	}

	if err := m.Connection.Flush(); err != nil {
		return common.NewErrorWithOther(common.ErrCodeInternal, err)
	}

	if _, err := m.Connection.Receive(); err != nil {
		return common.NewErrorWithOther(common.ErrCodeInternal, err)
	}

	return nil
}

func (m *Redis) GetObject(obj interface{}, key string) error {

	if err := m.Connection.Send("GET", key); err != nil {
		return common.NewErrorWithOther(common.ErrCodeInternal, err)
	}

	if err := m.Connection.Flush(); err != nil {
		return common.NewErrorWithOther(common.ErrCodeInternal, err)
	}

	value, err := m.Connection.Receive()
	if err != nil {
		return common.NewErrorWithOther(common.ErrCodeInternal, err)
	}
	str := value.(string)
	data := []byte(str)
	common.ReadJSON(obj, data)

	return nil
}

func New() *Redis {
	return &Redis{}
}
