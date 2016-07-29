package goutils

import (
	"fmt"
)

var DataSourceInstance DataSource

//==== DataSource config ====//

type DataSourceConfig interface {
	GetAdapter() string
}

type DataSourceConfigCommon struct {
	Adapter string `json:"adapter" osenv:"DB_ADAPTER"`
}

func (this *DataSourceConfigCommon) GetAdapter() string {
	return this.Adapter
}

//==== DataSource  ====//

type DataSource interface {
	InitDB() error
	GetConnection() (interface{}, error)
	GetConfig() DataSourceConfig
}

func NewDataSource(config DataSourceConfig) (DataSource, error) {

	adepter := config.GetAdapter()

	if adepter == "mysql" {
		config := config.(*MySQLConfig)
		return NewMySQL(config)
	} else if adepter == "mongo" {
		config := config.(*MongoConfig)
		return NewMongo(config)
	} else {
		errMsg := fmt.Sprintf("Not supported this adepter: %s", adepter)
		return nil, NewError(ErrCodeInternal, errMsg)
	}
}

func InitDataSource(configPath string) {

	//导入配置文件
	var data []byte
	var err error
	var config DataSourceConfig

	config = &DataSourceConfigCommon{}
	if configPath != "" {
		data, err = GetFileData(configPath)
		CheckFatal(err)

		err = ReadJSON(config, data)
		CheckFatal(err)
	}
	err = ReadOSEnv(config)
	CheckFatal(err)

	if config.GetAdapter() == "mysql" {
		config = &MySQLConfig{Adapter: "mysql"}

	} else if config.GetAdapter() == "mongo" {
		config = &MongoConfig{Adapter: "mongo"}

	} else {
		errMsg := fmt.Sprintf("Not found this adapter: %s", config.GetAdapter())
		err := NewError(ErrCodeInternal, errMsg)
		CheckFatal(err)
	}

	err = ReadJSON(config, data)
	CheckFatal(err)

	err = ReadOSEnv(config)
	CheckFatal(err)

	DataSourceInstance, err = NewDataSource(config)
	CheckFatal(err)

	err = DataSourceInstance.InitDB()
	CheckFatal(err)

}
