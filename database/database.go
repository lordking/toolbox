package database

import (
	"encoding/json"

	"github.com/lordking/toolbox/common"
)

type Database interface {
	NewConfig() interface{}
	ValidateBefore() error
	Connect() error
	GetConnection() interface{}
	Close() error
}

func ConfigureWithPath(db Database, path string) error {

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

func Configure(db Database, obj interface{}) error {

	config := db.NewConfig()

	var err error

	data, _ := json.Marshal(obj)
	if err := json.Unmarshal(data, &config); err != nil {
		return common.NewError(common.ErrCodeInternal, err.Error())
	}

	if err := db.ValidateBefore(); err != nil {
		return err
	}

	return err
}
