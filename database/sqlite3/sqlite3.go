package sqlite3

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lordking/toolbox/common"
)

type (
	Config struct {
		File string `json:"file" env:"SQLITE_FILE"`
	}

	SQLite struct {
		Config     *Config
		Connection *sql.DB
	}
)

func (m *SQLite) NewConfig() interface{} {
	m.Config = &Config{}
	return m.Config
}

func (m *SQLite) ValidateBefore() error {

	if m.Config.File == "" {
		return common.NewError(common.ErrCodeInternal, "Not found `file` in config file and `SQLITE_FILE` in env")
	}

	return nil
}

func (m *SQLite) Connect() error {

	var (
		db  *sql.DB
		err error
	)

	if db, err = sql.Open("sqlite3", m.Config.File); err != nil {
		return common.NewErrorWithOther(common.ErrCodeInternal, err)
	}

	m.Connection = db

	return nil
}

func (m *SQLite) GetConnection() interface{} {
	return m.Connection
}

func (m *SQLite) Close() error {
	if err := m.Connection.Close(); err != nil {
		return common.NewErrorWithOther(common.ErrCodeInternal, err)
	}

	return nil
}

func New() *SQLite {
	return &SQLite{}
}
