package goutils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
	Adapter      string
	Host         string `json:"host" osenv:"MYSQL_HOST"`
	Port         string `json:"port" osenv:"MYSQL_PORT"`
	Username     string `json:"username" osenv:"MYSQL_USERNAME"`
	Password     string `json:"password" osenv:"MYSQL_PASSWORD"`
	Database     string `json:"database" osenv:"MYSQL_DATABASE"`
	MaxOpenConns int    `json:"maxOpenConns" osenv:"MYSQL_MAXOPENCONNS"`
	MaxIdleConns int    `json:"maxIdleConns" osenv:"MYSQL_MAXIDLECONNS"`
}

func (this *MySQLConfig) GetAdapter() string {
	return this.Adapter
}

type MySQL struct {
	mysqldriver *sql.DB
	url         string
	config      *MySQLConfig
}

func (m *MySQL) InitDB() error {
	db, err := sql.Open("mysql", m.url)
	if err != nil {
		return ToError(ErrCodeInternal, err)
	}

	db.SetMaxOpenConns(m.config.MaxOpenConns) //最大打开的连接数，默认值为0表示不限制
	db.SetMaxIdleConns(m.config.MaxIdleConns) //闲置的连接数

	err = db.Ping()
	if err != nil {
		return ToError(ErrCodeInternal, err)
	}

	m.mysqldriver = db

	return nil
}

func (m *MySQL) GetConnection() (interface{}, error) {
	return m.mysqldriver, nil
}

func (m *MySQL) GetConfig() DataSourceConfig {
	return m.config
}

func NewMySQL(config *MySQLConfig) (*MySQL, error) {
	if config.Host == "" {
		return nil, NewError(ErrCodeInternal, "host not exist")
	}

	if config.Port == "" {
		return nil, NewError(ErrCodeInternal, "port not exist")
	}

	if config.Username == "" {
		return nil, NewError(ErrCodeInternal, "username not exist")
	}

	if config.Password == "" {
		return nil, NewError(ErrCodeInternal, "password not exist")
	}

	if config.Database == "" {
		return nil, NewError(ErrCodeInternal, "database not exist")
	}

	if config.MaxOpenConns < 0 {
		return nil, NewError(ErrCodeInternal, "MaxOpenConns not exist")
	}

	if config.MaxIdleConns < 0 {
		return nil, NewError(ErrCodeInternal, "MaxIdleConns not exist")
	}

	url := config.Username + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.Database + "?charset=utf8"

	return &MySQL{url: url, config: config}, nil
}
