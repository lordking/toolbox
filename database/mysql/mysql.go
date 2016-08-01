package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/lordking/toolbox/common"
)

type (
	Config struct {
		Host         string `json:"host" env:"MYSQL_HOST"`
		Port         string `json:"port" env:"MYSQL_PORT"`
		Username     string `json:"username" env:"MYSQL_USERNAME"`
		Password     string `json:"password" env:"MYSQL_PASSWORD"`
		Database     string `json:"database" env:"MYSQL_DATABASE"`
		MaxOpenConns int    `json:"maxOpenConns" env:"MYSQL_MAXOPENCONNS"`
		MaxIdleConns int    `json:"maxIdleConns" env:"MYSQL_MAXIDLECONNS"`
	}

	MySQL struct {
		Config     *Config
		Connection *sql.DB
	}
)

func (m *MySQL) GetConfig() interface{} {
	return m.Config
}

func (m *MySQL) ValidateBefore() err {

	if m.Config.Host == "" {
		return common.NewError(common.ErrCodeInternal, "host not exist")
	}

	if m.Config.Port == "" {
		return common.NewError(common.ErrCodeInternal, "port not exist")
	}

	if m.Config.Username == "" {
		return common.NewError(common.ErrCodeInternal, "username not exist")
	}

	if m.Config.Password == "" {
		return common.NewError(common.ErrCodeInternal, "password not exist")
	}

	if m.Config.Database == "" {
		return common.NewError(common.ErrCodeInternal, "database not exist")
	}

	if m.Config.MaxOpenConns < 0 {
		return common.NewError(common.ErrCodeInternal, "MaxOpenConns not exist")
	}

	if m.Config.MaxIdleConns < 0 {
		return common.NewError(common.ErrCodeInternal, "MaxIdleConns not exist")
	}
}

func (m *MySQL) Connect() error {

	var (
		db  *sql.DB
		err error
	)

	if db, err = sql.Open("mysql", m.url()); err != nil {
		return common.NewErrorWithOther(common.ErrCodeInternal, err)
	}

	db.SetMaxOpenConns(m.Config.MaxOpenConns) //最大打开的连接数，默认值为0表示不限制
	db.SetMaxIdleConns(m.Config.MaxIdleConns) //闲置的连接数

	if err = db.Ping(); err != nil {
		return common.NewErrorWithOther(common.ErrCodeInternal, err)
	}

	m.Connection = db

	return nil
}

func (m *MySQL) GetConnection() interface{} {
	return m.Connection
}

func (m *MySQL) url() string {
	return m.Config.Username + ":" + m.Config.Password + "@tcp(" + m.Config.Host + ":" + m.Config.Port + ")/" + m.Config.Database + "?charset=utf8"
}

func (m *MySQL) Close() error {
	if err := m.Connection.Close(); err != nil {
		return common.NewErrorWithOther(common.ErrCodeInternal, err)
	}

	return nil
}

func New() *MySQL {
	return &MySQL{Config: &Config{}}
}
