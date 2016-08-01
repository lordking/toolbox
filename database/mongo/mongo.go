package mongo

import (
	"time"

	"gopkg.in/mgo.v2"

	"github.com/lordking/toolbox/common"
)

const (
	connTimeout = time.Second * 5
)

type (
	Config struct {
		Url      string `json:"url" env:"MONGO_URL"`
		Database string `json:"database" env:"MONGO_DATABASE"`
	}

	Mongo struct {
		Config     *Config
		Connection *mgo.Session
	}
)

func (m *Mongo) NewConfig() interface{} {
	m.Config = &Config{}
	return m.Config
}

func (m *Mongo) ValidateBefore() error {

	if m.Config.Url == "" {
		return common.NewError(common.ErrCodeInternal, "Not found `url` in config file and `MONGO_URL` in env")
	}

	if m.Config.Database == "" {
		return common.NewError(common.ErrCodeInternal, "Not found `database` in config file and `MONGO_DATABASE` in env")
	}

	return nil
}

func (m *Mongo) Connect() error {
	session, err := mgo.DialWithTimeout(m.Config.Url, connTimeout)
	if err != nil {
		return common.NewErrorWithOther(common.ErrCodeInternal, err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	m.Connection = session.Copy() //必须使用copy，否则return后会自动关闭
	return nil
}

func (m *Mongo) GetConnection() interface{} {
	return m.Connection
}

func (m *Mongo) GetCollection(name string) (*mgo.Collection, error) {

	if name == "" {
		return nil, common.NewError(common.ErrCodeInternal, "name is empty")
	}

	collection := m.Connection.DB(m.Config.Database).C(name)

	return collection, nil
}

func (m *Mongo) Close() error {
	m.Connection.Close()
	return nil
}

func New() *Mongo {
	return &Mongo{}
}
