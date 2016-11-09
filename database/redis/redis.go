package redis

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
	"github.com/lordking/toolbox/common"
	"github.com/lordking/toolbox/database"
	"github.com/lordking/toolbox/log"
)

var _ database.Database = (*Redis)(nil)

type (
	Config struct {
		Host        string `json:"host" env:"REDIS_HOST"`
		Port        string `json:"port" env:"REDIS_PORT"`
		Password    string `json:"password" env:"REDIS_PASSWORD"`
		MaxIdle     int    `json:"maxIdle" env:"REDIS_MAX_IDLE"`
		IdleTimeout int64  `json:"idleTimeout" env:"REDIS_IDLE_TIMEOUT"`
	}

	conn struct {
		redis.Conn
	}

	Redis struct {
		ReceiveQueue chan []byte
		config       *Config
		conn         redis.Conn
	}
)

func (m *Redis) NewConfig() interface{} {
	m.config = &Config{}
	return m.config
}

func (m *Redis) ValidateBefore() error {

	if m.config.Host == "" {
		return common.NewError(common.ErrCodeInternal, "Not found `host` in config file and `REDIS_HOST` in env")
	}

	if m.config.Port == "" {
		return common.NewError(common.ErrCodeInternal, "Not found `port` in config file and `REDIS_PORT` in env")
	}

	if m.config.MaxIdle == 0 {
		return common.NewError(common.ErrCodeInternal, "Not found `maxIdle` in config file and `REDIS_MAX_IDLE` in env")
	}

	if m.config.IdleTimeout == 0 {
		return common.NewError(common.ErrCodeInternal, "Not found `idleTimeout` in config file and `REDIS_IDLE_TIMEOUT` in env")
	}

	return nil
}

func (m *Redis) Connect() error {

	address := m.config.Host + m.config.Port

	c, err := redis.Dial("tcp", address)
	if err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	if m.config.Password != "" {
		_, err = c.Do("AUTH", m.config.Password)

		if err != nil {
			c.Close()
			return common.NewError(common.ErrCodeInternal, err.Error())
		}
	}

	m.conn = c

	return nil
}

func (m *Redis) GetConnection() interface{} {
	return m.conn
}

func (m *Redis) Close() error {
	if err := m.conn.Close(); err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	return nil
}

func (m *Redis) SetObject(key string, value interface{}, expire int) error {

	json, _ := json.Marshal(value)
	str := string(json)

	if err := m.conn.Send("SET", key, str); err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	if expire > 0 {
		if err := m.conn.Send("EXPIRE", key, expire); err != nil {
			return common.NewError(common.ErrCodeInternal, err.Error())
		}
	}

	if err := m.conn.Flush(); err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	if _, err := m.conn.Receive(); err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	return nil
}

func (m *Redis) GetObject(obj interface{}, key string) error {

	if err := m.conn.Send("GET", key); err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	if err := m.conn.Flush(); err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	value, err := m.conn.Receive()
	if err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	common.ReadJSON(obj, value.([]byte))

	return nil
}

func (m *Redis) DeleteObject(key string) error {

	if err := m.conn.Send("DEL", key); err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	if err := m.conn.Flush(); err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	_, err := m.conn.Receive()
	if err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	return nil
}

func (m *Redis) PublishObject(channel string, value interface{}) error {

	json, _ := json.Marshal(value)
	str := string(json)

	if err := m.conn.Send("PUBLISH", channel, str); err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	if err := m.conn.Flush(); err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	_, err := m.conn.Receive()
	if err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	return nil
}

func (m *Redis) Subscribe(channel string) (redis.PubSubConn, error) {

	psc := redis.PubSubConn{Conn: m.conn}
	err := psc.Subscribe(channel)

	return psc, err
}

func (m *Redis) Receive(psc redis.PubSubConn) {

	m.ReceiveQueue = make(chan []byte)
	go func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				m.ReceiveQueue <- v.Data
			case redis.Subscription:
				log.Debugf("subscribe: %s, count: %d", v.Channel, v.Count)
			case error:
				log.Error("Error:", v.Error())
			}
		}
	}()

}

func New() *Redis {
	return &Redis{}
}
