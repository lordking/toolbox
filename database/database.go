package database

import "github.com/lordking/toolbox/common"

type Database interface {
	NewConfig() interface{}
	ValidateBefore() error
	Connect() error
	GetConnection() interface{}
	Close() error
}

func Configure(configKey string, db Database) error {

	config := db.NewConfig()
	common.ReadConfigFromKey(configKey, config)
	err := db.ValidateBefore()

	return err
}
