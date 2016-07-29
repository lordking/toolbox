package goutils

import (
	"time"

	"gopkg.in/mgo.v2"
)

// === MongoConfig ===//

type MongoConfig struct {
	Adapter  string
	Url      string `json:"url" osenv:"MONGO_URL"`
	Database string `json:"database" osenv:"MONGO_DATABASE"`
}

func (this *MongoConfig) GetAdapter() string {
	return this.Adapter
}

// === Mongo ===//

var connTimeout = time.Second * 5

type Mongo struct {
	session *mgo.Session
	config  *MongoConfig
}

func (m *Mongo) InitDB() error {
	session, err := mgo.DialWithTimeout(m.config.Url, connTimeout)
	if err != nil {
		return ToError(ErrCodeInternal, err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	m.session = session.Copy() //必须使用copy，否则return后会自动关闭
	return nil
}

func (m *Mongo) GetConfig() DataSourceConfig {
	return m.config
}

func (m *Mongo) GetConnection() (interface{}, error) {
	return m.session, nil
}

func (m *Mongo) GetCollection(name string) (*mgo.Collection, error) {

	if name == "" {
		return nil, NewError(ErrCodeInternal, "name is empty")
	}

	collection := m.session.DB(m.config.Database).C(name)

	return collection, nil
}

func NewMongo(config *MongoConfig) (*Mongo, error) {

	if config.Url == "" {
		return nil, NewError(ErrCodeInternal, "Not found `MONGO_URL` in os env or `url` in config")
	}

	if config.Database == "" {
		return nil, NewError(ErrCodeInternal, "Not found `MONGO_DATABASE` in os env or `database` in config")
	}

	return &Mongo{config: config}, nil
}
