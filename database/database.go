package database

import "github.com/lordking/toolbox/common"

type Database interface {
	NewConfig() interface{}
	ValidateBefore() error
	Connect() error
	GetConnection() interface{}
	Close() error
}

func Configure(db Database, path string) error {

	if _, ok := db.(Database); !ok {
		return common.NewError(common.ErrCodeInternal, "Not found interface")
	}

	config := db.NewConfig()

	var err error

	if err := common.ReadConfig(config, path); err != nil {
		return err
	}

	if err := db.ValidateBefore(); err != nil {
		return err
	}

	return err
}
