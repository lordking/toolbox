package database

import "github.com/lordking/toolbox/common"

type Datasource interface {
	NewConfig() interface{}
	ValidateBefore() error
	Connect() error
	GetConnection() interface{}
	Close() error
}

var Instance Datasource

func CreateInstance(db Datasource, path string) error {

	config := db.NewConfig()

	var err error

	if err := common.ReadConfig(config, "./db.json"); err != nil {
		return err
	}

	if err := db.ValidateBefore(); err != nil {
		return err
	}

	Instance = db

	return err
}
